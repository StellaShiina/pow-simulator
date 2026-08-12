// Package report 负责生成仿真链结构报告和网络连接拓扑报告。
//
// 从 Bootnode 收集的区块 DAG 数据中分析链结构（分叉、最长链、矿工统计），
// 并将结果格式化为可读的文本报告写入文件。
package report

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/stellashiina/pow-routine/config"
	"github.com/stellashiina/pow-routine/p2p"
)

// ==================== 报告输入 ====================

// Input 聚合了生成报告所需的全部数据。
type Input struct {
	StartTime       time.Time
	Blocks          map[[32]byte]*p2p.BlockRepr
	PeerAssignments map[int][]int
}

// ==================== 轻量分析 ====================

// ChainSummary 链的轻量分析摘要，供控制台快速展示。
type ChainSummary struct {
	TotalBlocks   int
	LongestHeight int
	TotalForks    int
	ActiveTips    int
}

// Analyze 返回链的轻量分析摘要（不包含全量 DAG）。
func Analyze(blocks map[[32]byte]*p2p.BlockRepr) (*ChainSummary, error) {
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no blocks recorded")
	}

	children := make(map[[32]byte]int)
	heightCount := make(map[int]int)
	tips := 0
	longest := 0

	for _, b := range blocks {
		children[b.PreHash]++
		heightCount[b.Height]++
		if b.Height > longest {
			longest = b.Height
		}
	}

	for _, b := range blocks {
		if children[b.Hash] == 0 {
			tips++
		}
	}

	forks := 0
	for _, count := range heightCount {
		if count > 1 {
			forks++
		}
	}

	return &ChainSummary{
		TotalBlocks:   len(blocks),
		LongestHeight: longest,
		TotalForks:    forks,
		ActiveTips:    tips,
	}, nil
}

// ==================== DAG 分析 ====================

// chainNode 用于链重构的节点。
type chainNode struct {
	block    *p2p.BlockRepr
	children []*chainNode
	depth    int
}

// dagAnalysis 存储一次 DAG 分析的结果。
type dagAnalysis struct {
	genesis        *p2p.BlockRepr
	tips           []*p2p.BlockRepr
	blocks         map[[32]byte]*p2p.BlockRepr
	children       map[[32]byte][]*p2p.BlockRepr
	depths         map[[32]byte]int
	heightCount    map[int]int
	heightMiners   map[int][]int
	heights        []int
	minerCounts    map[int]int
	totalBlocks    int
	longestHeight  int
	tipPathLengths map[[32]byte]int
	totalForks     int
}

// analyzeDAG 扫描区块集合，构建 DAG 并返回分析结果。
func analyzeDAG(blocks map[[32]byte]*p2p.BlockRepr) (*dagAnalysis, error) {
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no blocks recorded")
	}

	blockMap := make(map[[32]byte]*p2p.BlockRepr)
	var genesis *p2p.BlockRepr
	for h, b := range blocks {
		blockMap[h] = b
		if b.Height == 0 {
			genesis = b
		}
	}
	if genesis == nil {
		return nil, fmt.Errorf("no genesis block found")
	}

	children := make(map[[32]byte][]*p2p.BlockRepr)
	for _, b := range blockMap {
		children[b.PreHash] = append(children[b.PreHash], b)
	}

	depths := computeDepths(genesis.Hash, children)

	var tips []*p2p.BlockRepr
	for _, b := range blockMap {
		if len(children[b.Hash]) == 0 {
			tips = append(tips, b)
		}
	}
	sort.Slice(tips, func(i, j int) bool {
		return tips[i].Height > tips[j].Height
	})

	heightCount := make(map[int]int)
	heightMiners := make(map[int][]int)
	for _, b := range blockMap {
		heightCount[b.Height]++
		heightMiners[b.Height] = append(heightMiners[b.Height], b.MinerID)
	}

	minerCounts := make(map[int]int)
	for _, b := range blockMap {
		if b.MinerID >= 0 {
			minerCounts[b.MinerID]++
		}
	}

	longestHeight := 0
	if len(tips) > 0 {
		longestHeight = tips[0].Height
	}

	tipPathLengths := make(map[[32]byte]int)
	for _, tip := range tips {
		length := pathLength(tip.Hash, blockMap)
		tipPathLengths[tip.Hash] = length
	}

	heights := make([]int, 0, len(heightCount))
	for h := range heightCount {
		heights = append(heights, h)
	}
	sort.Ints(heights)

	return &dagAnalysis{
		genesis:        genesis,
		tips:           tips,
		blocks:         blockMap,
		children:       children,
		depths:         depths,
		heightCount:    heightCount,
		heightMiners:   heightMiners,
		heights:        heights,
		minerCounts:    minerCounts,
		totalBlocks:    len(blockMap),
		longestHeight:  longestHeight,
		tipPathLengths: tipPathLengths,
		totalForks:     countForks(heightCount),
	}, nil
}

