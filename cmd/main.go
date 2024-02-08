package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iliasgal/network-monitor/pkg/db"
	"github.com/iliasgal/network-monitor/pkg/metrics"
)

func main() {
	// Create a channel to receive errors
	errorChan := make(chan error)

	// Setup channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	// Notify for SIGINT and SIGTERM
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start packet capture in its own goroutine
	go metrics.PacketCapture()

	host := "google.com"
	count := 4
	ticker := time.NewTicker(5 * time.Second) // Ping every 5 seconds

	for {
		select {
		case <-ticker.C:
			pingStats, err := metrics.PingHost(host, count)
			if err != nil {
				log.Fatal(err)
				return
			}
			db.WritePingMetricsToInfluxDB(pingStats, errorChan)
		case info := <-metrics.PacketInfoChan:
			db.WritePacketInfoToInfluxDB(info, errorChan)
		case err := <-errorChan:
			log.Printf("Error occurred: %v", err)
		case <-signals:
			log.Println("Termination signal received, closing resources.")
			db.CloseInfluxDBClient()
			ticker.Stop()
			return
		}
	}
}
