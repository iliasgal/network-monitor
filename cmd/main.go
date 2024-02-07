package main

import (
	"time"

	met "github.com/iliasgal/network-monitor/pkg/metrics"
)

func main() {

	// host := "google.com"
	// count := 4
	for {
		// metrics, err := met.PingHost(host, count)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		go met.PacketCapture()

		// // Write ping metrics to InfluxDB
		// go db.WritePingMetricsToInfluxDB(metrics.AvgLatency, metrics.Jitter, metrics.PacketLoss)

		time.Sleep(5 * time.Second)
	}
}
