package node

import (
	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/core"
)

// Each round(tick), a node will try to mine, if success, send to network and update its tip immediately.
// After mining, receive new blocks messages from network.
// Randomly pick the highest block as self tip.
// Note that node won't receive a block mined by itself(see simulator), if current tip is one of the highest block, the miner will keep mining based on its own block.

type Strategy int

const (
	Honest Strategy = iota
	Fork
	Selfish
)

type ExpVar struct {
	FinishFork   bool
	PrivateChain []*core.Block // In selfish mining
	Release      []*core.Block // Release block in selfish mining
}

// To reduce memory consumption, blockchain information is not enabled at the node level; it only exists in the simulation center.
type Node struct {
	ID int
	// KnownBlocks *core.Blockchain // Store a blockchain map
	Tip       *core.Block
	HashPower float64       // Simulate machine power，1 by default
	Inbox     []*core.Block // Receive new block from network
	Strategy  Strategy
	ExpVar    *ExpVar
}

func NewNode(id int, genesis *core.Block) *Node {
	return &Node{
		ID:        id,
		Tip:       genesis,
		HashPower: 1,
		// KnownBlocks: core.NewBlockchain(genesis),
	}
}

func (n *Node) Mine(tick uint64) *core.Block {
	probability := config.Difficulty * n.HashPower
	if config.RandSrc.Float64() < probability {
		blk := core.NewBlock(n.Tip, n.ID, tick)
		// n.KnownBlocks.Blocks[blk.Hash] = blk
		n.Tip = blk
		switch n.Strategy {
		case Fork:
			if blk.Height == 7 {
				n.ExpVar.FinishFork = true
				return blk
			} else {
				return nil
			}
		case Selfish:
			n.ExpVar.PrivateChain = append(n.ExpVar.PrivateChain, blk)
			return nil
		}
		return blk
	}
	return nil
}

func (n *Node) Receive(blk *core.Block) {
	// n.KnownBlocks.Blocks[blk.Hash] = blk
	// If inbox is empty, compare with current tip, if higher, append to inbox
	// This ensures that miners will prioritize continuing on their own blocks!!
	// If inbox isn't empty, always make sure only the highest blocks can be appended to inbox.
	if len(n.Inbox) == 0 && blk.Height > n.Tip.Height {
		n.Inbox = append(n.Inbox, blk)
	} else if len(n.Inbox) != 0 {
		if blk.Height >= n.Inbox[0].Height {
			if blk.Height > n.Inbox[0].Height {
				n.Inbox = n.Inbox[:0]
			}
			n.Inbox = append(n.Inbox, blk)
		}
	}
	if n.Strategy == Selfish {
		for len(n.ExpVar.PrivateChain) > 0 && n.ExpVar.PrivateChain[0].Height <= blk.Height {
			n.ExpVar.Release = append(n.ExpVar.Release, n.ExpVar.PrivateChain[0])
			n.ExpVar.PrivateChain = n.ExpVar.PrivateChain[1:]
		}
	}
}

// Randomly pick in inbox(if not empty).
// This simulates the out-of-order arrival of information from different nodes in a real network.
func (n *Node) UpdateTip() {
	switch n.Strategy {
	case Honest:
		if len(n.Inbox) != 0 {
			n.Tip = n.Inbox[config.RandSrc.IntN(len(n.Inbox))]
		}
		n.Inbox = n.Inbox[:0]
	case Fork:
		return
	case Selfish:
		if len(n.Inbox) != 0 {
			n.Tip = n.Inbox[config.RandSrc.IntN(len(n.Inbox))]
		}
		n.Inbox = n.Inbox[:0]
	}
}
