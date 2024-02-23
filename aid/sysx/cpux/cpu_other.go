//go:build !linux

package cpux

func RefreshCpu() uint64 {
	return 0
}
