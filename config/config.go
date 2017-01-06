package config

type Config struct {
	Zookeeper `json:"zookeeper"`
}

type Zookeeper struct {
	Servers         string                       `json:"servers"`
	NormalServerMap map[string]map[string]string `json:"normal"`
	GrpcServerMap   map[string]map[string]string `json:"grpc"`
}

type NormalServer struct {
	UVERManager string `json:"uvermanager"`
	UDatabase   string `json:"udatabase"`
}

type GrpcServer struct {
	HelloServer string `json:"hellomanager"`
}
