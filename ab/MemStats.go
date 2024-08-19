package ab

import (
	"fmt"
	"runtime"
)

var prevHeapSize uint64

func MemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	heapSize := m.HeapAlloc
	heapMaxSize := m.HeapSys
	heapFreeSize := m.HeapIdle
	diff := int64(heapSize) - int64(prevHeapSize)
	prevHeapSize = heapSize

	fmt.Printf("Heap %d MaxSize %d Free %d Diff %d\n", heapSize, heapMaxSize, heapFreeSize, diff)
}
