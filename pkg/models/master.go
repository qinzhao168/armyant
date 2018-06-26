package models

type MasterConfig struct {
	HttpPort       string
	GrpcPort       string
	Storage        string
	DataDir        string
	InfluxDB       string
	InfluxUrl      string
	InfluxUsername string
	InfluxPassword string
}
