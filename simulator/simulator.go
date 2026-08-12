// 1. init randsrc, genesis, chain
// 2. new and start bootnode
// 3. create nodes and starts
// 4. monitor if it reached target
// 5. collect data
package simulator

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/stellashiina/pow-routine/bootnode"
	"github.com/stellashiina/pow-routine/config"
	"github.com/stellashiina/pow-routine/core"
	"github.com/stellashiina/pow-routine/node"
)

const bootnodePort = 9000

type Simulator struct {
	Bootnode  *bootnode.Bootnode
	Nodes     []*node.Node
	Chain     *core.Blockchain
	startTime time.Time
}

func NewSimulator() *Simulator {
	config.Init()
	genesis := core.NewGenesis()
	return &Simulator{
		Bootnode: bootnode.NewBootnode(bootnodePort),
		Nodes:    make([]*node.Node, 0, config.NodeCount),
		Chain:    core.NewBlockchain(genesis),
	}
}

func (sim *Simulator) Run(duration time.Duration) {
	sim.startTime = time.Now()
	genesis := sim.Chain.Genesis

	// 1. start bootnode
	slog.Info("[Simulator] starting bootnode...")
	if err := sim.Bootnode.Start(); err != nil {
		slog.Error("[Simulator] start bootnode: ", "err", err)
	}

	select {
	case <-sim.Bootnode.Ready:
		break
	case <-time.After(100 * time.Microsecond):
		slog.Error("[Simulator] start bootnode: timeout")
		return
	}

	sim.Bootnode.SetGenesis(genesis.Hash)

	// 2. create nodes
	slog.Info("[Simulator] creating nodes...", "NodeCount", config.NodeCount)
	for i := 0; i < config.NodeCount; i++ {
		n := node.NewNode(i)
		sim.Nodes = append(sim.Nodes, n)
	}

	// 3. start all nodes
	bootAddr := fmt.Sprintf("127.0.0.1:%d", bootnodePort)
	slog.Info("[Simulator] starting all nodes...")
	for _, n := range sim.Nodes {
		if err := n.Start(bootAddr, genesis); err != nil {
			slog.Error("[Simulator] node start failed:", "NodeID", n.ID, "err", err)
		}
		time.Sleep(20 * time.Millisecond)
	}

	// 4. wait for bootnode to asign peers
	slog.Info("[Simulator] waiting for peer assignment...")
	select {
	case <-sim.Bootnode.PeersAssigned:
		slog.Info("[Simulator] peer assignment complete")
	case <-time.After(10 * time.Second):
		slog.Warn("[Simulator] WARNING: peer assignment timeout, proceeding anyway")
	}

	// 5. simulate
	slog.Info("[Simulator] running for", "Duration", duration, "TargetHeight", config.TargetHeight)
	sim.waitForCompletion(duration)

	// 6. stop nodes
	slog.Info("[Simulator] stopping all nodes...")
	for _, n := range sim.Nodes {
		n.Stop()
	}

	// 7. write report
	slog.Info("[Simulator] generating chain report...")
	if err := sim.Bootnode.WriteReport(); err != nil {
		slog.Error("[Simulator] write report:", "err", err)
	}
	if err := sim.Bootnode.WriteJSONReport(); err != nil {
		slog.Error("[Simulator] write JSON report:", "err", err)
	}

	sim.Bootnode.Stop()

	// 8. print basic data
	sim.printStats()
}

// Wait for the simulation termination condition: target altitude reached first or timeout.
func (sim *Simulator) waitForCompletion(duration time.Duration) {
	select {
	case <-sim.Bootnode.Done:
		slog.Info("[Simulator] target height reached, stopping early.", "TargetHeight", config.TargetHeight)
	case <-time.After(duration):
		slog.Warn("[Simulator] timeout reached.", "Duration", duration)
	}
}

// print basic data
func (sim *Simulator) printStats() {
	elapsed := time.Since(sim.startTime).Seconds()
	nodeCount, blockCount := sim.Bootnode.Stats()

	analysis, err := sim.Bootnode.Analyze()

	fmt.Println("\n============================================")
	fmt.Println("         SIMULATION RESULTS")
	fmt.Println("============================================")
	fmt.Printf("  Nodes:          %d\n", nodeCount)
	fmt.Printf("  Difficulty:     %.6f\n", config.Difficulty)
	fmt.Printf("  Density:        %.2f\n", config.Density)
	fmt.Printf("  Target Height:  %d\n", config.TargetHeight)
	fmt.Printf("  Run Time:       %.2f seconds\n", elapsed)
	fmt.Printf("  Blocks Mined:   %d\n", blockCount)
	if blockCount > 0 {
		fmt.Printf("  Block Rate:     %.2f blocks/sec\n", float64(blockCount)/elapsed)
	}
	if err == nil {
		fmt.Printf("  Longest Chain:  %d blocks\n", analysis.LongestHeight)
		fmt.Printf("  Total Forks:    %d\n", analysis.TotalForks)
		fmt.Printf("  Active Tips:    %d\n", analysis.ActiveTips)
	}
	fmt.Println("============================================")
}
