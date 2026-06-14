package openvpn

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/utils"
)

// safeUsername 限制用户名只含证书/文件名安全字符，挡住 shell 注入：
// 吊销流程会把用户名拼进 bash 命令（/CN=<user>），必须先校验。
var safeUsername = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

// serverDir 返回 /etc/openvpn/server（所有服务端文件所在目录）。
func serverDir() string {
	return filepath.Dir(constants.ServerConfigPath)
}

// EnsureServerHelperFiles 把镜像里的辅助脚本/配置（tls-verify.sh、crl.cnf）刷到
// 持久卷 /etc/openvpn/server/ 下。
//
// 为什么需要它：/etc/openvpn 是持久卷，辅助文件只在「全新初始化」(cmd/environment.go
// generateCertificates) 时从 <cwd>/file/ 复制一次。已存在的旧卷里没有 tls-verify.sh /
// crl.cnf——一旦 server.conf 渲染出 `tls-verify` / `crl-verify` 引用它们却找不到文件，
// OpenVPN 会拒绝启动 = 全员锁死。所以每次改配置/启动前都同步一遍（幂等覆盖，
// 让镜像里更新过的脚本逻辑也能刷进旧卷）。
//
// 只同步无状态的代码/配置文件；blacklist.txt 是运行时状态（暂停名单），绝不覆盖。
func EnsureServerHelperFiles() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}
	srcDir := filepath.Join(wd, "file")
	dstDir := serverDir()

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建服务端目录失败: %v", err)
	}

	for _, name := range []string{"tls-verify.sh", "crl.cnf"} {
		src := filepath.Join(srcDir, name)
		if _, statErr := os.Stat(src); statErr != nil {
			// 源文件不在（例如本机开发、非容器环境）→ 跳过，不报错。
			fmt.Printf("辅助文件源缺失，跳过同步: %s\n", src)
			continue
		}
		dst := filepath.Join(dstDir, name)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("同步辅助文件 %s 失败: %v", name, err)
		}
		// 0755：tls-verify.sh 要被 OpenVPN（降权后 nobody）读取并执行。
		if err := os.Chmod(dst, 0755); err != nil {
			return fmt.Errorf("设置辅助文件权限失败 %s: %v", dst, err)
		}
		fmt.Printf("已同步辅助文件: %s\n", dst)
	}
	return nil
}

// EnsureCRLSetup 幂等地搭好 openssl CA 数据库并保证 crl.pem 存在。
//
// 关键防锁死：server.conf 一旦带 `crl-verify <file>`，该文件缺失/格式错/过期都会让
// OpenVPN 拒绝所有连接。所以必须「先有一份有效(初始为空)的 CRL，再渲染出 crl-verify」。
// 本函数在写 server.conf / 重启 / 冷启动之前调用，确保 crl.pem 已就位。
//
// 没有 CA（全新环境还没 generateCertificates）时直接返回 nil——此时也没有 server.conf
// 引用 CRL，无需求。
func EnsureCRLSetup() error {
	// 没有 CA 证书/密钥就没法（也不需要）建 CRL。
	if _, err := os.Stat(constants.ServerCACertPath); os.IsNotExist(err) {
		return nil
	}
	if _, err := os.Stat(constants.ServerCAKeyPath); os.IsNotExist(err) {
		return nil
	}

	// crl.cnf 必须在场（gencrl 要用），顺手把辅助文件刷进卷。
	if err := EnsureServerHelperFiles(); err != nil {
		return err
	}

	dbDir := constants.ServerCRLDBDir
	if err := os.MkdirAll(filepath.Join(dbDir, "newcerts"), 0755); err != nil {
		return fmt.Errorf("创建 CA 数据库目录失败: %v", err)
	}

	// index.txt：openssl ca 的「账本」，初始为空文件。
	indexPath := filepath.Join(dbDir, "index.txt")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		if err := os.WriteFile(indexPath, []byte{}, 0644); err != nil {
			return fmt.Errorf("创建 index.txt 失败: %v", err)
		}
	}

	// index.txt.attr：unique_subject=no，允许同一 CN 多次出现（重建用户时不报冲突）。
	attrPath := filepath.Join(dbDir, "index.txt.attr")
	if _, err := os.Stat(attrPath); os.IsNotExist(err) {
		if err := os.WriteFile(attrPath, []byte("unique_subject = no\n"), 0644); err != nil {
			return fmt.Errorf("创建 index.txt.attr 失败: %v", err)
		}
	}

	// crlnumber：CRL 序号计数器，openssl gencrl 每次自增，初值 1000。
	crlNumberPath := filepath.Join(dbDir, "crlnumber")
	if _, err := os.Stat(crlNumberPath); os.IsNotExist(err) {
		if err := os.WriteFile(crlNumberPath, []byte("1000\n"), 0644); err != nil {
			return fmt.Errorf("创建 crlnumber 失败: %v", err)
		}
	}

	// serial：签发序号计数器。revoke/gencrl 用不到，但 crl.cnf 里 serial= 指向它，
	// 个别 openssl 版本会校验其存在；幂等补一个初值，避免「打不开 serial 文件」报错。
	serialPath := filepath.Join(dbDir, "serial")
	if _, err := os.Stat(serialPath); os.IsNotExist(err) {
		if err := os.WriteFile(serialPath, []byte("1000\n"), 0644); err != nil {
			return fmt.Errorf("创建 serial 失败: %v", err)
		}
	}

	// crl.pem 不存在 → 生成初始空 CRL（防地雷：crl-verify 引用前先有有效文件）。
	if _, err := os.Stat(constants.ServerCRLPath); os.IsNotExist(err) {
		cmd := fmt.Sprintf("openssl ca -config %s -gencrl -out %s", constants.ServerCRLConfig, constants.ServerCRLPath)
		if err := utils.ExecCommand(cmd); err != nil {
			return fmt.Errorf("生成初始 CRL 失败: %v", err)
		}
		fmt.Printf("已生成初始空 CRL: %s\n", constants.ServerCRLPath)
	}
	return nil
}

