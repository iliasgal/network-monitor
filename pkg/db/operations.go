package db

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func WritePingMetricsToInfluxDB(avgLatency, jitter, packetLoss float64) error {
	// Get non-blocking write client
	writeAPI := influxClient.WriteAPI(influxOrg, influxBucket)

	// Create a point and add to batch
	p := influxdb2.NewPointWithMeasurement("ping_metrics").
		AddField("avg_latency_ms", avgLatency).
		AddField("jitter_ms", jitter).
		AddField("packet_loss", packetLoss).
		SetTime(time.Now())

	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	return nil
}
