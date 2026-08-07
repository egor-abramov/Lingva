package config

type Config struct {
	GRPC GRPCClient `yaml:"grpc"`
}

type GRPCClient struct {
	Port int `yaml:"port" default:"50052"`
}
