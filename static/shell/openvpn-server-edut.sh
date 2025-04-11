#!/bin/bash

# 定义变量
OPENVPN_KEY_DIR="/etc/openvpn/keys"
OPENVPN_BASE_DIR="/etc/openvpn"
OPENVPN_RSA_BITS=2048
OPENVPN_USE_PREGNERATED_DH_PARAMS=false
CI_BUILD=false

####### color code ########
red="31m"
green="32m"
yellow="33m"
blue="36m"
fuchsia="35m"

colorEcho() {
    color=$1
    echo -e "\033[${color}${@:2}\033[0m"
}

checkSys() {
    # 检查是否为Root
    [ $(id -u) != "0" ] && { colorEcho ${red} "Error: You must be root to run this script"; exit 1; }

    arch=$(uname -m 2>/dev/null)
    if [[ $arch != x86_64 && $arch != aarch64 ]]; then
        colorEcho $yellow "not support $arch machine".
        exit 1
    fi

    if [[ `command -v apt-get` ]]; then
        package_manager='apt-get'
    elif [[ `command -v dnf` ]]; then
        package_manager='dnf'
    elif [[ `command -v yum` ]];then
        package_manager='yum'
    else
        colorEcho $red "Not support OS!"
        exit 1
    fi

    # 缺失/usr/local/bin路径时自动添加
    [[ -z `echo $PATH | grep /usr/local/bin` ]] && { echo 'export PATH=$PATH:/usr/local/bin' >>/etc/bashrc; source /etc/bashrc; }
}

# 安装依赖
installDependent() {
    if [[ ${package_manager} == 'dnf' || ${package_manager} == 'yum' ]]; then
        ${package_manager} install epel-release -y
        ${package_manager} install openvpn openssl wget -y
    else
        ${package_manager} update -y
        ${package_manager} install openvpn openssl wget -y
    fi
}

