package main

import (
	"fmt"
	"time"

	"github.com/stellashiina/pow-routine/config"
	"github.com/stellashiina/pow-routine/simulator"
)

func main() {
	fmt.Println("=== PoW Routine-Based P2P Simulator ===")
	fmt.Printf("Parameters: nodes=%d, difficulty=%.6f, density=%.2f, targetHeight=%d\n",
		config.NodeCount, config.Difficulty, config.Density, config.TargetHeight)

	sim := simulator.NewSimulator()
	sim.Run(30 * time.Second) // Run with timeout
}