// computeDepths 递归计算每个节点的最大子树深度（自底向上）。
func computeDepths(hash [32]byte, children map[[32]byte][]*p2p.BlockRepr) map[[32]byte]int {
	depths := make(map[[32]byte]int)
	var dfs func(h [32]byte) int
	dfs = func(h [32]byte) int {
		if d, ok := depths[h]; ok {
			return d
		}
		maxChildDepth := 0
		for _, child := range children[h] {
			cd := dfs(child.Hash) + 1
			if cd > maxChildDepth {
				maxChildDepth = cd
			}
		}
		depths[h] = maxChildDepth
		return maxChildDepth
	}
	dfs(hash)
	return depths
}

// pathLength 计算从指定 block 到 genesis 的实际步数。
func pathLength(hash [32]byte, blocks map[[32]byte]*p2p.BlockRepr) int {
	length := 0
	current := hash
	for {
		b, ok := blocks[current]
		if !ok || b.Height == 0 {
			break
		}
		length++
		current = b.PreHash
		if length > 100000 {
			break
		}
	}
	return length
}

// countForks 统计分叉点数量（某个高度有 >1 个 block）。
func countForks(heightCount map[int]int) int {
	forks := 0
	for _, count := range heightCount {
		if count > 1 {
			forks++
		}
	}
	return forks
}

// traceChain 从 tip 回溯到 genesis，返回 block 列表（genesis → tip 顺序）。
func traceChain(tip *p2p.BlockRepr, blocks map[[32]byte]*p2p.BlockRepr) []*p2p.BlockRepr {
	var chain []*p2p.BlockRepr
	current := tip
	for current != nil {
		chain = append(chain, current)
		if current.Height == 0 {
			break
		}
		next, ok := blocks[current.PreHash]
		if !ok {
			break
		}
		current = next
	}
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain
}

// ==================== 报告分隔线常量 ====================

const (
	sep40 = "  ----------------------------------------\n"
	sep48 = "  ------------------------------------------------\n"
	sep58 = "  ----------------------------------------------------------\n"
	sep68 = "  --------------------------------------------------------------------\n"
)

// ==================== 报告文件前缀 ====================

const logFilePrefix = "result-routine"

// ==================== 链结构报告 ====================

