#!/bin/bash

# 生成CA私钥和证书
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=OpenVPN-CA"

# 生成服务器私钥和证书签名请求
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=OpenVPN-Server"

# 使用CA证书签名服务器证书
openssl x509 -req -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

# 清理临时文件
rm server.csr ca.srl

# 设置适当的权限
chmod 600 *.key
chmod 644 *.crt

echo "证书生成完成！" 