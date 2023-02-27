package config

type RabbitMqConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Exchange string `mapstructure:"exchange" json:"exchange"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
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