// WriteChainReport 生成完整的链结构报告并写入文件。
func WriteChainReport(input Input) error {
	analysis, err := analyzeDAG(input.Blocks)
	if err != nil {
		return fmt.Errorf("chain analysis: %w", err)
	}

	now := time.Now()
	filename := filepath.Join("results", "log", fmt.Sprintf("%s-%s.log", logFilePrefix, now.Format("20060102-150405.000000000")))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create log file: %w", err)
	}
	defer f.Close()

	elapsed := time.Since(input.StartTime).Seconds()

	p := func(format string, args ...any) {
		fmt.Fprintf(f, format, args...)
	}

	p("============================================================\n")
	p("      PoW ROUTINE SIMULATION REPORT\n")
	p("============================================================\n\n")

	// ----- 参数 -----
	p("--- PARAMETERS ---\n")
	p("  Nodes:              %d\n", config.NodeCount)
	p("  Difficulty:         %.6f\n", config.Difficulty)
	p("  Density:            %.2f\n", config.Density)
	p("  Target Height:      %d\n", config.TargetHeight)
	p("  Seed:               %d\n", config.Seed)
	p("  Run Time:           %.2f seconds\n\n", elapsed)

	// ----- 概览 -----
	p("--- OVERVIEW ---\n")
	p("  Genesis Hash:       %s\n", analysis.genesis.HashString())
	blocksExclGenesis := analysis.totalBlocks - 1
	p("  Total Blocks:       %d (excl. genesis)\n", blocksExclGenesis)
	p("  Total Forks:        %d\n", analysis.totalForks)
	p("  Longest Chain:      %d blocks\n", analysis.longestHeight)
	p("  Active Tips:        %d\n", len(analysis.tips))
	if analysis.totalBlocks > 0 {
		p("  Average Block Rate: %.2f blocks/sec\n", float64(blocksExclGenesis)/elapsed)
	}
	p("\n")

	// ----- 矿工统计 -----
	p("--- PER-NODE MINING STATS ---\n")
	p("  %-6s | %-14s | %-6s\n", "NODE", "BLOCKS MINED", "SHARE")
	p(sep40)

	type minerEntry struct {
		id    int
		count int
	}
	var miners []minerEntry
	for id, count := range analysis.minerCounts {
		miners = append(miners, minerEntry{id, count})
	}
	sort.Slice(miners, func(i, j int) bool {
		return miners[i].count > miners[j].count
	})
	for _, m := range miners {
		share := float64(m.count) / float64(analysis.totalBlocks) * 100
		p("  %-6d | %-14d | %5.1f%%\n", m.id, m.count, share)
	}
	p("\n")

	// ----- 最长链 -----
	if len(analysis.tips) > 0 {
		longestTip := analysis.tips[0]
		longestChain := traceChain(longestTip, analysis.blocks)
		p("--- LONGEST CHAIN (%d blocks, tip by miner %d) ---\n",
			len(longestChain), longestTip.MinerID)
		p("  %-6s | %-6s | %-6s | %-18s | %-18s\n", "HEIGHT", "MINER", "SEQ", "HASH", "PREV HASH")
		p(sep68)
		for i, blk := range longestChain {
			seq := i
			if seq == 0 {
				p("  %-6d | %-6d | %-6d | %-18s | %-18s\n",
					blk.Height, blk.MinerID, seq, blk.HashString(), "(genesis)")
			} else {
				p("  %-6d | %-6d | %-6d | %-18s | %-18s\n",
					blk.Height, blk.MinerID, seq, blk.HashString(), blk.PreHashString())
			}
		}
		p("\n")

		// ----- 所有 Tips（分叉末端）-----
		if len(analysis.tips) > 1 {
			p("--- ALL TIPS (%d forks) ---\n", len(analysis.tips))
			p("  %-6s | %-6s | %-6s | %-6s | %-18s\n", "TIP#", "HEIGHT", "MINER", "LENGTH", "HASH")
			p(sep58)
			for i, tip := range analysis.tips {
				pathLen := analysis.tipPathLengths[tip.Hash]
				if i == 0 {
					p("  %-6d | %-6d | %-6d | %-6d | %-18s  ← longest\n",
						i, tip.Height, tip.MinerID, pathLen, tip.HashString())
				} else {
					p("  %-6d | %-6d | %-6d | %-6d | %-18s\n",
						i, tip.Height, tip.MinerID, pathLen, tip.HashString())
				}
			}
			p("\n")
		}
	}

	// ----- Fork 分析（按高度分组）-----
	p("--- FORK ANALYSIS ---\n")
	p("  %-6s | %-6s | %-18s\n", "HEIGHT", "BLOCKS", "MINERS")
	p(sep48)

	forkCount := 0
	for _, h := range analysis.heights {
		count := analysis.heightCount[h]
		miners := analysis.heightMiners[h]
		uniqueMiners := make(map[int]bool)
		for _, m := range miners {
			uniqueMiners[m] = true
		}
		minerList := make([]int, 0, len(uniqueMiners))
		for m := range uniqueMiners {
			minerList = append(minerList, m)
		}
		sort.Ints(minerList)

		minerStr := fmt.Sprintf("%v", minerList)
		if count > 1 {
			forkCount++
			p("  %-6d | %-6d | %-18s  ← FORK\n", h, count, minerStr)
		} else {
			p("  %-6d | %-6d | %-18s\n", h, count, minerStr)
		}
	}
	if forkCount == 0 {
		p("  (no forks detected)\n")
	}
	p("\n")

	// ----- 链 DAG 文本图（简化）-----
	p("--- CHAIN DAG (Simplified Tree View) ---\n")
	p("  Each row: [HEIGHT] HASH.. (MINER) +childCount\n\n")

	printDAGTree(f, analysis)

	p("\n============================================================\n")
	p("  Report generated at: %s\n", now.Format(time.RFC3339))
	p("============================================================\n")

	// 追加连接拓扑报告
	writeConnectionReport(f, input.PeerAssignments, config.NodeCount)

	return nil
}

