package exp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stellashiina/pow-simulator/exp"
)

var randSeed2 = &[2]uint64{7, 21}

func runSelfishSuite(t *testing.T, name string, difficulty float64, nodeCount int) {
	t.Helper()
	rates := []float64{0.4, 0.3, 0.2, 0.1}
	results := make([]exp.SelfishMiningResult, 0, len(rates))
	var log strings.Builder
	fmt.Fprintf(&log, "============================================================\n")
	fmt.Fprintf(&log, "      PoW TICK SIMULATION - SELFISH MINING\n")
	fmt.Fprintf(&log, "============================================================\n\n")
	fmt.Fprintf(&log, "--- PARAMETERS ---\n  Difficulty: %.6f\n  Node Count: %d\n  Random Seed: %v\n\n", difficulty, nodeCount, *randSeed2)
	fmt.Fprintf(&log, "--- RESULTS ---\n  %-16s | %-10s | %-16s | %-14s\n", "MALICIOUS RATE", "ROUNDS", "INCOME BLOCKS", "PROFITABILITY")
	fmt.Fprintf(&log, "  ------------------------------------------------------------------------\n")
	for _, rate := range rates {
		result := exp.Selfish(rate, difficulty, nodeCount, randSeed2)
		results = append(results, result)
		fmt.Fprintf(&log, "  %-16.2f | %-10d | %-16d | %12.2f%%\n", rate, result.Rounds, result.IncomeBlocks, result.Profitability*100)
	}
	fmt.Fprintf(&log, "\n============================================================\n")
	writeExperimentFiles(t, name, log.String(), selfishSuiteResult{Experiment: name, Results: results})
	t.Log(log.String())
}

type selfishSuiteResult struct {
	Experiment string                    `json:"experiment"`
	Results    []exp.SelfishMiningResult `json:"results"`
}

func TestSelfishVar1(t *testing.T) { runSelfishSuite(t, "selfish-mining-var1", 0.001, 100) }
func TestSelfishVar2(t *testing.T) { runSelfishSuite(t, "selfish-mining-var2", 0.001, 50) }
func TestSelfishVar3(t *testing.T) { runSelfishSuite(t, "selfish-mining-var3", 0.0005, 100) }
