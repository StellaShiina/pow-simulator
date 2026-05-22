package simulator

import (
	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/core"
	"github.com/stellashiina/pow-simulator/network"
	"github.com/stellashiina/pow-simulator/node"
)

type ExpResult struct {
	Strategy    node.Strategy
	ForkSuccess bool
}

// The simulation uses a centrally broadcast network to simulate a P2P network.
// The advantage of doing this is that the simulation experiment is controllable.
// Point-to-point latency can be perfectly addressed by introducing latency at the network layer.
type Simulator struct {
	Network *network.Network
	Chain   *core.Blockchain
	Tips    []*core.Block // Store all the highest blocks
	Tick    uint64
}

func NewSimulator() *Simulator {
	genesis := core.NewGenesis()
	return &Simulator{
		Network: network.NewNetwork(config.NodesCount, genesis),
		Chain:   core.NewBlockchain(genesis),
		Tips:    []*core.Block{genesis},
	}
}

// Step 1 Miners try to mine
// ------ Mining occurs sequentially within a tick, but the propagation at the network layers closely simulates a real network.
// ------ I tried using a concurrent model for simulation, but the results were very poor. See the README for details.
// Step 2 Record new blocks and send them to network, to further introduce latency to packets in network model(see network.go).
// Step 3 Make network send packets(see network.go)
// Step 4 Miners update tip.
func (sim *Simulator) Run() *ExpResult {
	for {
		sim.Tick++
		for i, n := range sim.Network.Nodes {
			blk := n.Mine(sim.Tick)
			if blk != nil {
				// fmt.Printf("[Tick %d] New Block: Miner=%d, Height=%d\n", sim.Tick, blk.MinerID, blk.Height)
				if i == 0 && n.Strategy == node.Fork && n.ExpVar.FinishFork {
					if sim.Tips[0].Height <= 6 {
						return &ExpResult{
							Strategy:    node.Fork,
							ForkSuccess: true,
						}
					} else {
						return &ExpResult{
							Strategy:    node.Fork,
							ForkSuccess: false,
						}
					}
				}
				sim.Chain.Blocks[blk.Hash] = blk
				if blk.Height >= sim.Tips[0].Height {
					if blk.Height > sim.Tips[0].Height {
						sim.Tips = sim.Tips[:0]
					}
					sim.Tips = append(sim.Tips, blk)
				}
				sim.Network.ProcessPacket(sim.Tick, blk)
			}
		}
		sim.Network.SendPackets(sim.Tick)
		// Check if the selfish node has blocks to release
		if sim.Network.Nodes[0].Strategy == node.Selfish {
			for _, blk := range sim.Network.Nodes[0].ExpVar.Release {
				// fmt.Printf("[Tick %d] Get Block From Selfish Miner: Miner=%d, Height=%d\n", sim.Tick, blk.MinerID, blk.Height)
				sim.Chain.Blocks[blk.Hash] = blk
				if blk.Height >= sim.Tips[0].Height {
					if blk.Height > sim.Tips[0].Height {
						sim.Tips = sim.Tips[:0]
					}
					sim.Tips = append(sim.Tips, blk)
				}
				sim.Network.ProcessPacket(sim.Tick, blk)
				sim.Network.SendPackets(sim.Tick)
			}
			sim.Network.Nodes[0].ExpVar.Release = sim.Network.Nodes[0].ExpVar.Release[:0]
		}
		for _, n := range sim.Network.Nodes {
			n.UpdateTip()
		}
		if sim.Tips[0].Height == config.Target {
			// IF block speed is really high, there may be many blocks remain in privatechain!
			// for _, blk := range sim.Network.Nodes[0].ExpVar.PrivateChain {
			// 	if blk.Height > config.Target+1 {
			// 		break
			// 	}
			// 	fmt.Printf("[Tick %d] Get Block From Selfish Miner: Miner=%d, Height=%d\n", sim.Tick, blk.MinerID, blk.Height)
			// 	sim.Chain.Blocks[blk.Hash] = blk
			// 	if blk.Height >= sim.Tips[0].Height {
			// 		if blk.Height > sim.Tips[0].Height {
			// 			sim.Tips = sim.Tips[:0]
			// 		}
			// 		sim.Tips = append(sim.Tips, blk)
			// 	}
			// 	sim.Network.ProcessPacket(sim.Tick, blk)
			// 	sim.Network.SendPackets(sim.Tick)
			// }
			break
		}
	}
	return nil
}
