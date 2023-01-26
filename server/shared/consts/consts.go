package consts

import "time"

const (
	TikGok         = "TikGok"
	MySqlDSN       = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	UserMigrateDSN = "root:123456@tcp(localhost:3306)/TikGok?charset=utf8mb4&parseTime=True&loc=Local"

	JWTIssuer        = "TikGok"
	ThirtyDays       = 60 * 60 * 24 * 30
	AuthorizationKey = "authorization"
	TokenRefreshTime = time.Hour
	Claims           = "claims"
	AccountID        = "accountID"

	ApiConfigPath         = "./server/cmd/api/config.yaml"
	UserConfigPath        = "./server/cmd/user/config.yaml"
	InteractionConfigPath = "./server/cmd/interaction/config.yaml"
	SocialityConfigPath   = "./server/cmd/sociality/config.yaml"
	VideoConfigPath       = "./server/cmd/video/config.yaml"

	NacosSnowflakeNode    = 1
	UserSnowflakeNode     = 2
	VideoSnowflakeNode    = 3
	CommentSnowflakeNode  = 4
	FavoriteSnowflakeNode = 5

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

	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0"
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"

	TCP = "tcp"

	FreePortAddress = "localhost:0"

	InvalidComment = 0
	ValidComment   = 1

	IsNotLike = 0
	IsLike    = 1

	MySQLImage         = "mysql:latest"
	MySQLContainerPort = "3306/tcp"
	MySQLContainerIP   = "127.0.0.1"
	MySQLPort          = "0"
	MySQLAdmin         = "root"
)
