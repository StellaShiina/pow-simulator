// Package node implements the simulation node.
//
// Each node runs in a dedicated goroutine, maintaining its own mining loop and P2P connections.
// Node workflow:
//  1. Startup → Create TCP listener → Register with Bootnode → Wait for peer assignment
//  2. Receive neighbor list from Bootnode → Establish TCP connections with neighbors
//  3. Start mining goroutine: Continuously attempt to mine (based on difficulty probability)
//  4. Upon mining a block, broadcast it to all neighbors via TCP; relay new blocks upon receipt
//  5. Simultaneously listen for blocks from neighbors and update the local chain head
//
// Key design features:
//   - Uses ChaCha8-based pseudorandomness to simulate hashing instead of performing actual SHA256 calculations
//   - Incorporates time.Sleep in the mining loop to control the attempt rate, ensuring the block generation speed does not exceed the network propagation speed
//   - Propagates blocks between nodes via gossip relay; forwards new blocks to other neighbors (excluding the source)
//   - Uses a MsgHello handshake to ensure inbound connections can identify the dialing party
package node

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/stellashiina/pow-routine/config"
	"github.com/stellashiina/pow-routine/core"
	"github.com/stellashiina/pow-routine/p2p"
)

const (
	// Times to retry
	bootnodeRetry = 3
	peerDialRetry = 2
	// Sleep if fail to mine
	mineInterval = 500 * time.Microsecond
)

type Node struct {
	ID    int
	Tip   *core.Block
	tipMu sync.RWMutex

	// Only record block hash
	chain   map[[32]byte]struct{}
	chainMu sync.Mutex

	listener net.Listener
	peers    map[int]net.Conn
	peersMu  sync.Mutex

	bootConn net.Conn

	// controllers
	stopCh    chan struct{}
	stopped   bool
	stoppedMu sync.Mutex

	// statistic
	blocksMined int
}

func NewNode(id int) *Node {
	return &Node{
		ID:     id,
		Tip:    nil, // Init in simulator
		peers:  make(map[int]net.Conn),
		stopCh: make(chan struct{}),
		chain:  make(map[[32]byte]struct{}),
	}
}

// start node -> start listener -> connect to bootnode -> start miner
func (n *Node) Start(bootAddr string, genesis *core.Block) error {
	n.Tip = genesis

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("node %d listen: %w", n.ID, err)
	}
	n.listener = listener

	// start listen
	go n.listenLoop()

	// connect to bootnode
	if err := n.connectBootnode(bootAddr); err != nil {
		return fmt.Errorf("node %d connect bootnode: %w", n.ID, err)
	}

	// start mining loop
	go n.mineLoop()

	return nil
}

// set stopped -> stop listener -> close conn to peers -> disconnect to bootnode
func (n *Node) Stop() {
	n.stoppedMu.Lock()
	if n.stopped {
		n.stoppedMu.Unlock()
		return
	}
	n.stopped = true
	close(n.stopCh)
	n.stoppedMu.Unlock()

	if n.listener != nil {
		n.listener.Close()
	}
	n.peersMu.Lock()
	for _, conn := range n.peers {
		conn.Close()
	}
	n.peersMu.Unlock()
	if n.bootConn != nil {
		n.bootConn.Close()
	}
}

// ==================== NETWORK ====================

// connect to bootnode -> register -> wait for peers
func (n *Node) connectBootnode(bootAddr string) error {
	conn, err := p2p.DialWithRetry(bootAddr, bootnodeRetry)
	if err != nil {
		return err
	}
	n.bootConn = conn

	// register
	addr := n.listener.Addr().String()
	reg := p2p.RegisterPayload{NodeID: n.ID, Addr: addr}
	if err := p2p.SendMessage(conn, p2p.MsgRegisterPayload, reg); err != nil {
		return fmt.Errorf("register: %w", err)
	}
	slog.Debug("[Node] registered at bootnode", "NodeID", n.ID, "NodeAddr", addr)

	// request for peer list
	msg, err := p2p.ReadMsg(conn)
	if err != nil {
		return fmt.Errorf("read ack: %w", err)
	}
	if msg.Type != p2p.MsgPeerList {
		return fmt.Errorf("expected peer list ack, got %s", msg.Type)
	}

	// wait for peer list
	go n.bootListenerLoop(conn)

	return nil
}

func (n *Node) bootListenerLoop(conn net.Conn) {
	for {
		msg, err := p2p.ReadMsg(conn)
		if err != nil {
			select {
			case <-n.stopCh:
				return
			default:
				return
			}
		}
		switch msg.Type {
		case p2p.MsgPeerAssignment:
			n.handlePeerAssignment(msg.Data)
		case p2p.MsgBlock:
			// TODO
		}
	}
}

// when get peers list from bootnode, connect to all peers
func (n *Node) handlePeerAssignment(data json.RawMessage) {
	var assignment p2p.PeerList
	if err := json.Unmarshal(data, &assignment); err != nil {
		slog.Warn("[Node] bad peer assignment:", "NodeID", n.ID, "err", err)
		return
	}

	for _, peer := range assignment.Peers {
		if peer.NodeID == n.ID {
			continue
		}
		n.connectToPeer(peer.NodeID, peer.Addr)
	}

	n.peersMu.Lock()
	peerCount := len(n.peers)
	n.peersMu.Unlock()
	slog.Debug("[Node] peer assignment applied: peers connected", "NodeID", n.ID, "PeerCount", peerCount)
}