// WriteJSONReport writes the same simulation snapshot as machine-readable JSON.
func WriteJSONReport(input Input) error {
	analysis, err := analyzeDAG(input.Blocks)
	if err != nil {
		return fmt.Errorf("chain analysis: %w", err)
	}

	type jsonBlock struct {
		Height  int    `json:"height"`
		MinerID int    `json:"miner_id"`
		Hash    string `json:"hash"`
		PreHash string `json:"pre_hash"`
	}
	type jsonPeerAssignment struct {
		NodeID int   `json:"node_id"`
		Peers  []int `json:"peers"`
	}
	type jsonReport struct {
		GeneratedAt string               `json:"generated_at"`
		Parameters  map[string]any       `json:"parameters"`
		Overview    map[string]any       `json:"overview"`
		MinerBlocks map[int]int          `json:"miner_blocks"`
		Blocks      []jsonBlock          `json:"blocks"`
		Tips        []jsonBlock          `json:"tips"`
		Peers       []jsonPeerAssignment `json:"peer_assignments"`
	}
	fullHash := func(hash [32]byte) string { return hex.EncodeToString(hash[:]) }

	blocks := make([]jsonBlock, 0, len(analysis.blocks))
	for _, block := range analysis.blocks {
		blocks = append(blocks, jsonBlock{
			Height: block.Height, MinerID: block.MinerID,
			Hash: fullHash(block.Hash), PreHash: fullHash(block.PreHash),
		})
	}
	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i].Height != blocks[j].Height {
			return blocks[i].Height < blocks[j].Height
		}
		return blocks[i].Hash < blocks[j].Hash
	})

	tips := make([]jsonBlock, 0, len(analysis.tips))
	for _, block := range analysis.tips {
		tips = append(tips, jsonBlock{
			Height: block.Height, MinerID: block.MinerID,
			Hash: fullHash(block.Hash), PreHash: fullHash(block.PreHash),
		})
	}

	peerIDs := make([]int, 0, len(input.PeerAssignments))
	for nodeID := range input.PeerAssignments {
		peerIDs = append(peerIDs, nodeID)
	}
	sort.Ints(peerIDs)
	peers := make([]jsonPeerAssignment, 0, len(peerIDs))
	for _, nodeID := range peerIDs {
		assigned := append([]int(nil), input.PeerAssignments[nodeID]...)
		peers = append(peers, jsonPeerAssignment{NodeID: nodeID, Peers: assigned})
	}

	elapsed := time.Since(input.StartTime).Seconds()
	if elapsed <= 0 {
		elapsed = 1e-9
	}
	data := jsonReport{
		GeneratedAt: time.Now().Format(time.RFC3339Nano),
		Parameters: map[string]any{
			"node_count": config.NodeCount, "difficulty": config.Difficulty,
			"density": config.Density, "target_height": config.TargetHeight, "seed": config.Seed,
		},
		Overview: map[string]any{
			"genesis_hash": fullHash(analysis.genesis.Hash), "total_blocks": analysis.totalBlocks - 1,
			"total_forks": analysis.totalForks, "longest_height": analysis.longestHeight,
			"active_tips": len(analysis.tips), "elapsed_seconds": elapsed,
			"blocks_per_second": float64(analysis.totalBlocks-1) / elapsed,
		},
		MinerBlocks: analysis.minerCounts, Blocks: blocks, Tips: tips, Peers: peers,
	}

	now := time.Now()
	filename := filepath.Join("results", "json", fmt.Sprintf("%s-%s.json", logFilePrefix, now.Format("20060102-150405.000000000")))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("create JSON directory: %w", err)
	}
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encode JSON report: %w", err)
	}
	if err := os.WriteFile(filename, append(encoded, '\n'), 0o644); err != nil {
		return fmt.Errorf("write JSON report: %w", err)
	}
	return nil
}

