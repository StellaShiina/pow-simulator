package exp

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/node"
	"github.com/stellashiina/pow-simulator/simulator"
)

func ForkAttack(maliciousRate, difficulty float64, nodeCount int, randSeed *[2]uint64) {
	round, success := 10000, 0
	config.Difficulty = difficulty
	config.NodesCount = nodeCount
	var genRandSrc *rand.Rand
	if randSeed != nil {
		genRandSrc = rand.New(rand.NewPCG(randSeed[0], randSeed[1]))
	} else {
		genRandSrc = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()+1)))
	}
	for range round {
		config.RandSrc = rand.New(rand.NewPCG(genRandSrc.Uint64(), genRandSrc.Uint64()))
		sim := simulator.NewSimulator()
		maliciousNode := sim.Network.Nodes[0]
		maliciousNode.HashPower = CalculateNodePower(maliciousRate)
		maliciousNode.Strategy = node.Fork
		maliciousNode.ExpVar = &node.ExpVar{}
		res := sim.Run()
		if res.ForkSuccess {
			success++
		}
	}
	fmt.Printf("[MaliciousRate %f] Fork success rate %f\n", maliciousRate, float64(success)/float64(round))
}

func CalculateNodePower(rate float64) float64 {
	return float64(config.NodesCount-1) / (1 - rate) * rate
}