// dial a tcp connect to neighbour and try to handshake
func (n *Node) connectToPeer(peerID int, addr string) bool {
	n.peersMu.Lock()
	if _, exists := n.peers[peerID]; exists {
		n.peersMu.Unlock()
		return false
	}
	n.peersMu.Unlock()

	conn, err := p2p.DialWithRetry(addr, peerDialRetry)
	if err != nil {
		return false
	}

	// send handshake msg to recognize each other
	if err := p2p.SendMessage(conn, p2p.MsgHello, p2p.HelloPayload{NodeID: n.ID}); err != nil {
		conn.Close()
		return false
	}

	n.peersMu.Lock()
	n.peers[peerID] = conn
	n.peersMu.Unlock()

	go n.peerReadLoop(peerID, conn)

	return true
}

func (n *Node) listenLoop() {
	for {
		conn, err := n.listener.Accept()
		if err != nil {
			select {
			case <-n.stopCh:
				return
			default:
				continue
			}
		}
		go n.handleInboundConn(conn)
	}
}

func (n *Node) handleInboundConn(conn net.Conn) {
	msg, err := p2p.ReadMsg(conn)
	if err != nil {
		conn.Close()
		return
	}
	if msg.Type != p2p.MsgHello {
		conn.Close()
		return
	}
	var hello p2p.HelloPayload
	if err := json.Unmarshal(msg.Data, &hello); err != nil {
		conn.Close()
		return
	}
	peerID := hello.NodeID

	n.peersMu.Lock()
	if _, exists := n.peers[peerID]; !exists {
		n.peers[peerID] = conn
	} else {
		conn.Close()
		n.peersMu.Unlock()
		return
	}
	n.peersMu.Unlock()

	go n.peerReadLoop(peerID, conn)
}

func (n *Node) peerReadLoop(peerID int, conn net.Conn) {
	defer func() {
		n.peersMu.Lock()
		delete(n.peers, peerID)
		n.peersMu.Unlock()
		conn.Close()
	}()

	for {
		msg, err := p2p.ReadMsg(conn)
		if err != nil {
			return
		}
		switch msg.Type {
		case p2p.MsgBlock:
			var bp p2p.BlockPayload
			if err := json.Unmarshal(msg.Data, &bp); err != nil {
				continue
			}
			n.receiveBlock(&bp, peerID)
		}
	}
}

// Process received blocks: skip or continue propagating; update if the height is greater.
func (n *Node) receiveBlock(bp *p2p.BlockPayload, fromPeerID int) {
	n.chainMu.Lock()
	if _, seen := n.chain[bp.Hash]; seen {
		n.chainMu.Unlock()
		return
	}
	n.chain[bp.Hash] = struct{}{}
	n.chainMu.Unlock()

	// update tip
	n.tipMu.Lock()
	if bp.Height > n.Tip.Height {
		n.Tip = &core.Block{
			Height:  bp.Height,
			MinerID: bp.MinerID,
			Hash:    bp.Hash,
			PreHash: bp.PreHash,
		}
	}
	n.tipMu.Unlock()

	// Relay to all neighbors (except the source) to implement gossip flooding.
	n.peersMu.Lock()
	for peerID, conn := range n.peers {
		if peerID == fromPeerID {
			continue
		}
		if err := p2p.SendMessage(conn, p2p.MsgBlock, bp); err != nil {
			slog.Warn("[Node] relay block to peer failed", "NodeID", n.ID, "PeerID", peerID, "err", err)
		}
	}
	n.peersMu.Unlock()
}

// ==================== MINING ====================

func (n *Node) mineLoop() {
	for {
		select {
		case <-n.stopCh:
			return
		default:
			if n.tryMine() {
				// TODO
				// May do something, now msg is sent in tryMine()
			}
			time.Sleep(mineInterval)
		}
	}
}

// If success, update self tip -> update self chain -> braodcast to peers -> report to bootnode
func (n *Node) tryMine() bool {
	if config.RandSrc.Float64() >= config.Difficulty {
		return false
	}

	// Success and produce now block
	n.tipMu.Lock()
	currentTip := n.Tip
	blk := core.NewBlock(currentTip, n.ID)
	n.Tip = blk
	n.blocksMined++
	n.tipMu.Unlock()

	bp := p2p.BlockPayload{
		NodeID:  n.ID,
		Height:  blk.Height,
		MinerID: blk.MinerID,
		Hash:    blk.Hash,
		PreHash: blk.PreHash,
	}

	n.chainMu.Lock()
	n.chain[bp.Hash] = struct{}{}
	n.chainMu.Unlock()

	n.broadcastBlock(&bp)

	n.reportToBootnode(&bp)

	return true
}

func (n *Node) broadcastBlock(bp *p2p.BlockPayload) {
	n.peersMu.Lock()
	defer n.peersMu.Unlock()

	for peerID, conn := range n.peers {
		if err := p2p.SendMessage(conn, p2p.MsgBlock, bp); err != nil {
			slog.Warn("[Node] send block to peer failed", "NodeID", n.ID, "PeerID", peerID, "err", err)
		}
	}
}

func (n *Node) reportToBootnode(bp *p2p.BlockPayload) {
	if n.bootConn == nil {
		return
	}
	if err := p2p.SendMessage(n.bootConn, p2p.MsgBlock, bp); err != nil {
		// connection to bootnode may be closed
	}
}
