package config

type Config struct {
	Nacos struct {
		NamespaceId string
		UserName    string
		Password    string
		IpAddr      string
		Port        int
		DataId      string
		Group       string
	}
	Mysql struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
}
