//go:build !linux

package dataFS

func getAvailableSize(_ string) int64 { return 0 }
