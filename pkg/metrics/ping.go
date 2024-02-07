package metrics

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/iliasgal/network-monitor/pkg/model"
)

// PingHost executes the ping command and calculates average latency, packet loss, and jitter.
func PingHost(host string, count int) (*model.PingMetrics, error) {
	cmd := exec.Command("ping", "-c", fmt.Sprint(count), host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute ping: %s", err)
	}

	latencies, err := extractLatencies(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to extract latencies: %s", err)
	}

	avgLatency := calculateAverageLatency(latencies)
	jitter := calculateJitter(latencies)
	packetLoss, err := calculatePacketLoss(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to calculate packet loss: %s", err)
	}

	// write ping metrics to struct
	pingMetrics := model.PingMetrics{
		AvgLatency: avgLatency,
		Jitter:     jitter,
		PacketLoss: packetLoss,
	}

	return &pingMetrics, nil
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

func calculatePacketLoss(output string) (float64, error) {
	// Regex to extract packet loss percentage from ping output
	packetLossRegex := regexp.MustCompile(`(\d+)% packet loss`)
	match := packetLossRegex.FindStringSubmatch(output)

	if len(match) > 1 {
		packetLoss, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse packet loss value: %s", err)
		}
		return packetLoss, nil
	}

	return 0, nil
}
