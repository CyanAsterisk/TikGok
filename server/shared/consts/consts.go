package consts

const (
	TikGok         = "TikGok"
	MySqlDSN       = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	UserMigrateDSN = "root:123456@tcp(localhost:3306)/TikGok?charset=utf8mb4&parseTime=True&loc=Local"

	JWTIssuer        = "TikGok"
	ThirtyDays       = 60 * 60 * 24 * 30
	AuthorizationKey = "token"
	Claims           = "claims"
	AccountID        = "accountID"

	ApiConfigPath         = "./server/cmd/api/config.yaml"
	UserConfigPath        = "./server/cmd/user/config.yaml"
	InteractionConfigPath = "./server/cmd/interaction/config.yaml"
	SocialityConfigPath   = "./server/cmd/sociality/config.yaml"
	VideoConfigPath       = "./server/cmd/video/config.yaml"
	ChatConfigPath        = "./server/cmd/chat/config.yaml"

	NacosSnowflakeNode    = 1
	UserSnowflakeNode     = 2
	VideoSnowflakeNode    = 3
	CommentSnowflakeNode  = 4
	FavoriteSnowflakeNode = 5
	FollowSnowflakeNode   = 6
	MinioSnowflakeNode    = 7
	ChatSnowflakeNode     = 8

	ApiGroup         = "API_GROUP"
	UserGroup        = "USER_GROUP"
	InteractionGroup = "INTERACTION_GROUP"
	SocialityGroup   = "SOCIALITY_GROUP"
	VideoGroup       = "VIDEO_GROUP"
	ChatGroup        = "CHAT_GROUP"

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

	InvalidComment = 2
	ValidComment   = 1

	IsNotLike = 2
	IsLike    = 1
	Like      = 1
	UnLike    = 2

	IsNotFollow = 2
	IsFollow    = 1

	SentMessage    = 1
	ReceiveMessage = 0

	MySQLImage         = "mysql:latest"
	MySQLContainerPort = "3306/tcp"
	MySQLContainerIP   = "127.0.0.1"
	MySQLPort          = "0"
	MySQLAdmin         = "root"
	DockerTestMySQLPwd = "123456"

	RedisImage         = "redis:latest"
	RedisContainerPort = "6379/tcp"
	RedisContainerIP   = "127.0.0.1"
	RedisPort          = "0"

	VideosLimit = 10

	RabbitMqURI = "amqp://%s:%s@%s:%d/"

	RedisSocialClientDB   = 1
	RedisVideoClientDB    = 2
	RedisCommentClientDB  = 3
	RedisFavoriteClientDB = 4
	RedisUserClientDB     = 5

	RedisFollowerSuffix = "follower"
	RedisFollowSuffix   = "follow"
	AllVideoSortSetKey  = "all-video-list-key"
	VideoIdFiled        = "videoId"
	CommentJsonFiled    = "commentJson"

	UsernameFiled = "username"
	CryptPwdFiled = "password"

	FollowList   = 0
	FollowerList = 1
	FriendsList  = 2

	FollowCount   = 0
	FollowerCount = 1
)
