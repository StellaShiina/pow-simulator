package main

import (
	"fmt"

	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/simulator"
)

func main() {
	// var level slog.LevelVar
	// level.Set(slog.LevelWarn)

	// logger := slog.New(
	// 	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 		Level: &level,
	// 	}),
	// )
	// slog.SetDefault(logger)
	// sim := simulator.NewSimulator(2)
	fmt.Println("Pow-simulator")
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
			fmt.Printf("Node Count: %d, Difficulty: %f. Block Speed %f ticks/block\n", nc, d, float64(sim.Tick)/float64(config.Target))
		}
	}
	// sim := simulator.NewSimulator(config.NodesCount)
	// sim.Run(config.RoundCount)
	// for k, v := range sim.Chain.Blocks {
	// 	fmt.Printf("hash: %x; Miner: %d; Height: %d; PreHash: %x\n", k, v.MinerID, v.Height, v.PreHash)
	// }
	// blks := sim.Tips
	// for i, blk := range blks {
	// 	fmt.Printf("[Chain %d]\n", i)
	// 	for blk != nil {
	// 		fmt.Printf("[TICK %d] hash: %x; Miner: %d; Height: %d; PreHash: %x\n", blk.Tick, blk.Hash, blk.MinerID, blk.Height, blk.PreHash)
	// 		blk = sim.Chain.Blocks[blk.PreHash]
	// 	}
	// 	fmt.Println()
	// }
	// exp.ForkAttack()
}
