package metrics

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

func PingHost(host string, count int) error {
	// Execute ping command
	cmd := exec.Command("ping", "-c", fmt.Sprint(count), host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute ping: %s", err)
	}

	avgLatency := calculateAverageLatency(string(output), count)
	packetLoss := calculatePacketLoss(string(output))
	jitter := calculateJitter(string(output))

	// Output results
	fmt.Printf("Average Latency: %.2f ms\n", avgLatency)
	fmt.Printf("Packet Loss: %s%%\n", packetLoss)
	fmt.Printf("Jitter: %.2f ms\n", jitter)
	return nil
}

func calculatePacketLoss(output string) string {
	packetLossRegex := regexp.MustCompile(`(\d+)% packet loss`)
	packetLossMatch := packetLossRegex.FindStringSubmatch(output)

	if len(packetLossMatch) > 1 {
		return packetLossMatch[1]
	}
	return "0"
}

func calculateAverageLatency(output string, count int) float64 {
	// Regex to find latency values
	latencyRegex := regexp.MustCompile(`time=(\d+\.\d+)`)
	latencies := latencyRegex.FindAllStringSubmatch(output, -1)

	var totalLatency float64

	// Calculate total latency
	for _, latencyStr := range latencies {
		latency, err := strconv.ParseFloat(latencyStr[1], 64)
		if err != nil {
			fmt.Printf("failed to parse latency value: %s\n", err)
			count--
		}
		totalLatency += latency
	}

	// Calculate average latency
	avgLatency := totalLatency / float64(count)

	return avgLatency
}

func calculateJitter(output string) float64 {
	// Regex to find latency values
	latencyRegex := regexp.MustCompile(`time=(\d+\.\d+)`)
	latencies := latencyRegex.FindAllStringSubmatch(output, -1)

	if len(latencies) < 2 {
		// Not enough data to calculate jitter
		return 0.0
	}

	var totalJitter float64
	var previousLatency float64

	// Initialize previousLatency with the first latency value
	firstLatency, err := strconv.ParseFloat(latencies[0][1], 64)
	if err != nil {
		fmt.Printf("failed to parse first latency value: %s\n", err)
		return 0.0
	}
	previousLatency = firstLatency

	for i := 1; i < len(latencies); i++ {
		currentLatency, err := strconv.ParseFloat(latencies[i][1], 64)
		if err != nil {
			fmt.Printf("failed to parse latency value: %s\n", err)
			continue
		}

		jitter := currentLatency - previousLatency
		if jitter < 0 {
			jitter = -jitter
		}
		totalJitter += jitter
		previousLatency = currentLatency
	}

	avgJitter := totalJitter / float64(len(latencies)-1)

	return avgJitter
}
