package common

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

var dbPath string

func init() {
	// 数据库路径，可通过环境变量 LEVELDB_PATH 覆盖
	dp := os.Getenv("LEVELDB_PATH")
	if dp == "" {
		dp = "/var/lib/openvpn-manager"
	}
	dbPath = dp
}

// DBPath 返回 LevelDB 存储路径
func DBPath() string {
	return dbPath
}

// GetValue 获取 LevelDB 中的值
func GetValue(key string) (string, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()
	result, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// SetValue 设置leveldb值
func SetValue(key string, value string) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Put([]byte(key), []byte(value), nil)
}

// DelValue 删除值
func DelValue(key string) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Delete([]byte(key), nil)
}
