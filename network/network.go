package network

import (
	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/core"
	"github.com/stellashiina/pow-simulator/node"
)

// Packet
type Packet struct {
	ToID  int
	block *core.Block
}

// Use time buckets to simulate network latency, measured in ticks.
// Point-to-point random delay simulation was implemented with the definition of data packets.
type Network struct {
	Nodes   []*node.Node
	Packets map[uint64][]*Packet // Time(tick) Bucket. In round <tick>, network will send Packets[<tick>]
}

func NewNetwork(nodeNum int, genesis *core.Block) *Network {
	nodes := make([]*node.Node, nodeNum)
	for i := range nodeNum {
		nodes[i] = node.NewNode(i, genesis)
	}
	return &Network{
		Nodes:   nodes,
		Packets: make(map[uint64][]*Packet),
	}
}

// When new block come to network, random delay to each node is set and store in time buckets
func (nw *Network) ProcessPacket(tick uint64, blk *core.Block) {
	for _, n := range nw.Nodes {
		if blk.MinerID != n.ID {
			ramdomDelay := uint64(config.RandSrc.IntN(config.DelayRange))
			nw.Packets[tick+ramdomDelay] = append(nw.Packets[tick+ramdomDelay], &Packet{
				ToID:  n.ID,
				block: blk,
			})
		}
	}
}

func (nw *Network) SendPackets(tick uint64) {
	for _, p := range nw.Packets[tick] {
		nw.Nodes[p.ToID].Receive(p.block)
	}
	delete(nw.Packets, tick)
}