// printDAGTree 从 genesis 开始递归打印 DAG 的文本树状图。
func printDAGTree(f *os.File, analysis *dagAnalysis) {
	type dfsItem struct {
		hash  [32]byte
		block *p2p.BlockRepr
		level int
	}

	stack := []dfsItem{{hash: analysis.genesis.Hash, block: analysis.genesis, level: 0}}
	visited := make(map[[32]byte]bool)

	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[item.hash] {
			continue
		}
		visited[item.hash] = true

		indent := ""
		for i := 0; i < item.level; i++ {
			indent += "  "
		}

		children := analysis.children[item.hash]
		isLongest := false
		if len(analysis.tips) > 0 {
			longestChain := traceChain(analysis.tips[0], analysis.blocks)
			for _, b := range longestChain {
				if b.Hash == item.hash {
					isLongest = true
					break
				}
			}
		}

		marker := " "
		if isLongest {
			marker = "*"
		}

		if len(children) > 0 {
			fmt.Fprintf(f, "%s[%s] %s | H=%d M=%d | children=%d\n",
				indent, marker, item.block.HashString(),
				item.block.Height, item.block.MinerID, len(children))
		} else {
			fmt.Fprintf(f, "%s[%s] %s | H=%d M=%d | TIP\n",
				indent, marker, item.block.HashString(),
				item.block.Height, item.block.MinerID)
		}

		sort.Slice(children, func(i, j int) bool {
			return analysis.depths[children[i].Hash] > analysis.depths[children[j].Hash]
		})

		for i := len(children) - 1; i >= 0; i-- {
			if !visited[children[i].Hash] {
				stack = append(stack, dfsItem{
					hash:  children[i].Hash,
					block: children[i],
					level: item.level + 1,
				})
			}
		}
	}
}

// ==================== 连接拓扑报告 ====================

// writeConnectionReport 输出节点间连接拓扑信息。
func writeConnectionReport(f *os.File, peerAssignments map[int][]int, nodeCount int) {
	p := func(format string, args ...any) {
		fmt.Fprintf(f, format, args...)
	}

	p("\n============================================================\n")
	p("      CONNECTION TOPOLOGY REPORT\n")
	p("============================================================\n\n")

	if len(peerAssignments) == 0 {
		p("  (no peer assignment data available)\n")
		return
	}

	// 目标 peer 数
	targetCount := int(float64(nodeCount-1)*config.Density + 0.999)
	if targetCount < 1 {
		targetCount = 1
	}

	p("--- PEER ASSIGNMENTS ---\n")
	p("  Target peers per node: %d (density=%.2f)\n\n", targetCount, config.Density)
	p("  %-6s | %-8s | %-8s | %-40s\n", "NODE", "ASSIGNED", "TARGET", "PEER LIST")
	p(sep68)

	// 统计
	minDegree := nodeCount
	maxDegree := 0
	totalDegree := 0

	nodeIDs := make([]int, 0, len(peerAssignments))
	for id := range peerAssignments {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Ints(nodeIDs)

	for _, nodeID := range nodeIDs {
		peers := peerAssignments[nodeID]
		degree := len(peers)
		totalDegree += degree
		if degree < minDegree {
			minDegree = degree
		}
		if degree > maxDegree {
			maxDegree = degree
		}

		// 截断显示 peer 列表
		peerStr := fmt.Sprintf("%v", peers)
		if len(peerStr) > 40 {
			peerStr = peerStr[:37] + "..."
		}

		status := "✓"
		if degree != targetCount {
			status = "✗"
		}
		p("  %-6d | %-8d | %-8d | %-40s %s\n", nodeID, degree, targetCount, peerStr, status)
	}
	p("\n")

	// 汇总统计
	avgDegree := float64(totalDegree) / float64(nodeCount)
	p("--- DEGREE STATISTICS ---\n")
	p("  Min degree:  %d\n", minDegree)
	p("  Max degree:  %d\n", maxDegree)
	p("  Avg degree:  %.1f\n", avgDegree)
	p("  Target:      %d\n", targetCount)

	allSatisfied := minDegree == targetCount && maxDegree == targetCount
	if allSatisfied {
		p("  Status:      ✓ All nodes meet density target\n")
	} else {
		p("  Status:      ✗ Some nodes do NOT meet density target\n")
	}
	p("\n")

	// 邻接矩阵（紧凑格式）
	if nodeCount <= 30 {
		p("--- ADJACENCY MATRIX ---\n")
		p("     ")
		for i := 0; i < nodeCount; i++ {
			p("%2d ", i)
		}
		p("\n")
		// 列头分隔线
		fmt.Fprintf(f, "     %s\n", strings.Repeat("-", nodeCount*3))
		for i := 0; i < nodeCount; i++ {
			peers := peerAssignments[i]
			peerSet := make(map[int]bool, len(peers))
			for _, pid := range peers {
				peerSet[pid] = true
			}
			p("  %2d |", i)
			for j := 0; j < nodeCount; j++ {
				if i == j {
					p(" · ")
				} else if peerSet[j] {
					p(" ● ")
				} else {
					p(" · ")
				}
			}
			p("\n")
		}
		p("  Legend: ● = connected, · = not connected\n")
		p("\n")
	} else {
		p("  (adjacency matrix skipped for >30 nodes)\n\n")
	}
}
