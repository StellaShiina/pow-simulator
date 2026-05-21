package core

import (
	"github.com/stellashiina/pow-simulator/config"
)

// A concise representation that does not distinguish between the head and the body, directly recording key simulation parameters.
type Block struct {
	Tick    uint64 // When(in which tick) the block is mined
	Height  int
	MinerID int
	Hash    [32]byte // Since it's not a real hash calculation, I use its own hash field to simulate it.
	PreHash [32]byte
}

func NewGenesis() *Block {
	var hash [32]byte
	config.RandHashSrc.Read(hash[:])
	return &Block{
		Tick:    0,
		Height:  0,
		MinerID: -1,
		Hash:    hash,
		PreHash: [32]byte{},
	}
}

func NewBlock(parent *Block, minerID int, tick uint64) *Block {
	var hash [32]byte
	config.RandHashSrc.Read(hash[:])
	return &Block{
		Tick:    tick,
		Height:  parent.Height + 1,
		MinerID: minerID,
		Hash:    hash,
		PreHash: parent.Hash,
	}
}
