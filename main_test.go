package main_test

import (
	"testing"

	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/simulator"
)

func TestParameters(t *testing.T) {
	nodeCounts := []int{10, 50, 100, 200, 500, 1000}
	difficulties := []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.25, 0.3, 0.5}
	for _, nc := range nodeCounts {
		for _, d := range difficulties {
			if float64(nc)*d > 10 {
				continue
			}
			config.NodesCount = nc
			config.Difficulty = d
			sim := simulator.NewSimulator()
			sim.Run()
			t.Logf("Node Count: %d, Difficulty: %f. Block Speed %f ticks/block\n", nc, d, float64(sim.Tick)/float64(config.Target))
		}
	}
}
