package pkg

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/yiqinguo/armyant/pkg/models"
)

func ParseMasterArgs() models.MasterConfig {
	var config models.MasterConfig
	kingpin.Flag("http-port", "listen http port").Default("8080").StringVar(&config.HttpPort)
	kingpin.Flag("grpc-port", "listen grpc port").Default("50051").StringVar(&config.GrpcPort)
	kingpin.Flag("storage", "storage type").Default("influxdb").StringVar(&config.Storage)
	kingpin.Flag("datadir", "local storage dir").Default("armyant.db").StringVar(&config.DataDir)
	kingpin.Flag("influx-url", "influx url").StringVar(&config.InfluxUrl)
	kingpin.Flag("influx-db", "influx db").StringVar(&config.InfluxDB)
	kingpin.Flag("influx-username", "influx username").StringVar(&config.InfluxUsername)
	kingpin.Flag("influx-password", "influx password").StringVar(&config.InfluxPassword)

	kingpin.Parse()
	return config
}
