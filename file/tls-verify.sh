#!/bin/bash
# OpenVPN tls-verify 脚本：按证书 CN 拉黑（替代旧的 auth-user-pass-verify + 假密码方案）。
# OpenVPN 对证书链每一层调用：tls-verify <cert_depth> <X509_subject>
# 并为每层设置环境变量 X509_<depth>_CN。我们只校验叶子证书(depth 0)。
#
# 语义：CN 在黑名单里 → exit 1（拒绝本次 TLS 握手，连不上）；否则 exit 0（放行）。
# 本脚本只读检查，不修改黑名单、不 kill——「暂停那一刻」的即时断开由 PauseClient
# 经管理接口发 kill 完成；黑名单的增删由 暂停/恢复 接口负责。
set -u

cert_depth="${1:-}"
x509_subject="${2:-}"

# 只校验叶子证书；CA / 中间层(depth>0)直接放行。
if [ "$cert_depth" != "0" ]; then
    exit 0
fi

# CN 优先取 OpenVPN 注入的环境变量；回退到从 subject 串里解析 /CN=xxx。
cn="${X509_0_CN:-}"
if [ -z "$cn" ] && [ -n "$x509_subject" ]; then
    # subject 形如 "C=..,O=..,CN=username"；抓最后一个 CN= 字段，去掉前后空白。
    cn=$(printf '%s' "$x509_subject" | sed -n 's/.*CN=\([^,/]*\).*/\1/p' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
fi

BLACKLIST_FILE="${OPENVPN_BLACKLIST_FILE:-/etc/openvpn/server/blacklist.txt}"

# 日志：OpenVPN 以降权用户(nobody)运行本脚本，/var/log 往往不可写。
# 选一个可写路径，并保证「记日志」永远不会让脚本非零退出（set -u 下也安全）。
LOG_FILE="${OPENVPN_TLS_VERIFY_LOG:-/var/log/openvpn-tls-verify.log}"
if ! { [ -w "$LOG_FILE" ] 2>/dev/null || { [ -d "$(dirname "$LOG_FILE")" ] && [ -w "$(dirname "$LOG_FILE")" ]; } 2>/dev/null; }; then
    LOG_FILE="/tmp/openvpn-tls-verify.log"
fi
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S'): $*" >> "$LOG_FILE" 2>/dev/null || true
}

# CN 解析不出来：放行（纯证书已经由 OpenVPN 的 TLS 校验保证身份；这里只做黑名单叠加）。
if [ -z "$cn" ]; then
    log "depth0 but empty CN (subject=$x509_subject); allowing."
    exit 0
fi

if [ ! -f "$BLACKLIST_FILE" ]; then
    log "Blacklist file $BLACKLIST_FILE not found. Allowing CN $cn."
    exit 0
fi

if grep -qx "$cn" "$BLACKLIST_FILE"; then
    log "CN $cn is blacklisted. Denying TLS handshake."
    exit 1
fi

log "CN $cn not blacklisted. Allowing."
exit 0
