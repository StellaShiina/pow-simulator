package node_test

import (
	"testing"

	"github.com/stellashiina/pow-simulator/core"
	"github.com/stellashiina/pow-simulator/node"
)

func TestNode(t *testing.T) {
	genesis := core.NewGenesis()
	n := node.NewNode(0, genesis)
	// for {
	// 	if blk := n.Mine(1); blk != nil {
	// 		if n.KnownBlocks.Blocks[blk.Hash] != blk {
	// 			t.Error("KnownBlocks incorrect")
	// 			t.Fail()
	// 		}
	// 		break
	// 	}
	// }
	oldblk, newBlk, newBlk2 := core.NewBlock(genesis, -1, 2), core.NewBlock(n.Tip, n.ID, 2), core.NewBlock(n.Tip, n.ID, 2)
	n.Receive(genesis)
	n.Receive(oldblk)
	n.Receive(newBlk2)
	n.Receive(newBlk)
	n.UpdateTip()
	if n.Tip == genesis {
		t.Error("Wrong Tip")
		t.Fail()
	}
	if n.Tip != newBlk && n.Tip != newBlk2 {
		t.Error("Wrong Tip")
		t.Fail()
	}
}
