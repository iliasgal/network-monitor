package metrics

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// PingHost executes the ping command and calculates average latency, packet loss, and jitter.
func PingHost(host string, count int) error {
	cmd := exec.Command("ping", "-c", fmt.Sprint(count), host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute ping: %s", err)
	}

	latencies, err := extractLatencies(string(output))
	if err != nil {
		return fmt.Errorf("failed to extract latencies: %s", err)
	}

	avgLatency := calculateAverageLatency(latencies)
	jitter := calculateJitter(latencies)
	packetLoss := calculatePacketLoss(string(output))

	fmt.Printf("Average Latency: %.2f ms\n", avgLatency)
	fmt.Printf("Jitter: %.2f ms\n", jitter)
	fmt.Printf("Packet Loss: %s%%\n", packetLoss)
	return nil
}

func extractLatencies(output string) ([]float64, error) {
	// Regex to extract latencies from ping output
	latencyRegex := regexp.MustCompile(`time=(\d+\.\d+)`)
	matches := latencyRegex.FindAllStringSubmatch(output, -1)

	var latencies []float64
	for _, match := range matches {
		latency, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return nil, err
		}
		latencies = append(latencies, latency)
	}
	return latencies, nil
}

func calculateAverageLatency(latencies []float64) float64 {
	var totalLatency float64
	for _, latency := range latencies {
		totalLatency += latency
	}
	return totalLatency / float64(len(latencies))
}

func calculateJitter(latencies []float64) float64 {
	// Jitter requires at least two latencies
	if len(latencies) < 2 {
		return 0.0
	}

	// Jitter is the average of the absolute difference between latencies
	var totalJitter float64
	for i := 1; i < len(latencies); i++ {
		diff := latencies[i] - latencies[i-1]
		if diff < 0 {
			diff = -diff
		}
		totalJitter += diff
	}
	return totalJitter / float64(len(latencies)-1)
}

func calculatePacketLoss(output string) string {
	// Regex to extract packet loss from ping output
	packetLossRegex := regexp.MustCompile(`(\d+)% packet loss`)
	match := packetLossRegex.FindStringSubmatch(output)
	if len(match) > 1 {
		return match[1]
	}
	return "0"
}
