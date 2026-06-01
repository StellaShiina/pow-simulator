// Encode using JSON format and transmit via TCP.
// [4-byte big-endian length][JSON payload]。
package p2p

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	// Node register itself to bootnode
	MsgRegisterPayload = "register_payload"
	// Bootnode distribute peer list to node
	MsgPeerList = "peer_list"

	MsgBlock = "block"

	// Bootnode proactively push newly added neighbors to registered nodes.
	MsgPeerUpdate = "peer_update"

	// Bootnode Once registration is complete, distribute the final peer assignments to the respective nodes.
	MsgPeerAssignment = "peer_assignment"

	// A handshake message sent after the TCP connection between nodes is established, identifying the dialer.
	MsgHello = "hello"

	// Request chain information from neighbors (currently unused; reserved for future expansion).
	MsgChainRequest = "chain_req"

	// Response chain information (currently unused; reserved for future expansion).
	MsgChainResponse = "chain_resp"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type RegisterPayload struct {
	NodeID int    `json:"node_id"`
	Addr   string `json:"addr"`
}

// nodes handshake between each other
type HelloPayload struct {
	NodeID int `json:"node_id"`
}

// peer
type PeerEntry struct {
	NodeID int    `json:"node_id"`
	Addr   string `json:"addr"`
}

type PeerList struct {
	Peers []PeerEntry `json:"peers"`
}

type BlockPayload struct {
	NodeID  int      `json:"node_id"`
	Height  int      `json:"height"`
	MinerID int      `json:"miner_id"`
	Hash    [32]byte `json:"hash"`
	PreHash [32]byte `json:"pre_hash"`
}

// Lightweight block representation, used for reconstruction from network messages.
type BlockRepr struct {
	Height  int
	MinerID int
	Hash    [32]byte
	PreHash [32]byte
	Arrived time.Time `json:"-"`
}

func (b *BlockRepr) HashString() string {
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x%02x%02x",
		b.Hash[0], b.Hash[1], b.Hash[2], b.Hash[3],
		b.Hash[4], b.Hash[5], b.Hash[6], b.Hash[7])
}

func (b *BlockRepr) PreHashString() string {
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x%02x%02x",
		b.PreHash[0], b.PreHash[1], b.PreHash[2], b.PreHash[3],
		b.PreHash[4], b.PreHash[5], b.PreHash[6], b.PreHash[7])
}
