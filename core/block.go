package core

import (
	"github.com/stellashiina/pow-routine/config"
)

// A concise representation that does not distinguish between the head and the body, directly recording key simulation parameters.
type Block struct {
	Height  int
	MinerID int
	Hash    [32]byte
	PreHash [32]byte
}

func NewGenesis() *Block {
	var hash [32]byte
	config.HashRandSrc.Read(hash[:])
	return &Block{
		Height:  0,
		MinerID: -1,
		Hash:    hash,
		PreHash: [32]byte{},
	}
}

func NewBlock(parent *Block, minerID int) *Block {
	var hash [32]byte
	config.HashRandSrc.Read(hash[:])
	return &Block{
		Height:  parent.Height + 1,
		MinerID: minerID,
		Hash:    hash,
		PreHash: parent.Hash,
	}
}
