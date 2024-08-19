package ab

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer()
	fmt.Printf("Elapsed Time in Milliseconds: %d\n", timer.ElapsedTimeMillis())
	fmt.Printf("Elapsed Time in Seconds: %.2f\n", timer.ElapsedTimeSecs())
	fmt.Printf("Elapsed Time in Minutes: %.2f\n", timer.ElapsedTimeMins())

	time.Sleep(2 * time.Second)

	fmt.Println("Restarting Timer...")
	fmt.Printf("Elapsed Time in Milliseconds after restart: %d\n", timer.ElapsedRestartMs())
}
