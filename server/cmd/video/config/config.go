package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Salt     string `mapstructure:"salt" json:"salt"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}

type RabbitMqConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Exchange string `mapstructure:"exchange" json:"exchange"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name               string               `mapstructure:"name" json:"name"`
	Host               string               `mapstructure:"host" json:"host"`
	RedisInfo          RedisConfig          `mapstructure:"redis" json:"redis"`
	MysqlInfo          MysqlConfig          `mapstructure:"mysql" json:"mysql"`
	RabbitMqInfo       RabbitMqConfig       `mapstructure:"rabbitmq" json:"rabbitmq"`
	OtelInfo           OtelConfig           `mapstructure:"otel" json:"otel"`
	InteractionSrvInfo InteractionSrvConfig `mapstructure:"interaction_srv" json:"interaction_srv"`
	UserSrvInfo        UserSrvConfig        `mapstructure:"user_srv" json:"user_srv"`
}

type InteractionSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type UserSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
