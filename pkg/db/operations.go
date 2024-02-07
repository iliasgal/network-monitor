package db

import (
	"time"

	"github.com/iliasgal/network-monitor/pkg/model"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func WritePingMetricsToInfluxDB(avgLatency, jitter, packetLoss float64) {
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
}

func WritePacketInfoToInfluxDB(info *model.PacketInfo) {
	// Get a non-blocking write client
	writeAPI := influxClient.WriteAPI(influxOrg, influxBucket)

	// Create a new point with the measurement name
	p := influxdb2.NewPointWithMeasurement("network_traffic").
		AddTag("packet_type", info.PacketType).
		AddTag("src_ip", info.SrcIP).
		AddTag("dst_ip", info.DstIP).
		AddTag("src_port", info.SrcPort).
		AddTag("dst_port", info.DstPort).
		AddField("packet_size", info.Size).
		SetTime(time.Now())

	// Write the point asynchronously
	writeAPI.WritePoint(p)

	// Ensure all writes are sent
	writeAPI.Flush()
}
