package exp_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stellashiina/pow-simulator/exp"
)

var randSeed1 = &[2]uint64{7, 21}

type forkSuiteResult struct {
	Experiment string                 `json:"experiment"`
	Results    []exp.ForkAttackResult `json:"results"`
}

func runForkSuite(t *testing.T, name string, difficulty float64, nodeCount int) {
	t.Helper()
	rates := []float64{0.4, 0.3, 0.2, 0.1}
	results := make([]exp.ForkAttackResult, 0, len(rates))
	var log strings.Builder
	fmt.Fprintf(&log, "============================================================\n")
	fmt.Fprintf(&log, "      PoW TICK SIMULATION - FORK ATTACK\n")
	fmt.Fprintf(&log, "============================================================\n\n")
	fmt.Fprintf(&log, "--- PARAMETERS ---\n  Difficulty: %.6f\n  Node Count: %d\n  Random Seed: %v\n\n", difficulty, nodeCount, *randSeed1)
	fmt.Fprintf(&log, "--- RESULTS ---\n  %-16s | %-10s | %-12s | %-12s\n", "MALICIOUS RATE", "ROUNDS", "SUCCESSES", "SUCCESS RATE")
	fmt.Fprintf(&log, "  ------------------------------------------------------------\n")
	for _, rate := range rates {
		result := exp.ForkAttack(rate, difficulty, nodeCount, randSeed1)
		results = append(results, result)
		fmt.Fprintf(&log, "  %-16.2f | %-10d | %-12d | %10.2f%%\n", rate, result.Rounds, result.Successes, result.SuccessRate*100)
	}
	fmt.Fprintf(&log, "\n============================================================\n")
	writeExperimentFiles(t, name, log.String(), forkSuiteResult{Experiment: name, Results: results})
	t.Log(log.String())
}

func writeExperimentFiles(t *testing.T, name, logText string, data any) {
	t.Helper()
	timestamp := time.Now().Format("20060102-150405.000000000")
	base := filepath.Join("..", "results")
	logPath := filepath.Join(base, "log", fmt.Sprintf("%s-%s.log", name, timestamp))
	jsonPath := filepath.Join(base, "json", fmt.Sprintf("%s-%s.json", name, timestamp))
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(logPath, []byte(logText), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := exp.WriteJSON(jsonPath, data); err != nil {
		t.Fatal(err)
	}
	t.Logf("Reports written: %s and %s", logPath, jsonPath)
}

func TestForkATKVar1(t *testing.T) { runForkSuite(t, "fork-attack-var1", 0.001, 100) }
func TestForkATKVar2(t *testing.T) { runForkSuite(t, "fork-attack-var2", 0.001, 50) }
func TestForkATKVar3(t *testing.T) { runForkSuite(t, "fork-attack-var3", 0.0005, 100) }
