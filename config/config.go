package config

type Config struct {
	Zookeeper `json:"zookeeper"`
}

type Zookeeper struct {
	Servers string `json:"servers"`
	UVERManager string `json:"uvermanager"`
	UDatabase string `json:"udatabase"`
}

