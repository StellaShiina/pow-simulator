package exp

import (
	"math/rand/v2"
	"time"

	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/core"
	"github.com/stellashiina/pow-simulator/node"
	"github.com/stellashiina/pow-simulator/simulator"
)

type SelfishMiningResult struct {
	Experiment    string    `json:"experiment"`
	MaliciousRate float64   `json:"malicious_rate"`
	Difficulty    float64   `json:"difficulty"`
	NodeCount     int       `json:"node_count"`
	Rounds        int       `json:"rounds"`
	IncomeBlocks  int       `json:"income_blocks"`
	Profitability float64   `json:"profitability"`
	RandomSeed    [2]uint64 `json:"random_seed"`
}

func Selfish(maliciousRate, difficulty float64, nodeCount int, randSeed *[2]uint64) SelfishMiningResult {
	round := 1000
	config.Difficulty = difficulty
	config.NodesCount = nodeCount
	var genRandSrc *rand.Rand
	if randSeed != nil {
		genRandSrc = rand.New(rand.NewPCG(randSeed[0], randSeed[1]))
	} else {
		genRandSrc = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()+1)))
	}
	selfishIncome := 0
	for range round {
		config.RandSrc = rand.New(rand.NewPCG(genRandSrc.Uint64(), genRandSrc.Uint64()))
		sim := simulator.NewSimulator()
		maliciousNode := sim.Network.Nodes[0]
		maliciousNode.HashPower = CalculateNodePower(maliciousRate)
		maliciousNode.Strategy = node.Selfish
		maliciousNode.ExpVar = &node.ExpVar{
			PrivateChain: make([]*core.Block, 0),
			Release:      make([]*core.Block, 0),
		}
		sim.Run()
		blks := sim.Tips
		sum := 0
		for _, blk := range blks {
			// fmt.Printf("[Chain %d]\n", i)
			for blk != nil {
				// fmt.Printf("[TICK %d] hash: %x; Miner: %d; Height: %d; PreHash: %x\n", blk.Tick, blk.Hash, blk.MinerID, blk.Height, blk.PreHash)
				if blk.MinerID == 0 {
					sum++
				}
				blk = sim.Chain.Blocks[blk.PreHash]
			}
		}
		selfishIncome += sum / len(sim.Tips)
	}
	result := SelfishMiningResult{
		Experiment:    "selfish_mining",
		MaliciousRate: maliciousRate,
		Difficulty:    difficulty,
		NodeCount:     nodeCount,
		Rounds:        round,
		IncomeBlocks:  selfishIncome,
		Profitability: float64(selfishIncome) / float64(round*config.Target),
	}
	if randSeed != nil {
		result.RandomSeed = *randSeed
	}
	return result
}