configureOpenVPN() {
    mkdir -p "${OPENVPN_KEY_DIR}/certs"
    chmod 0755 "${OPENVPN_KEY_DIR}"

    # 生成openssl server/ca扩展文件
    cat <<EOF > "${OPENVPN_KEY_DIR}/openssl-server.ext"
basicConstraints = CA:FALSE
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
extendedKeyUsage = serverAuth
keyUsage = digitalSignature,keyEncipherment
EOF

    cat <<EOF > "${OPENVPN_KEY_DIR}/openssl-ca.ext"
basicConstraints = CA:TRUE
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
keyUsage = cRLSign, keyCertSign
EOF

    cat <<EOF > "${OPENVPN_KEY_DIR}/openssl-client.ext"
basicConstraints = CA:FALSE
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
extendedKeyUsage = clientAuth
keyUsage = digitalSignature
EOF

    chmod 0400 "${OPENVPN_KEY_DIR}/openssl-server.ext" "${OPENVPN_KEY_DIR}/openssl-ca.ext" "${OPENVPN_KEY_DIR}/openssl-client.ext"

    # 生成CA密钥和证书
    openssl req -nodes -newkey rsa:${OPENVPN_RSA_BITS} -keyout "${OPENVPN_KEY_DIR}/ca.key" -out "${OPENVPN_KEY_DIR}/ca.csr" -subj "/CN=awasome-openvpn/"
    chmod 0400 "${OPENVPN_KEY_DIR}/ca.key"
    openssl x509 -req -in "${OPENVPN_KEY_DIR}/ca.csr" -out "${OPENVPN_KEY_DIR}/ca.crt" -signkey "${OPENVPN_KEY_DIR}/ca.key" -days 3650 -extfile "${OPENVPN_KEY_DIR}/openssl-ca.ext"

    # 生成服务器密钥和证书
    openssl req -nodes -newkey rsa:${OPENVPN_RSA_BITS} -keyout "${OPENVPN_KEY_DIR}/server.key" -out "${OPENVPN_KEY_DIR}/server.csr" -subj "/CN=awasome-openvpn/"
    chmod 0400 "${OPENVPN_KEY_DIR}/server.key"
    openssl x509 -req -in "${OPENVPN_KEY_DIR}/server.csr" -out "${OPENVPN_KEY_DIR}/server.crt" -CA "${OPENVPN_KEY_DIR}/ca.crt" -CAkey "${OPENVPN_KEY_DIR}/ca.key" -days 3650 -CAcreateserial -extfile "${OPENVPN_KEY_DIR}/openssl-server.ext"

    # 生成tls-auth密钥
    openvpn --genkey --secret "${OPENVPN_KEY_DIR}/ta.key"
    chmod 0400 "${OPENVPN_KEY_DIR}/ta.key"

    # 生成DH参数
    if [ "${OPENVPN_USE_PREGNERATED_DH_PARAMS}" = true ]; then
        cp dh.pem "${OPENVPN_KEY_DIR}/dh.pem"
        chmod 0400 "${OPENVPN_KEY_DIR}/dh.pem"
    else
        openssl dhparam -out "${OPENVPN_KEY_DIR}/dh.pem" ${OPENVPN_RSA_BITS}
    fi

    # 安装ca.conf配置文件
    cat <<EOF > "${OPENVPN_KEY_DIR}/ca.conf"
# OpenVPN CA configuration
[ ca ]
default_ca = CA_default

[ CA_default ]
dir = ${OPENVPN_KEY_DIR}
certs = \$dir/certs
new_certs_dir = \$dir/newcerts
database = \$dir/index.txt
serial = \$dir/serial
RANDFILE = \$dir/private/.rand

private_key = \$dir/private/ca.key
certificate = \$dir/ca.crt

crlnumber = \$dir/crlnumber
crl = \$dir/crl.pem
crl_extensions = crl_ext
default_crl_days = 30

default_md = sha256

name_opt = ca_default
cert_opt = ca_default
default_days = 365
preserve = no
policy = policy_strict

[ policy_strict ]
countryName = match
stateOrProvinceName = match
organizationName = match
organizationalUnitName = optional
commonName = supplied
emailAddress = optional

[ crl_ext ]
# CRL extension (optional)
EOF
    chmod 0744 "${OPENVPN_KEY_DIR}/ca.conf"

    # 创建初始证书吊销列表序列号
    if [ ! -f "${OPENVPN_KEY_DIR}/crlnumber" ]; then
        echo "00" > "${OPENVPN_KEY_DIR}/crlnumber"
    fi

    # 创建证书吊销列表数据库
    if [ ! -f "${OPENVPN_KEY_DIR}/index.txt" ]; then
        touch "${OPENVPN_KEY_DIR}/index.txt"
        chmod 0644 "${OPENVPN_KEY_DIR}/index.txt"
    fi

    # 设置证书吊销列表
    openssl ca -config "${OPENVPN_KEY_DIR}/ca.conf" -gencrl -keyfile "${OPENVPN_KEY_DIR}/ca.key" -cert "${OPENVPN_KEY_DIR}/ca.crt" -out "${OPENVPN_KEY_DIR}/crl.pem"

    # 安装crl-cron脚本
    cat <<EOF > "${OPENVPN_BASE_DIR}/crl-cron.sh"
#!/bin/bash
# Check if CRL needs to be renewed

OPENVPN_KEY_DIR="${OPENVPN_KEY_DIR}"

cd \${OPENVPN_KEY_DIR}
source ca.conf

openssl ca -config ca.conf -gencrl -keyfile private/ca.key -cert ca.crt -out crl.pem

echo "CRL has been renewed."
EOF
    chmod 0744 "${OPENVPN_BASE_DIR}/crl-cron.sh"

    # 检查crontab
    if ! command -v crontab &> /dev/null; then
        if [ -f /etc/redhat-release ]; then
            yum install -y cronie
        fi
    fi

    # 添加cron任务每周检查CRL是否需要更新
    if [ "${CI_BUILD}" = false ]; then
        (crontab -l 2>/dev/null; echo "0 0 * * 6 sh ${OPENVPN_BASE_DIR}/crl-cron.sh") | crontab -
    fi

    # 创建OpenVPN服务器配置文件
    cat <<EOF > /etc/openvpn/server.conf
port 1194
proto tcp
dev tun
ca ${OPENVPN_KEY_DIR}/ca.crt
cert ${OPENVPN_KEY_DIR}/server.crt
key ${OPENVPN_KEY_DIR}/server.key
dh ${OPENVPN_KEY_DIR}/dh.pem
server 10.8.0.0 255.255.255.0
ifconfig-pool-persist ipp.txt
push "dhcp-option DNS 8.8.8.8"
push "dhcp-option DNS 8.8.4.4"
keepalive 10 120
tls-auth ${OPENVPN_KEY_DIR}/ta.key 0
auth SHA256
cipher AES-256-GCM
proto tcp
user nobody
group nogroup
persist-key
persist-tun
status openvpn-status.log
verb 3
script-security 2
crl-verify ${OPENVPN_KEY_DIR}/crl.pem
EOF

    systemctl start openvpn@server
    systemctl enable openvpn@server

    echo 1 > /proc/sys/net/ipv4/ip_forward
    sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/g' /etc/sysctl.conf

    iptables -t nat -A POSTROUTING -s 10.8.0.0/24 -o eth0 -j MASQUERADE
    iptables-save > /etc/iptables.rules

    echo "iptables-restore < /etc/iptables.rules" >> /etc/rc.local

    systemctl restart openvpn@server
    colorEcho ${green} "OpenVPN installation and configuration completed!"
}

