package ping

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

	// Output results
	fmt.Printf("Average Latency: %.2f ms\n", avgLatency)
	fmt.Printf("Packet Loss: %s%%\n", packetLoss)
	return nil
}

func calculateAverageLatency(output string, count int) float64 {
	// Regex to find average latency from ping output
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

func calculatePacketLoss(output string) string {
	packetLossRegex := regexp.MustCompile(`(\d+)% packet loss`)
	packetLossMatch := packetLossRegex.FindStringSubmatch(output)

	if len(packetLossMatch) > 1 {
		return packetLossMatch[1]
	}
	return "0"
}
