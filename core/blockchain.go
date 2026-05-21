package core

type Blockchain struct {
	Genesis *Block
	Blocks  map[[32]byte]*Block
}

// Instead of using regular slices to record potential blockchain forks
// Map is used in the simulation to track the block using hash values.
func NewBlockchain(genesis *Block) *Blockchain {
	blocks := make(map[[32]byte]*Block)
	blocks[genesis.Hash] = genesis
	return &Blockchain{
		Genesis: genesis,
		Blocks:  blocks,
	}
}
