package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stellashiina/pow-simulator/config"
	"github.com/stellashiina/pow-simulator/simulator"
)

type parameterResult struct {
	NodeCount     int     `json:"node_count"`
	Difficulty    float64 `json:"difficulty"`
	Ticks         uint64  `json:"ticks"`
	TicksPerBlock float64 `json:"ticks_per_block"`
}

type parameterReport struct {
	Experiment string            `json:"experiment"`
	Results    []parameterResult `json:"results"`
}

func TestParameters(t *testing.T) {
	nodeCounts := []int{10, 50, 100, 200, 500, 1000}
	difficulties := []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.25, 0.3, 0.5}
	results := make([]parameterResult, 0)
	var log strings.Builder
	fmt.Fprintln(&log, "============================================================")
	fmt.Fprintln(&log, "      PoW TICK SIMULATION - PARAMETER SWEEP")
	fmt.Fprintln(&log, "============================================================")
	fmt.Fprintln(&log, "\n--- RESULTS ---")
	fmt.Fprintf(&log, "  %-12s | %-12s | %-10s | %-16s\n", "NODE COUNT", "DIFFICULTY", "TICKS", "TICKS / BLOCK")
	fmt.Fprintln(&log, "  ------------------------------------------------------------")
	for _, nc := range nodeCounts {
		for _, d := range difficulties {
			if float64(nc)*d > 10 {
				continue
			}
			config.NodesCount = nc
			config.Difficulty = d
			sim := simulator.NewSimulator()
			sim.Run()
			result := parameterResult{NodeCount: nc, Difficulty: d, Ticks: sim.Tick, TicksPerBlock: float64(sim.Tick) / float64(config.Target)}
			results = append(results, result)
			fmt.Fprintf(&log, "  %-12d | %-12.6f | %-10d | %16.2f\n", nc, d, result.Ticks, result.TicksPerBlock)
		}
	}
	fmt.Fprintln(&log, "\n============================================================")
	writeParameterFiles(t, log.String(), parameterReport{Experiment: "parameter_sweep", Results: results})
	t.Log(log.String())
}

func writeParameterFiles(t *testing.T, logText string, data parameterReport) {
	t.Helper()
	timestamp := time.Now().Format("20060102-150405.000000000")
	logPath := filepath.Join("results", "log", "parameter-sweep-"+timestamp+".log")
	jsonPath := filepath.Join("results", "json", "parameter-sweep-"+timestamp+".json")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(logPath, []byte(logText), 0o644); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(jsonPath, append(encoded, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("Reports written: %s and %s", logPath, jsonPath)
}
