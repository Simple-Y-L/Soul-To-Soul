package inits

import (
	"encoding/json"
	"kratosItems/cmd/basic/config"
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func InitNaCos() {
	NacosConf := config.Configs.Nacos

	//create clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         NacosConf.NamespaceId, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            NacosConf.UserName,
		Password:            NacosConf.Password,
	}

	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      NacosConf.IpAddr,
			ContextPath: "/nacos",
			Port:        uint64(NacosConf.Port),
			Scheme:      "http",
		},
	}

	// Another way of create config client for dynamic configuration (recommend)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: NacosConf.DataId,
		Group:  NacosConf.Group})
	json.Unmarshal([]byte(content), &config.Configs)
	configClient.ListenConfig(vo.ConfigParam{
		DataId: NacosConf.DataId,
		Group:  NacosConf.Group,
		OnChange: func(namespace, group, dataId, data string) {
			json.Unmarshal([]byte(content), &config.Configs)
		},
	})
	log.Println(config.Configs)

}
