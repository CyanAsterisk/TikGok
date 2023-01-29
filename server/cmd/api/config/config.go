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

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name               string       `mapstructure:"name" json:"name"`
	Host               string       `mapstructure:"host" json:"host"`
	Port               int          `mapstructure:"port" json:"port"`
	JWTInfo            JWTConfig    `mapstructure:"jwt" json:"jwt"`
	OtelInfo           OtelConfig   `mapstructure:"otel" json:"otel"`
	UserSrvInfo        RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	InteractionSrvInfo RPCSrvConfig `mapstructure:"interaction_srv" json:"interaction_srv"`
	SocialitySrvInfo   RPCSrvConfig `mapstructure:"sociality_srv" json:"sociality_srv"`
	VideoSrvInfo       RPCSrvConfig `mapstructure:"video_srv" json:"video_srv"`
	MinioInfo          MinioConfig  `mapstructure:"minio" json:"minio"`
}

type RPCSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint"`
	AccessKeyID     string `mapstructure:"access-key-id" json:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key" json:"secret-access-key"`
}
