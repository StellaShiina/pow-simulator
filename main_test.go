package main_test

import (
	"testing"
	"time"

	"github.com/stellashiina/pow-routine/config"
	"github.com/stellashiina/pow-routine/simulator"
)

func TestDefaultRun(t *testing.T) {
	config.NodeCount = 10
	config.Difficulty = 0.001
	config.Density = 0.2
	config.TargetHeight = 100
	config.Seed = 42

	sim := simulator.NewSimulator()
	sim.Run(10 * time.Second)
	t.Logf("Simulation completed: blocks=%d", func() int { _, c := sim.Bootnode.Stats(); return c }())
}

func TestParameterSweep(t *testing.T) {
	nodeCounts := []int{5, 10, 20}
	difficulties := []float64{0.005, 0.01, 0.02, 0.05}
	densities := []float64{0.3, 0.6, 1.0}

	for _, nc := range nodeCounts {
		for _, d := range difficulties {
			for _, den := range densities {
				config.NodeCount = nc
				config.Difficulty = d
				config.Density = den
				config.TargetHeight = 30
				config.Seed = 42

				sim := simulator.NewSimulator()
				sim.Run(8 * time.Second)

				_, blocks := sim.Bootnode.Stats()
				t.Logf("Nodes=%d, Diff=%.4f, Density=%.2f → Blocks=%d",
					nc, d, den, blocks)
			}
		}
	}
}
