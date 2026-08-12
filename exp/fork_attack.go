package exp

import (
	"math/rand/v2"
	"time"

	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/node"
	"github.com/stellashiina/pow-simulator/simulator"
)

type ForkAttackResult struct {
	Experiment    string    `json:"experiment"`
	MaliciousRate float64   `json:"malicious_rate"`
	Difficulty    float64   `json:"difficulty"`
	NodeCount     int       `json:"node_count"`
	Rounds        int       `json:"rounds"`
	Successes     int       `json:"successes"`
	SuccessRate   float64   `json:"success_rate"`
	RandomSeed    [2]uint64 `json:"random_seed"`
}

func ForkAttack(maliciousRate, difficulty float64, nodeCount int, randSeed *[2]uint64) ForkAttackResult {
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
	result := ForkAttackResult{
		Experiment:    "fork_attack",
		MaliciousRate: maliciousRate,
		Difficulty:    difficulty,
		NodeCount:     nodeCount,
		Rounds:        round,
		Successes:     success,
		SuccessRate:   float64(success) / float64(round),
	}
	if randSeed != nil {
		result.RandomSeed = *randSeed
	}
	return result
}

func CalculateNodePower(rate float64) float64 {
	return float64(config.NodesCount-1) / (1 - rate) * rate
}
