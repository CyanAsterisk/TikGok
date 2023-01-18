package consts

const (
	JWTIssuer  = "FreeCar"
	ThirtyDays = 60 * 60 * 24 * 30

	ApiConfigPath         = "./server/cmd/api/config.yaml"
	UserConfigPath        = "./server/cmd/user/config.yaml"
	InteractionConfigPath = "./server/cmd/interaction/config.yaml"
	SocialityConfigPath   = "./server/cmd/sociality/config.yaml"
	VideoConfigPath       = "./server/cmd/video/config.yaml"

	ApiGroup         = "API_GROUP"
	UserGroup        = "AUTH_GROUP"
	InteractionGroup = "INTERACTION_GROUP"
	SocialityGroup   = "SOCIALITY_GROUP"
	VideoGroup       = "VIDEO_GROUP"

	NacosLogDir   = "tmp/nacos/log"
	NacosCacheDir = "tmp/nacos/cache"
	NacosLogLevel = "debug"

	HlogFilePath = "./tmp/hlog/logs/"
	KlogFilePath = "./tmp/klog/logs/"

	MySqlDSN = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"

	TCP = "tcp"
)
