package core

type Blockchain struct {
	Genesis *Block
	Blocks  map[[32]byte]*Block
}

func NewBlockchain(genesis *Block) *Blockchain {
	blocks := make(map[[32]byte]*Block)
	blocks[genesis.Hash] = genesis
	return &Blockchain{
		Genesis: genesis,
		Blocks:  blocks,
	}
}

func (bc *Blockchain) AddBlock(blk *Block) {
	bc.Blocks[blk.Hash] = blk
}
