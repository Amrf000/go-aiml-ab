package ab

import (
	"runtime"
)

// TotalMemory returns the total amount of memory obtained from the system.
func TotalMemory() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Sys
}

// MaxMemory returns an approximation of the maximum amount of memory that can be allocated by the system.
func MaxMemory() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Sys // Go does not have a direct equivalent of Java's maxMemory
}

// FreeMemory returns the amount of heap memory currently unused.
func FreeMemory() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Frees
}
