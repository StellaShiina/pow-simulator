package exp_test

import (
	"testing"

	"github.com/stellashiina/pow-simulator/exp"
)

var randSeed2 = &[2]uint64{7, 21}

func TestSelfishVar1(t *testing.T) {
	difficulty := 0.001
	nodeCount := 100
	t.Logf("Current parameters, Difficulty: %f, Node Count: %d\n", difficulty, nodeCount)
	exp.Selfish(0.4, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.3, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.2, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.1, difficulty, nodeCount, randSeed2)
}

func TestSelfishVar2(t *testing.T) {
	difficulty := 0.001
	nodeCount := 50
	t.Logf("Current parameters, Difficulty: %f, Node Count: %d\n", difficulty, nodeCount)
	exp.Selfish(0.4, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.3, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.2, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.1, difficulty, nodeCount, randSeed2)
}

func TestSelfishVar3(t *testing.T) {
	difficulty := 0.0005
	nodeCount := 100
	t.Logf("Current parameters, Difficulty: %f, Node Count: %d\n", difficulty, nodeCount)
	exp.Selfish(0.4, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.3, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.2, difficulty, nodeCount, randSeed2)
	exp.Selfish(0.1, difficulty, nodeCount, randSeed2)
}