// RevokeClientCert 吊销某用户的证书并重生成 CRL（删除用户时调用，永久生效）。
//
// 我们的客户端证书是用 `openssl x509 -req` 直接签的，没进 openssl ca 的账本(index.txt)，
// 所以 `openssl ca -revoke` 找不到条目会失败。这里先按证书的 serial / notAfter 往 index.txt
// 补一条 Valid(V) 记录，再 revoke 把它翻成 Revoked(R)，最后 gencrl 把 R 条目写进 crl.pem。
//
// crl-verify 每次新连接都重读 crl.pem，无需重启 OpenVPN 即生效。
func RevokeClientCert(username string) error {
	if !safeUsername.MatchString(username) {
		return fmt.Errorf("非法用户名: %q", username)
	}

	crtPath := filepath.Join(constants.ClientConfigDir, username+".crt")
	if _, err := os.Stat(crtPath); os.IsNotExist(err) {
		// 没有证书可吊销（可能从未生成或已删）——不算错误。
		fmt.Printf("用户 %s 无证书文件，跳过吊销\n", username)
		return nil
	}

	if err := EnsureCRLSetup(); err != nil {
		return fmt.Errorf("CRL 环境准备失败: %v", err)
	}

	indexPath := filepath.Join(constants.ServerCRLDBDir, "index.txt")
	// 单条 bash：取 serial / notAfter → 若 index 里没有该 serial 则补一条 V 记录 →
	// openssl ca -revoke 翻成 R → gencrl 重生成 CRL。serial 是长 hex，用作去重键够唯一。
	script := fmt.Sprintf(`set -e
CRT=%q
CFG=%q
INDEX=%q
CRLOUT=%q
SERIAL=$(openssl x509 -in "$CRT" -noout -serial | cut -d= -f2)
if ! grep -qiF "$SERIAL" "$INDEX"; then
  ENDRAW=$(openssl x509 -in "$CRT" -noout -enddate | cut -d= -f2)
  ENDFMT=$(date -u -d "$ENDRAW" +%%y%%m%%d%%H%%M%%SZ)
  printf 'V\t%%s\t\t%%s\tunknown\t/CN=%s\n' "$ENDFMT" "$SERIAL" >> "$INDEX"
fi
openssl ca -config "$CFG" -revoke "$CRT" -batch
openssl ca -config "$CFG" -gencrl -out "$CRLOUT"`,
		crtPath, constants.ServerCRLConfig, indexPath, constants.ServerCRLPath, username)

	if err := utils.ExecCommand(script); err != nil {
		return fmt.Errorf("吊销证书失败: %v", err)
	}
	fmt.Printf("已吊销用户 %s 的证书并更新 CRL\n", username)
	return nil
}
