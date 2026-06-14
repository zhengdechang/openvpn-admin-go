//go:build linux

package controller

import (
	"math"
	"syscall"
)

// diskUsage 通过 statfs 获取磁盘使用情况（Linux）。
func diskUsage(path string) DiskInfo {
	di := DiskInfo{Path: path}
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return di
	}
	bsize := uint64(st.Bsize)
	di.Total = st.Blocks * bsize
	di.Free = st.Bavail * bsize
	if st.Blocks >= st.Bfree {
		di.Used = (st.Blocks - st.Bfree) * bsize
	}
	if di.Total > 0 {
		di.UsedPercent = math.Round(float64(di.Used)/float64(di.Total)*1000) / 10
	}
	return di
}
