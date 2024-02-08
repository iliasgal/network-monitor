package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Note the underscore, which imports the package for its side-effects only.
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iliasgal/network-monitor/pkg/db"
	"github.com/iliasgal/network-monitor/pkg/metrics"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Start packet capture in its own goroutine
	go metrics.PacketCapture()

	// Setup channel to listen for termination signals
	signals := make(chan os.Signal, 1)
	// Notify signals channel on SIGINT and SIGTERM
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

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
			go db.WritePingMetricsToInfluxDB(pingStats)
		case <-signals:
			// Received a termination signal, perform cleanup
			fmt.Println("Termination signal received, closing resources.")
			db.CloseInfluxDBClient()
			ticker.Stop()
			return
		}
	}
}
