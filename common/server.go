package common

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "openvpn-admin-go/constants"
)

// ServerConfig 结构体表示 OpenVPN 的配置
type ServerConfig struct {
    Port                int      `json:"port"`
    Proto               string   `json:"proto"`
    Dev                 string   `json:"dev"`
    CA                  string   `json:"ca"`
    Cert                string   `json:"cert"`
    Key                 string   `json:"key"`
    DH                  string   `json:"dh"`
    Server              string   `json:"server"`
    IfconfigPoolPersist string   `json:"ifconfig_pool_persist"`
    Push                []string `json:"push"`
    Keepalive           string   `json:"keepalive"`
    TLSAuth             string   `json:"tls_auth"`
    Cipher              string   `json:"cipher"`
    User                string   `json:"user"`
    Group               string   `json:"group"`
    PersistKey          bool     `json:"persist_key"`
    PersistTun          bool     `json:"persist_tun"`
    Status              string   `json:"status"`
    Verb                int      `json:"verb"`
}

// Load 加载配置文件
func Load(path string) []string {
    if path == "" {
        path = constants.OpenVPNConfigPath
    }
    file, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
        return nil
    }

    return lines
}

// Save 保存配置文件
func Save(lines []string, path string) bool {
    if path == "" {
        path = constants.OpenVPNConfigPath
    }
    file, err := os.Create(path)
    if err != nil {
        fmt.Println(err)
        return false
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(writer, line)
    }

    if err := writer.Flush(); err != nil {
        fmt.Println(err)
        return false
    }
    return true
}

// GetConfig 获取配置
func GetConfig() *ServerConfig {
    lines := Load("")
    config := &ServerConfig{}

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 1 {
            continue
        }
        switch fields[0] {
        case "port":
            fmt.Sscanf(fields[1], "%d", &config.Port)
        case "proto":
            config.Proto = fields[1]
        case "dev":
            config.Dev = fields[1]
        case "ca":
            config.CA = fields[1]
        case "cert":
            config.Cert = fields[1]
        case "key":
            config.Key = fields[1]
        case "dh":
            config.DH = fields[1]
        case "server":
            config.Server = fields[1] + " " + fields[2]
        case "ifconfig-pool-persist":
            config.IfconfigPoolPersist = fields[1]
        case "push":
            config.Push = append(config.Push, strings.Join(fields[1:], " "))
        case "keepalive":
            config.Keepalive = fields[1] + " " + fields[2]
        case "tls-auth":
            config.TLSAuth = fields[1] + " " + fields[2]
        case "cipher":
            config.Cipher = fields[1]
        case "user":
            config.User = fields[1]
        case "group":
            config.Group = fields[1]
        case "persist-key":
            config.PersistKey = true
        case "persist-tun":
            config.PersistTun = true
        case "status":
            config.Status = fields[1]
        case "verb":
            fmt.Sscanf(fields[1], "%d", &config.Verb)
        }
    }
    return config
}

// UpdateConfig 更新配置文件
func UpdateConfig(config *ServerConfig) bool {
    lines := Load("")
    newLines := []string{}

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) < 1 {
            newLines = append(newLines, line)
            continue
        }
        switch fields[0] {
        case "port":
            newLines = append(newLines, fmt.Sprintf("port %d", config.Port))
        case "proto":
            newLines = append(newLines, fmt.Sprintf("proto %s", config.Proto))
        case "dev":
            newLines = append(newLines, fmt.Sprintf("dev %s", config.Dev))
        case "ca":
            newLines = append(newLines, fmt.Sprintf("ca %s", config.CA))
        case "cert":
            newLines = append(newLines, fmt.Sprintf("cert %s", config.Cert))
        case "key":
            newLines = append(newLines, fmt.Sprintf("key %s", config.Key))
        case "dh":
            newLines = append(newLines, fmt.Sprintf("dh %s", config.DH))
        case "server":
            newLines = append(newLines, fmt.Sprintf("server %s", config.Server))
        case "ifconfig-pool-persist":
            newLines = append(newLines, fmt.Sprintf("ifconfig-pool-persist %s", config.IfconfigPoolPersist))
        case "push":
            for _, push := range config.Push {
                newLines = append(newLines, fmt.Sprintf("push \"%s\"", push))
            }
        case "keepalive":
            newLines = append(newLines, fmt.Sprintf("keepalive %s", config.Keepalive))
        case "tls-auth":
            newLines = append(newLines, fmt.Sprintf("tls-auth %s", config.TLSAuth))
        case "cipher":
            newLines = append(newLines, fmt.Sprintf("cipher %s", config.Cipher))
        case "user":
            newLines = append(newLines, fmt.Sprintf("user %s", config.User))
        case "group":
            newLines = append(newLines, fmt.Sprintf("group %s", config.Group))
        case "persist-key":
            if config.PersistKey {
                newLines = append(newLines, "persist-key")
            }
        case "persist-tun":
            if config.PersistTun {
                newLines = append(newLines, "persist-tun")
            }
        case "status":
            newLines = append(newLines, fmt.Sprintf("status %s", config.Status))
        case "verb":
            newLines = append(newLines, fmt.Sprintf("verb %d", config.Verb))
        default:
            newLines = append(newLines, line)
        }
    }

    return Save(newLines, "")
}

// UpdatePort 更新端口配置
func UpdatePort(newPort int) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Port = newPort
    return UpdateConfig(config)
}

// UpdateProto 更新协议配置
func UpdateProto(newProto string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Proto = newProto
    return UpdateConfig(config)
}

// UpdateDev 更新设备配置
func UpdateDev(newDev string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Dev = newDev
    return UpdateConfig(config)
}

// UpdateOpenSSL 更新 OpenSSL 配置
func UpdateOpenSSL(newCA, newCert, newKey, newDH string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.CA = newCA
    config.Cert = newCert
    config.Key = newKey
    config.DH = newDH
    return UpdateConfig(config)
}

// UpdateServer 更新服务器配置
func UpdateServer(newServer string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Server = newServer
    return UpdateConfig(config)
}

// UpdatePush 更新推送选项
func UpdatePush(newPush []string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Push = newPush
    return UpdateConfig(config)
}

// UpdateKeepalive 更新 keepalive 配置
func UpdateKeepalive(newKeepalive string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Keepalive = newKeepalive
    return UpdateConfig(config)
}

// UpdateTLSAuth 更新 TLS Auth 配置
func UpdateTLSAuth(newTLSAuth string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.TLSAuth = newTLSAuth
    return UpdateConfig(config)
}

// UpdateCipher 更新加密配置
func UpdateCipher(newCipher string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Cipher = newCipher
    return UpdateConfig(config)
}

// UpdateUser 更新用户配置
func UpdateUser(newUser string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.User = newUser
    return UpdateConfig(config)
}

// UpdateGroup 更新组配置
func UpdateGroup(newGroup string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Group = newGroup
    return UpdateConfig(config)
}

// UpdatePersistKey 更新 persist-key 配置
func UpdatePersistKey(newPersistKey bool) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.PersistKey = newPersistKey
    return UpdateConfig(config)
}

// UpdatePersistTun 更新 persist-tun 配置
func UpdatePersistTun(newPersistTun bool) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.PersistTun = newPersistTun
    return UpdateConfig(config)
}

// UpdateStatus 更新状态日志配置
func UpdateStatus(newStatus string) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Status = newStatus
    return UpdateConfig(config)
}

// UpdateVerb 更新日志等级配置
func UpdateVerb(newVerb int) bool {
    config := GetConfig()
    if config == nil {
        return false
    }
    config.Verb = newVerb
    return UpdateConfig(config)
}
