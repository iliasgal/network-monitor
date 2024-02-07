package main

import (
	"fmt"
	"time"

	"github.com/iliasgal/network-monitor/pkg/db"
	"github.com/iliasgal/network-monitor/pkg/metrics"
)

func main() {

	host := "google.com"
	count := 4
	for {
		metrics, err := metrics.PingHost(host, count)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write ping metrics to InfluxDB
		go db.WritePingMetricsToInfluxDB(metrics.AvgLatency, metrics.Jitter, metrics.PacketLoss)

		time.Sleep(5 * time.Second)
	}
}
