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

type RabbitMqConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Exchange string `mapstructure:"exchange" json:"exchange"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ServerConfig struct {
	Name               string              `mapstructure:"name" json:"name"`
	Host               string              `mapstructure:"host" json:"host"`
	Port               int                 `mapstructure:"port" json:"port"`
	JWTInfo            JWTConfig           `mapstructure:"jwt" json:"jwt"`
	OtelInfo           OtelConfig          `mapstructure:"otel" json:"otel"`
	ChatSrvInfo        RPCSrvConfig        `mapstructure:"chat_srv" json:"chat_srv"`
	UserSrvInfo        RPCSrvConfig        `mapstructure:"user_srv" json:"user_srv"`
	InteractionSrvInfo RPCSrvConfig        `mapstructure:"interaction_srv" json:"interaction_srv"`
	SocialitySrvInfo   RPCSrvConfig        `mapstructure:"sociality_srv" json:"sociality_srv"`
	VideoSrvInfo       RPCSrvConfig        `mapstructure:"video_srv" json:"video_srv"`
	UploadServiceInfo  UploadServiceConfig `mapstructure:"upload_srv" json:"upload_srv"`
}

type RPCSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key" json:"secret_access_key"`
	Bucket          string `mapstructure:"bucket" json:"bucket"`
	UrlPrefix       string `mapstructure:"url_prefix" json:"url_prefix"`
}

type UploadServiceConfig struct {
	MinioInfo    MinioConfig    `mapstructure:"minio" json:"minio"`
	RabbitMqInfo RabbitMqConfig `mapstructure:"rabbitmq" json:"rabbitmq"`
}
