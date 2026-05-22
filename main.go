package main

import (
	"fmt"

	"github.com/stellashiina/pow-simulator/simulator"
)

func main() {
	fmt.Println("Pow-simulator")
	sim := simulator.NewSimulator()
	sim.Run()
	chain := sim.Chain.Blocks
	for i, tip := range sim.Tips {
		fmt.Printf("[CHAIN %d]\n", i)
		blk := tip
		for blk != nil {
			fmt.Printf("[TICK %d] | MINER %d | HEIGHT %d | HASH %x... | PRE %x... |\n", blk.Tick, blk.MinerID, blk.Height, blk.Hash[:8], blk.PreHash[:8])
			blk = chain[blk.PreHash]
		}
	}
}
