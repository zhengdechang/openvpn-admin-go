#!/bin/bash

# 定义操作变量, 0为否, 1为是
help=0
remove=0
update=0
add_user=0
del_user=0
username=""

[[ -e /var/lib/openvpn-manager ]] && update=1

# CentOS 临时取消别名
[[ -f /etc/redhat-release && -z $(echo $SHELL | grep zsh) ]] && unalias -a

[[ -z $(echo $SHELL | grep zsh) ]] && shell_way="bash" || shell_way="zsh"

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

####### get params #########
while [[ $# > 0 ]]; do
    key="$1"
    case $key in
    --remove)
        remove=1
        ;;
    -h | --help)
        help=1
        ;;
    *)
        # unknown option
        ;;
    esac
    shift # past argument or value
done
#############################

help() {
    echo "bash $0 [-h|--help]"
    echo "  -h, --help           Show help"
    return 0
}

removeOpenVPN() {
    # 移除OpenVPN
    apt-get remove --purge -y openvpn || yum remove -y openvpn || dnf remove -y openvpn
    rm -rf /etc/openvpn
    systemctl daemon-reload
    colorEcho ${green} "uninstall success!"
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
        ${package_manager} update
        ${package_manager} install openvpn openssl wget -y
    fi
}

configureOpenVPN() {
    mkdir -p /etc/openvpn/server
    cd /etc/openvpn/server || exit

    # 生成 CA 证书和密钥
    openssl genpkey -algorithm RSA -out ca.key
    openssl req -x509 -new -key ca.key -sha256 -out ca.crt -days 3650 -subj "/C=US/ST=California/L=San Francisco/O=MyCompany/OU=IT/CN=myvpn.com"

    # 生成服务器证书和密钥
    openssl genpkey -algorithm RSA -out server.key
    openssl req -new -key server.key -out server.csr -subj "/C=US/ST=California/L=San Francisco/O=MyCompany/OU=IT/CN=myvpn.com"
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650 -sha256

    # 生成 Diffie-Hellman 参数
    openssl dhparam -out dh.pem 2048

    # 生成 ta.key
    openvpn --genkey secret ta.key

    # 手动创建 OpenVPN 服务器配置文件
    cat <<EOF > /etc/openvpn/server/server.conf
port 1194
proto tcp
dev tun
ca /etc/openvpn/server/ca.crt
cert /etc/openvpn/server/server.crt
key /etc/openvpn/server/server.key
dh /etc/openvpn/server/dh.pem
server 10.8.0.0 255.255.255.0
ifconfig-pool-persist ipp.txt
push "redirect-gateway def1 bypass-dhcp"
push "dhcp-option DNS 8.8.8.8"
push "dhcp-option DNS 8.8.4.4"
keepalive 10 120
tls-auth /etc/openvpn/server/ta.key 0
cipher AES-256-CBC
user nobody
group nogroup
persist-key
persist-tun
status openvpn-status.log
verb 3
script-security 2
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
    cd /etc/openvpn/server || exit

    # 生成客户端证书和密钥
    openssl genpkey -algorithm RSA -out ${user}.key
    openssl req -new -key ${user}.key -out ${user}.csr -subj "/C=US/ST=California/L=San Francisco/O=MyCompany/OU=IT/CN=${user}"
    openssl x509 -req -in ${user}.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out ${user}.crt -days 3650 -sha256

    # 生成客户端配置文件
    cat <<EOF > /etc/openvpn/client/${user}.ovpn
client
tls-client
auth SHA256
cipher AES-256-CBC
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
$(cat /etc/openvpn/server/ca.crt)
</ca>
<tls-auth>
$(cat /etc/openvpn/server/ta.key)
</tls-auth>
<cert>
$(cat /etc/openvpn/server/${user}.crt)
</cert>
<key>
$(cat /etc/openvpn/server/${user}.key)
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
    cd /etc/openvpn/server || exit

    rm -f ${user}.crt ${user}.csr ${user}.key
    rm -f /etc/openvpn/client/${user}.ovpn

    colorEcho ${green} "User ${user} deleted successfully!"
}

main() {
    while true; do
        echo "请选择操作:"
        echo "1) 安装 OpenVPN"
        echo "2) 新增用户"
        echo "3) 删除用户"
        echo "4) 退出"
        read -rp "输入选项 [1-4]: " option
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
            echo "退出"
            exit 0
            ;;
        *)
            echo "无效的选项，请重新输入"
            ;;
        esac
    done
}

main