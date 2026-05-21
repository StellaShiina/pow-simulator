package config

import (
	"math/rand/v2"
)

var NodesCount int = 100

var RoundCount uint64 = 10000

var Difficulty float64 = 0.001

var Target int = 1000

var DelayRange int = 3

var RandSrc = rand.New(rand.NewPCG(42, 42))

var RandHashSrc = rand.NewChaCha8([32]byte{4, 2})
