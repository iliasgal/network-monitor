package db

import (
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	influxToken  = os.Getenv("INFLUXDB_ADMIN_TOKEN")
	influxBucket = os.Getenv("INFLUXDB_BUCKET")
	influxOrg    = os.Getenv("INFLUXDB_ORG")
	influxURL    = os.Getenv("INFLUXDB_URL")
)

var influxClient influxdb2.Client

func init() {
	influxClient = influxdb2.NewClient(influxURL, influxToken)
}

func CloseInfluxDBClient() {
	influxClient.Close()
}
