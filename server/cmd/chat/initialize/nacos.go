package initialize

import (
	"net"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/config"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/bwmarrin/snowflake"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/utils"
	nacos "github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

// InitNacos to init nacos
func InitNacos(Port int) (registry.Registry, *registry.Info) {
	v := viper.New()
	v.SetConfigFile(consts.ChatConfigPath)
	if err := v.ReadInConfig(); err != nil {
		klog.Fatalf("read viper config failed: %s", err)
	}
	if err := v.Unmarshal(&config.GlobalNacosConfig); err != nil {
		klog.Fatalf("unmarshal err failed: %s", err)
	}
	klog.Infof("Config Info: %v", config.GlobalNacosConfig)

	// Read configuration information from nacos
	sc := []constant.ServerConfig{
		{
			IpAddr: config.GlobalNacosConfig.Host,
			Port:   config.GlobalNacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         config.GlobalNacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              consts.NacosLogDir,
		CacheDir:            consts.NacosCacheDir,
		LogLevel:            consts.NacosLogLevel,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		klog.Fatalf("create config client failed: %s", err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: config.GlobalNacosConfig.DataId,
		Group:  config.GlobalNacosConfig.Group,
	})
	if err != nil {
		klog.Fatalf("get config failed: %s", err.Error())
	}

	err = sonic.Unmarshal([]byte(content), &config.GlobalServerConfig)
	if err != nil {
		klog.Fatalf("nacos config failed: %s", err)
	}

	if config.GlobalServerConfig.Host == "" {
		config.GlobalServerConfig.Host, err = tools.GetLocalIPv4Address()
		if err != nil {
			klog.Fatalf("get localIpv4Addr failed:%s", err.Error())
		}
	}

	registryClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		klog.Fatalf("create registryClient err: %s", err.Error())
	}

	r := nacos.NewNacosRegistry(registryClient, nacos.WithGroup(consts.ChatGroup))

	sf, err := snowflake.NewNode(consts.NacosSnowflakeNode)
	if err != nil {
		klog.Fatalf("generate service name failed: %s", err)
	}
	info := &registry.Info{
		ServiceName: config.GlobalServerConfig.Name,
		Addr:        utils.NewNetAddr(consts.TCP, net.JoinHostPort(config.GlobalServerConfig.Host, strconv.Itoa(Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
	}
	return r, info
}