generateClientConfig() {
    local user=$1
    local server_ip=$(hostname -I | awk '{print $1}')
    mkdir -p /etc/openvpn/client
    cd "${OPENVPN_KEY_DIR}" || exit

    # 生成客户端证书和密钥
    openssl genpkey -algorithm RSA -out ${user}.key
    openssl req -new -key ${user}.key -out ${user}.csr -subj "/C=US/ST=California/L=San Francisco/O=MyCompany/OU=IT/CN=${user}"
    openssl x509 -req -in ${user}.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out ${user}.crt -days 3650 -sha256 -extfile "${OPENVPN_KEY_DIR}/openssl-client.ext"

    # 生成客户端配置文件
    cat <<EOF > /etc/openvpn/client/${user}.ovpn
client
tls-client
auth SHA256
cipher AES-256-GCM
remote-cert-tls server
tls-version-min 1.2

proto tcp
remote ${server_ip} 1194
dev tun

resolv-retry infinite
nobind
keepalive 10 120
persist-key
persist-tun
verb 3

route-method exe
route-delay 2

key-direction 1
<ca>
$(cat ${OPENVPN_KEY_DIR}/ca.crt)
</ca>
<tls-auth>
$(cat ${OPENVPN_KEY_DIR}/ta.key)
</tls-auth>
<cert>
$(cat ${OPENVPN_KEY_DIR}/${user}.crt)
</cert>
<key>
$(cat ${OPENVPN_KEY_DIR}/${user}.key)
</key>
EOF

    colorEcho ${green} "User ${user} added successfully!"
}

addOpenVPNUser() {
    local user=$1
    generateClientConfig "${user}"
}

delOpenVPNUser() {
    local user=$1
    cd "${OPENVPN_KEY_DIR}" || exit

    # 确保证书文件存在
    if [ ! -f "${OPENVPN_KEY_DIR}/${user}.crt" ]; then
        colorEcho ${red} "Error: Certificate file ${OPENVPN_KEY_DIR}/${user}.crt does not exist"
        return 1
    fi

    # 确保 crlnumber 文件存在
    if [ ! -f "${OPENVPN_KEY_DIR}/crlnumber" ]; then
        echo "01" > "${OPENVPN_KEY_DIR}/crlnumber"
    fi

    # 吊销证书并更新 CRL
    openssl ca -config "${OPENVPN_KEY_DIR}/ca.conf" -revoke "${OPENVPN_KEY_DIR}/${user}.crt" -keyfile "${OPENVPN_KEY_DIR}/ca.key" -cert "${OPENVPN_KEY_DIR}/ca.crt"
    if [ $? -ne 0 ]; then
        colorEcho ${red} "Error: Failed to revoke certificate ${OPENVPN_KEY_DIR}/${user}.crt"
        return 1
    fi

    openssl ca -config "${OPENVPN_KEY_DIR}/ca.conf" -gencrl -keyfile "${OPENVPN_KEY_DIR}/ca.key" -cert "${OPENVPN_KEY_DIR}/ca.crt" -out "${OPENVPN_KEY_DIR}/crl.pem"
    if [ $? -ne 0 ]; then
        colorEcho ${red} "Error: Failed to generate CRL"
        return 1
    fi

    # 删除证书文件
    rm -f ${user}.crt ${user}.csr ${user}.key
    rm -f /etc/openvpn/client/${user}.ovpn

    # 发送 SIGHUP 信号给 OpenVPN 进程以重新加载配置和 CRL 文件
    pkill -SIGHUP openvpn

    colorEcho ${green} "User ${user} deleted successfully!"
}

viewConfig() {
    read -rp "请输入要查看配置的用户名: " username
    if [ -f "/etc/openvpn/client/${username}.ovpn" ]; then
        cat "/etc/openvpn/client/${username}.ovpn"
    else
        colorEcho ${red} "Error: Client configuration file for user ${username} does not exist"
    fi
}

