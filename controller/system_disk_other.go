//go:build !linux

package controller

// diskUsage 非 Linux 平台占位实现（交叉编译用），返回空磁盘信息。
func diskUsage(path string) DiskInfo {
	return DiskInfo{Path: path}
}
