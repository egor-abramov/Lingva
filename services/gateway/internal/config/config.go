package config

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	GRPCClients GRPCClientsConfig `yaml:"grpc_clients"`
}

type ServerConfig struct {
	Port int `yaml:"port" default:"8080"`
}

type GRPCClientsConfig struct {
	Analyzer string `yaml:"analyzer" default:"localhost:50051"`
	Runner   string `yaml:"runner" default:"localhost:50052"`
}
