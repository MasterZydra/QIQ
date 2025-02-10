package stats

import (
	"GoPHP/cmd/goPHP/config"
	"fmt"
	"runtime"
	"time"
)

func getUsedMemory() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// runtime.GC()
	return m.Alloc
}

type Stats struct {
	start     time.Time
	stop      time.Time
	beforeMem uint64
	afterMem  uint64
}

func NewStats() *Stats {
	return &Stats{}
}

func (stats *Stats) Start() {
	stats.beforeMem = getUsedMemory()
	stats.start = time.Now()
}

func (stats *Stats) Stop() {
	stats.afterMem = getUsedMemory()
	stats.stop = time.Now()
}

func (stats *Stats) UsedMemory() uint64 {
	return stats.afterMem - stats.beforeMem
}

func (stats *Stats) ElapsedTime() time.Duration {
	return stats.stop.Sub(stats.start)
}

func (stats *Stats) PrintStats() {
	fmt.Printf(" Time: %12s     Memory: %5s\n",
		fmt.Sprintf("%5.3f ms",
			float64(stats.ElapsedTime().Nanoseconds())/1000000,
		),
		fmt.Sprintf("%5d KB",
			stats.UsedMemory()/1024,
		),
	)
}

func (stats *Stats) PrintStatsWithPrefix(prefix string) {
	fmt.Printf("%-15s ", prefix)
	stats.PrintStats()
}

func Start() *Stats {
	if !config.ShowStats {
		return nil
	}

	stat := NewStats()
	stat.Start()
	return stat
}

func StopAndPrint(stats *Stats, prefix string) {
	if !config.ShowStats {
		return
	}

	stats.Stop()
	stats.PrintStatsWithPrefix(prefix)
}