restartService() {
    systemctl restart openvpn@server
    colorEcho ${green} "OpenVPN service restarted successfully!"
}

stopService() {
    systemctl stop openvpn@server
    colorEcho ${green} "OpenVPN service stopped successfully!"
}

removeOpenVPN() {
    # 停止所有 OpenVPN 服务
    systemctl stop openvpn@server
    systemctl disable openvpn@server

    # 移除 OpenVPN
    if [[ ${package_manager} == 'apt-get' ]]; then
        apt-get remove --purge -y openvpn
    elif [[ ${package_manager} == 'yum' || ${package_manager} == 'dnf' ]]; then
        yum remove -y openvpn || dnf remove -y openvpn
    fi

    # 删除 OpenVPN 配置文件
    rm -rf /etc/openvpn

    # 重新加载 systemd 守护进程
    systemctl daemon-reload

    colorEcho ${green} "OpenVPN uninstalled successfully!"
}

modifyOpenVPNPort() {
    read -rp "请输入新的端口 (当前: 1194): " port

    # 修改OpenVPN服务器配置文件
    if [ -f /etc/openvpn/server.conf ]; then
        sed -i "s/^port .*/port ${port}/" /etc/openvpn/server.conf

        # 重新启动OpenVPN服务
        systemctl restart openvpn@server
        colorEcho ${green} "OpenVPN端口已修改并重新启动服务"
    else
        colorEcho ${red} "Error: OpenVPN配置文件不存在"
    fi
}

modifyOpenVPNHost() {
    read -rp "请输入新的服务器IP或域名 (当前: 10.8.0.0): " server_ip

    # 修改OpenVPN服务器配置文件
    if [ -f /etc/openvpn/server.conf ]; then
        sed -i "s/^server .*/server ${server_ip} 255.255.255.0/" /etc/openvpn/server.conf

        # 重新启动OpenVPN服务
        systemctl restart openvpn@server
        colorEcho ${green} "OpenVPN主机已修改并重新启动服务"
    else
        colorEcho ${red} "Error: OpenVPN配置文件不存在"
    fi
}

addRoute() {
    read -rp "请输入目的网络 (例如: 192.168.1.0): " dest_network
    read -rp "请输入子网掩码 (例如: 255.255.255.0): " netmask

    # 添加新的路由到OpenVPN服务器配置文件
    if [ -f /etc/openvpn/server.conf ]; then
        echo "push \"route ${dest_network} ${netmask}\"" >> /etc/openvpn/server.conf

        # 重新启动OpenVPN服务
        systemctl restart openvpn@server
        colorEcho ${green} "新的路由已添加并重新启动服务"
    else
        colorEcho ${red} "Error: OpenVPN配置文件不存在"
    fi
}

main() {
    while true; do
        echo "请选择操作:"
        echo "1) 安装 OpenVPN"
        echo "2) 新增用户"
        echo "3) 删除用户"
        echo "4) 卸载 OpenVPN"
        echo "5) 查看客户端配置"
        echo "6) 重启服务"
        echo "7) 停止服务"
        echo "8) 修改OpenVPN端口"
        echo "9) 修改OpenVPN主机"
        echo "10) 添加新路由"
        echo "11) 退出"
        read -rp "输入选项 [1-11]: " option
        case $option in
        1)
            echo "正在安装 OpenVPN..."
            checkSys
            installDependent
            configureOpenVPN
            ;;
        2)
            read -rp "请输入要新增的用户名: " username
            addOpenVPNUser "${username}"
            ;;
        3)
            read -rp "请输入要删除的用户名: " username
            delOpenVPNUser "${username}"
            ;;
        4)
            echo "正在卸载 OpenVPN..."
            removeOpenVPN
            ;;
        5)
            echo "查看客户端配置文件..."
            viewConfig
            ;;
        6)
            echo "重启 OpenVPN 服务..."
            restartService
            ;;
        7)
            echo "停止 OpenVPN 服务..."
            stopService
            ;;
        8)
            echo "修改OpenVPN端口..."
            modifyOpenVPNPort
            ;;
        9)
            echo "修改OpenVPN主机..."
            modifyOpenVPNHost
            ;;
        10)
            echo "添加新路由..."
            addRoute
            ;;
        11)
            echo "退出"
            exit 0
            ;;
        *)
            echo "无效的选项，请重新输入"
            ;;
        esac
    done
}