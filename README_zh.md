![TikGok](img/TikGok.png)

[English](README.md) | 中文

TikGok 是一个基于 Hertz 与 Kitex 的极简版抖音，按要求实现了全部的接口，并对业务进行了优化。除此之外还提升了项目的治理能力，例如引入了配置中心、服务中心、OTEL 等技术栈。

## 快速开始

TODO

## 项目实现

### 技术选型

- HTTP 框架使用 Hertz
- RPC 框架使用 Kitex
- 关系型数据库选用 MySQL
- 非关系型数据库选用 Redis
- 服务中心与配置中心均选用 Nacos
- 对象存储服务使用 Minio
- 消息队列使用 RabbitMQ
- 使用 Jaeger 与 Prometheus 进行链路追踪以及监控
- CI 使用 Github Actions

### 架构设计

#### 调用关系

![img](./img/call_relation.png)

#### 技术架构

TODO

#### 服务关系

TODO

#### MySQL 数据库设计

![img](./img/mysql.png)

##### 索引设计

###### User

`username` 设唯一索引

###### Favorite

`user_id` 与 `video_id` 设联合唯一索引

###### Video

`user_id` 设普通索引

###### Comment

`video_id` 设普通索引

##### Follow

`user_id` 与 `follower_id` 设联合唯一索引

###### Message

`to_user_id` 与 `from_user_id` 设联合索引

#### Redis 数据库设计

##### User

`用户id`作为 key，`用户信息`作为 value。

##### Video

`视频id`作为 key，`视频信息`作为 value。

`用户id`作为有序集合的 key，用户发布的`视频id`为 member ，排序分数为`发布时间`。

单独维护一个key为`consts.AllVideoSortSetKey`的有序集合，所有所有用户发布的`视频id`为有序集合的 member，排序分数为`发布时间`。

##### Interaction

- commentDB
    -  `评论id`作为 key，`评论信息`作为 value。

    -  `视频id`作为有序集合的 key，对该视频的`评论id`为 member，排序分数是`视频发布时间`。
- favoriteDB

`用户id`作为有序集合的 key，其点赞的`视频id`为 member，排序分数是`点赞时间`。

`视频Id`作为有序集合的 key，对其点赞的`用户id`为 member，排序分数为`点赞时间`。

##### Sociality

`用户id`` ``+`` ``consts.RedisFollowSuffix`为集合的 key，用户关注的人为集合元素

`用户id`` ``+`` ``consts.RedisFollowerSuffix` 为集合的 key，用户的分数为集合元素

#### 消息队列设计

所有打入 MySQL 的数据我们都会先发布至消息队列中，每一个服务都会有一个自己的订阅者协程，持续获取消息队列中的内容。这样可以避免流量过大时对 MySQL 造成冲击。

### 项目代码介绍

#### 项目代码结构

##### 主要结构

```Bash
├── docker-compose.yaml
├── otel-collector-config.yaml
├── go.mod
├── go.sum
├── server
│   ├── cmd
│   │   ├── api
│   │   ├── chat
│   │   ├── interaction
│   │   ├── sociality
│   │   ├── user
│   │   └── video
│   ├── idl
│   │   ├── api.thrift
│   │   ├── base.thrift
│   │   ├── chat.thrift
│   │   ├── errno.thrift
│   │   ├── interaction.thrift
│   │   ├── sociality.thrift
│   │   ├── user.thrift
│   │   └── video.thrift
│   └── shared
│       ├── Makefile
│       ├── consts
│       ├── errno
│       ├── kitex_gen
│       ├── middleware
│       ├── test
│       └── tools
```

##### 微服务内部结构

> 以 user 服务为例

```Bash
├── Makefile
├── config
│   └── config.go
├── config.yaml
├── dao
│   ├── user.go
│   └── user_test.go
├── global
│   └── global.go
├── handler.go
├── initialize
│   ├── chat_service.go
│   ├── db.go
│   ├── flag.go
│   ├── logger.go
│   ├── nacos.go
│   ├── redis.go
│   └── sociality_service.go
├── kitex.yaml
├── main.go
├── model
│   └── user.go
└── pkg
    ├── chat.go
    ├── md5.go
    ├── pack.go
    ├── redis.go
    ├── redis_test.go
    └── sociality.go
```

#### 数据库

##### MySQL

###### 初始化

```Go
// InitDB to init database
func InitDB() {
   c := global.ServerConfig.MysqlInfo
   dsn := fmt.Sprintf(consts.MySqlDSN, c.User, c.Password, c.Host, c.Port, c.Name)
   newLogger := logger.New(
      logrus.NewWriter(), // io writer
      logger.Config{
         SlowThreshold: time.Second,   // Slow SQL Threshold
         LogLevel:      logger.Silent, // Log level
         Colorful:      true,          // Disable color printing
      },
   )

   // global mode
   var err error
   global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
      NamingStrategy: schema.NamingStrategy{
         SingularTable: true,
      },
      Logger: newLogger,
   })
   if err != nil {
      klog.Fatalf("init gorm failed: %s", err)
   }
   if err := global.DB.Use(tracing.NewPlugin()); err != nil {
      klog.Fatalf("use tracing plugin failed: %s", err)
   }
}
```

MySQL 使用到了 GORM 进行操作，因此我们需要通过 GORM 来初始化 MySQL。值得一提的是这里的日志使用到了 GORM 提供的 Opentelemetry 插件中的 Logrus 日志，在后文中会再次介绍。

###### 使用

> 以 user 为例

我们首先在微服务下的 `model` 层中建立好数据模型。

```Go
type User struct {
   ID              int64  `gorm:"primarykey"`
   Username        string `gorm:"type:varchar(33);unique;not null"`
   Password        string `gorm:"type:varchar(33);not null"`
   Avatar          string `gorm:"type:varchar;not null"`
   BackGroundImage string `gorm:"type:varchar;not null"`
   Signature       string `gorm:"type:varchar;not null"`
}
```

接着在微服务下的 `dao` 层实现增删改查，以创建用户为例。

```Go
// CreateUser creates a user.
func (u *User) CreateUser(user *model.User) error {
   err := u.db.Model(&model.User{}).
      Where(&model.User{Username: user.Username}).First(&model.User{}).Error
   if err == nil {
      return ErrUserExist
   } else if err != gorm.ErrRecordNotFound {
      return err
   }
   return u.db.Model(&model.User{}).Create(user).Error
}
```

##### Redis

###### 初始化

```Go
func InitRedis() {
   global.RedisClient = redis.NewClient(&redis.Options{
      Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
      Password: global.ServerConfig.RedisInfo.Password,
      DB:       consts.RedisSocialClientDB,
   })
}
```

若需要多个客户端可以在初始化时按需求定制。

###### 使用

当我们需要 Redis 完成哪些服务时我们可以先在 `handler.go` 中进行定义，这里以 user 为例。

```Go
// RedisManager defines the redis interface.
type RedisManager interface {
   GetUserById(ctx context.Context, uid int64) (*model.User, error)
   BatchGetUserById(ctx context.Context, uidList []int64) ([]*model.User, error)
   CreateUser(ctx context.Context, user *model.User) error
}
```

接着在微服务下的 `pkg` 层中我们可以对这些接口进行实现，这里以创建用户为例。

```Go
// CreateUser creates a user.
func (r *RedisManager) CreateUser(ctx context.Context, user *model.User) error {
   uidStr := fmt.Sprintf("%d", user.ID)
   exists, err := r.redisClient.HExists(ctx, uidStr, consts.UsernameFiled).Result()
   if err != nil {
      return err
   }
   if exists {
      return errno.UserServerErr.WithMessage("user already exists")
   }
   batchData := make(map[string]string)
   batchData[consts.UsernameFiled] = user.Username
   batchData[consts.CryptPwdFiled] = user.Password
   return r.redisClient.HMSet(ctx, uidStr, batchData).Err()
}
```

##### Minio

###### 初始化

```Go
func initMinio() {
   mi := global.ServerConfig.UploadServiceInfo.MinioInfo
   // Initialize minio client object.
   mc, err := minio.New(mi.Endpoint, &minio.Options{
      Creds:  credentials.NewStaticV4(mi.AccessKeyID, mi.SecretAccessKey, ""),
      Secure: false,
   })
   if err != nil {
      klog.Fatalf("create minio client err: %s", err.Error())
   }
   exists, err := mc.BucketExists(context.Background(), mi.Bucket)
   if err != nil {
      klog.Fatal(err)
   }
   if !exists {
      err = mc.MakeBucket(context.Background(), mi.Bucket, minio.MakeBucketOptions{Region: "cn-north-1"})
      if err != nil {
         klog.Fatalf("make bucket err: %s", err.Error())
      }
   }
   policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + mi.Bucket + `/*"],"Sid": ""}]}`
   err = mc.SetBucketPolicy(context.Background(), mi.Bucket, policy)
   if err != nil {
      klog.Fatal("set bucket policy err:%s", err)
   }
   minioClient = mc
}
```

###### 使用

结合消息队列，实现异步上传视频和封面

```Go
func (s *Service) RunVideoUpload() error {
   taskCh, cleanUp, err := s.subscriber.Subscribe(context.Background())
   defer cleanUp()
   if err != nil {
      klog.Error("cannot subscribe", err)
      return err
   }
   for task := range taskCh {
      if err = getVideoCover(task.VideoTmpPath, task.CoverTmpPath); err != nil {
         klog.Errorf("get video cover err: videoTmpPath = %s", task.VideoTmpPath)
         continue
      }
      suffix, err := getFileSuffix(task.VideoTmpPath)
      if err != nil {
         klog.Errorf("get video suffix err:videoTmpPath = %s", task.VideoTmpPath)
         continue
      }
      buckName := s.config.MinioInfo.Bucket

      if _, err = s.minioClient.FPutObject(context.Background(), buckName, task.CoverUploadPath, task.CoverTmpPath, minio.PutObjectOptions{
         ContentType: "image/png",
      }); err != nil {
         klog.Error("upload cover image err", err)
         continue
      }
      _ = os.Remove(task.CoverTmpPath)
      if _, err = s.minioClient.FPutObject(context.Background(), buckName, task.VideoUploadPath, task.VideoTmpPath, minio.PutObjectOptions{
         ContentType: fmt.Sprintf("video/%s", suffix),
      }); err != nil {
         klog.Error("upload video err", err)
         continue
      }
      _ = os.Remove(task.VideoTmpPath)
   }
   return nil
}
```

#### 中间件

##### RabbitMQ

###### 初始化

```Go
// InitMq to init rabbitMQ
func InitMq() {
   c := global.ServerConfig.RabbitMqInfo
   amqpConn, err := amqp.Dial(fmt.Sprintf(consts.RabbitMqURI, c.User, c.Password, c.Host, c.Port))
   if err != nil {
      klog.Fatal("cannot dial amqp", err)
   }
   global.AmqpConn = amqpConn
}
```

###### 使用

在 `handler.go` 中定义好 Publish 的接口进行使用。

```Go
// Publisher defines the publisher interface.
type Publisher interface {
   Publish(context.Context, *sociality.DouyinRelationActionRequest) error
}
```

在 `pkg` 中的 `amqp.go` 中进行实现

```Go
// Publish publishes a message.
func (p *Publisher) Publish(_ context.Context, car *sociality.DouyinRelationActionRequest) error {
   body, err := sonic.Marshal(car)
   if err != nil {
      return fmt.Errorf("cannot marshal: %v", err)
   }

   return p.ch.Publish(p.exchange, "", false, false, amqp.Publishing{
      Body: body,
   })
}
```

在 `main.go` 中会开启一个协程对消息进行消费，其中消费的逻辑需要自行定义，这里以 sociality 为例。

```Go
func SubscribeRoutine(subscriber *Subscriber, dao *dao.Follow) {
   msgs, cleanUp, err := subscriber.Subscribe(context.Background())
   defer cleanUp()
   if err != nil {
      klog.Error("cannot subscribe", err)
   }
   for m := range msgs {
      fr, err := dao.FindRecord(m.ToUserId, m.UserId)
      if err == nil && fr == nil {
         err = dao.CreateFollow(&model.Follow{
            UserId:     m.ToUserId,
            FollowerId: m.UserId,
            ActionType: m.ActionType,
         })
         if err != nil {
            klog.Error("follow action error", err)
         }
      }
      if err != nil {
         klog.Error("follow error", err)
      }
      err = dao.UpdateFollow(m.ToUserId, m.UserId, m.ActionType)
      if err != nil {
         klog.Error("follow error", err)
      }
   }
}
```

##### Logger

```Go
// InitLogger to init logrus
func InitLogger() {
   // Customizable output directory.
   logFilePath := consts.KlogFilePath
   if err := os.MkdirAll(logFilePath, 0o777); err != nil {
      panic(err)
   }

   // Set filename to date
   logFileName := time.Now().Format("2006-01-02") + ".log"
   fileName := path.Join(logFilePath, logFileName)
   if _, err := os.Stat(fileName); err != nil {
      if _, err := os.Create(fileName); err != nil {
         panic(err)
      }
   }

   logger := kitexlogrus.NewLogger()
   // Provides compression and deletion
   lumberjackLogger := &lumberjack.Logger{
      Filename:   fileName,
      MaxSize:    20,   // A file can be up to 20M.
      MaxBackups: 5,    // Save up to 5 files at the same time.
      MaxAge:     10,   // A file can exist for a maximum of 10 days.
      Compress:   true, // Compress with gzip.
   }

   if runtime.GOOS == "linux" {
      logger.SetOutput(lumberjackLogger)
   }
   logger.SetLevel(klog.LevelDebug)

   klog.SetLogger(logger)
}
```

日志使用 Hertz/Kitex Opentelemetry 拓展中的 Logrus 日志库。当系统为 Linux 也就是线上环境时会重定向日志的输出，使用 Lumberjack 库对日志进行压缩与定期删除。当开发环境时会直接打印在控制台，方便 Debug。

##### Gzip

```Go
gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".mp4", ".png"})),
```

使用 Gzip 中间件资源进行压缩，并自定义不进行压缩的资源格式。

##### Pprof

使用 Pprof 中间件对项目进行测试。

```Go
pprof.Register(h)
```

使用以下命令来通过 Pprof 进行性能分析。

```Bash
go tool pprof -http=:8001 http://127.0.0.1:8080/debug/pprof/profile
```

![img](./img/fire1.png)
![img](./img/fire2.png)

可以看到优化后的火焰图性能更好，服务调用时间更短。

#### 服务治理

##### Nacos

Nacos 会同时承担服务中心与配置中心两种功能，以节约线上资源。

###### 初始化

```Go
// InitNacos to init nacos
func InitNacos(Port int) (registry.Registry, *registry.Info) {
   v := viper.New()
   v.SetConfigFile(consts.UserConfigPath)
   // ...

   configClient, err := clients.CreateConfigClient(map[string]interface{}{
      "serverConfigs": sc,
      "clientConfig":  cc,
   })
   // ...
   content, err := configClient.GetConfig(vo.ConfigParam{
      DataId: global.NacosConfig.DataId,
      Group:  global.NacosConfig.Group,
   })
   // ...
   err = sonic.Unmarshal([]byte(content), &global.ServerConfig)
   if err != nil {
      klog.Fatalf("nacos config failed: %s", err)
   }
   // ...
   registryClient, err := clients.NewNamingClient(
      vo.NacosClientParam{
         ClientConfig:  &cc,
         ServerConfigs: sc,
      },
   )
   // ...

   r := nacos.NewNacosRegistry(registryClient, nacos.WithGroup(consts.UserGroup))

   // ...
   return r, info
}
```

由于代码冗长，这里只提供关键代码，我们先通过 Viper 对 Nacos 进行配置，并初始化配置中心，接着进行服务中心的初始化，进行服务注册。

Kitex 与 Hertz 在优雅推迟时会自动进行服务取消注册。服务发现请关注下文 RPC 部分。

##### Opentelemetry

OpenTelemetry 要解决的是对可观测性的大一统，在我们的项目中，Trace 使用到的是 Jaeger，Metrics 使用到了 Prometheus，Logs 使用的是 Logrus（在 GORM 日志中配置的相同日志库）。

```Go
p := provider.NewOpenTelemetryProvider(
   provider.WithServiceName(global.ServerConfig.Name),
   provider.WithExportEndpoint(global.ServerConfig.OtelInfo.EndPoint),
   provider.WithInsecure(),
)
defer p.Shutdown(context.Background())
```

#### 安全

##### ErrNo

在项目中使用 ErrNo 来提供更多的错误信息但不会将系统内部的错误信息返回给前端。其中错误码在 IDL 中就进行了定义。

```Thrift
namespace go errno

enum Err {
    Success              = 0,
    ParamsErr            = 1,
    ServiceErr           = 2,
    RPCInteractionErr    = 10000,
    InteractionServerErr = 10001,
    RPCSocialityErr      = 20000,
    SocialityServerErr   = 20001,
    RPCUserErr           = 30000,
    UserServerErr        = 30001,
    UserAlreadyExistErr  = 30002,
    UserNotFoundErr      = 30003,
    AuthorizeFailErr     = 30004,
    RPCVideoErr          = 40000,
    VideoServerErr       = 40001,
    RPCChatErr           = 50000,
    ChatServerErr        = 50001,
}
```

同时错误信息也在 `shared/errno` 中进行了定义，这里以 `Success` 为例。

```Go
Success = NewErrNo(int64(errno.Err_Success), "success")
```

除此之外你也可以自定义错误信息，以 `SentMessage` 方法为例。

```Go
func (s *ChatServiceImpl) SentMessage(ctx context.Context, req *chat.DouyinMessageActionRequest) (resp *chat.DouyinMessageActionResponse, err error) {
    // ...
    if err != nil {
        klog.Error("publish message error", err)
        resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("sent message error"))
        return resp, nil
    }
    // ...
}
```

##### JWT

秘钥从配置中心中获得，未出现在代码中，实现了脱敏。在用户登录成功或注册成功时会生成 Token，并且 Token 中加入了用户的一些信息。

```Go
resp.UserId = usr.ID
resp.Token, err = s.Jwt.CreateToken(models.CustomClaims{
   ID: usr.ID,
   StandardClaims: jwt.StandardClaims{
      NotBefore: time.Now().Unix(),
      ExpiresAt: time.Now().Unix() + consts.ThirtyDays,
      Issuer:    consts.JWTIssuer,
   },
})
```

在网关层中我们会使用到 JWTAuth 中间件，对传入的 Token 进行校验。

```Go
func _publishMw() []app.HandlerFunc {
   return []app.HandlerFunc{
      middleware.JWTAuth(global.ServerConfig.JWTInfo.SigningKey),
   }
}
```

##### MD5

当我们进行用户注册时，密码不会进行明文存储，会先对密码进行 MD5 加盐加密。

```Go
// Md5Crypt uses MD5 encryption algorithm to add salt encryption.
func Md5Crypt(str string, salt ...interface{}) (CryptStr string) {
   if l := len(salt); l > 0 {
      slice := make([]string, l+1)
      str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
   }
   return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
```

在注册阶段会进行使用。

```Go
// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
   // ...
   usr := &model.User{
      ID:       sf.Generate().Int64(),
      Username: req.Username,
      Password: pkg.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt), // Encrypt password with md5.
   }
   // ...
}
```

在后续登录时会将用户输入的密码再次进行加密，将加密后的数据与数据库中的数据进行比对，若相同则说明密码正确，反之密码错误。

##### Limiter

使用 Limiter 中间件对项目进行限流

```Go
limiter.AdaptiveLimit(limiter.WithCPUThreshold(900)),
```

- 当CPU负载小于 90% 时：当前时间距离上次触发限流小于1s，则判断当前最大请求数是否大于过去最大负载情况，如果大于负载情况，然后限制流量。
- 当CPU负载大于 90% 时：判断当前请求数是否大于过去的最大负载，如果大于过去的最大负载，则进行限流。

#### 其他

##### RPC

当一个微服务需要调用别的微服务时需要进行 RPC 调用，在 `pkg` 中我们对需要的服务进行初始化。这里以 user 服务需要调用 chat 服务为例。

###### 初始化

```Go
// InitChat init chat service.
func InitChat() {
   // init resolver
   // Read configuration information from nacos
   sc := []constant.ServerConfig{
      {
         IpAddr: global.NacosConfig.Host,
         Port:   global.NacosConfig.Port,
      },
   }

   cc := constant.ClientConfig{
      // ...
   }

   nacosCli, err := clients.NewNamingClient(
      vo.NacosClientParam{
         ClientConfig:  &cc,
         ServerConfigs: sc,
      })
   r := nacos.NewNacosResolver(nacosCli, nacos.WithGroup(consts.ChatGroup))
   if err != nil {
      klog.Fatalf("new nacos client failed: %s", err.Error())
   }
   provider.NewOpenTelemetryProvider(
      // ...
   )

   // create a new client
   c, err := chat.NewClient(
      global.ServerConfig.ChatSrvInfo.Name,
      client.WithResolver(r),                                     // service discovery
      client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
      client.WithMuxConnection(1),                                // multiplexing
      client.WithSuite(tracing.NewClientSuite()),
      client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.ServerConfig.ChatSrvInfo.Name}),
   )
   if err != nil {
      klog.Fatalf("ERROR: cannot init client: %v\n", err)
   }
   global.ChatClient = c
}
```

其中比较重要的是我们会在这里进行服务发现去找到我们已经注册的服务，并且使用加权轮询算法的负载均衡。

###### 使用

当我们需要使用别的服务时，我们需要在 `handler.go` 中定义好相关服务的接口作为我们的防腐层，不在目前服务的逻辑中出现调用逻辑，直接使用接口进行操作，以 user 调用 interaction 中的服务为例。

```Go
type InteractionManager interface {
   GetInteractInfo(ctx context.Context, userId int64) (*base.UserInteractInfo, error)
   BatchGetInteractInfo(ctx context.Context, userIdList []int64) ([]*base.UserInteractInfo, error)
}
```

接着我们需要在 `pkg` 层中对接口进行实现，以 `GetInteractionInfo` 方法为例。

```Go
func (i *InteractionManager) GetInteractInfo(ctx context.Context, userId int64) (*base.UserInteractInfo, error) {
   resp, err := i.client.GetUserInteractInfo(ctx, &interaction.DouyinGetUserInteractInfoRequest{
      UserId: userId,
   })
   if err != nil {
      return nil, err
   }
   if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
      return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
   }
   return resp.InteractInfo, nil
}
```

##### Unit Test

我们对每一个服务的数据操作都进行了测试，为了让测试的数据库不影响业务数据库，我们选择了使用 Docker 容器进行单元测试。首先我们在 Docker 中启动一个 MySQL 或 Redis 的容器，接着在此容器中对数据库进行初始化，接着就可以进行测试了。在测试结束后会自动删除掉此容器，防止占用空间。

下面以测试 user 的 MySQL 操作为例，首先我们需要在 Docker 中运行 MySQL 数据库。

```Go
// RunWithMySQLInDocker runs the tests with
// a MySQL instance in a docker container.
func RunWithMySQLInDocker(t *testing.T) (cleanUpFunc func(), db *gorm.DB, err error) {
  // ...

   ctx := context.Background()
   resp, err := c.ContainerCreate(ctx, &container.Config{
      // ...
   }, nil, nil, "")
   if err != nil {
      return func() {}, nil, err
   }
   containerID := resp.ID
   cleanUpFunc = func() {
      err = c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
         Force: true,
      })
      if err != nil {
         t.Error("remove test docker failed", err)
      }
   }

   err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
   if err != nil {
      return cleanUpFunc, nil, err
   }

   inspRes, err := c.ContainerInspect(ctx, containerID)
   if err != nil {
      return cleanUpFunc, nil, err
   }
   hostPort := inspRes.NetworkSettings.Ports[consts.MySQLContainerPort][0]
   port, _ := strconv.Atoi(hostPort.HostPort)
   mysqlDSN := fmt.Sprintf(consts.MySqlDSN, consts.MySQLAdmin, consts.DockerTestMySQLPwd, hostPort.HostIP, port, consts.TikGok)
   // Init mysql
   time.Sleep(10 * time.Second)
   db, err = gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{
      // ...
   })
   // ...
}
```

这里我忽略掉了一些无用的代码，大概逻辑就是通过调用 API 在 Docker 中新建一个 MySQL 的容器，接着在 `defer` 中进行 `ContainerRemove` 的操作。最后就实现了我们想要的结果。

在测试中我们大量使用表格驱动测试以测试多种不同的情况。

```Go
func TestUserLifecycle(t *testing.T) {
   cleanUp, db, err := test.RunWithMySQLInDocker(t)
   defer cleanUp()
   // ...

   dao := NewUser(db)
   // ...

   cases := []struct {
      name       string
      op         func() (string, error)
      wantErr    bool
      wantResult string
   }{
      // ...
    },
}

   for _, cc := range cases {
      result, err := cc.op()
      if cc.wantErr {
         if err == nil {
            t.Errorf("%s:want error;got none", cc.name)
         } else {
            continue
         }
      }
      if err != nil {
         t.Errorf("%s:operation failed: %v", cc.name, err)
      }
      if result != cc.wantResult {
         t.Errorf("%s:result err: want %s,got %s", cc.name, cc.wantResult, result)
      }
   }
}
```

## 测试结果

### 单元测试

我们对项目的数据操作，例如在 MySQL 中的增删改查和在 Redis 中的键值对操作都进行了单元测试，其中测试方法在上文也提到过，通过新建一个 Docker 容器来进行测试。测试结果是全部通过，在 CI 中也有所体现。

![img](https://jxi4fut4kr.feishu.cn/space/api/box/stream/download/asynccode/?code=YTVhY2M1NGE2OTliNWU2Yzk0OGI4YTBhMDQyZjgxMDNfVFFGRUk5d3ljN0diejFzSXJreTZnY3JIQTJUUjRGTUdfVG9rZW46Ym94Y24wM1JSOGNwN0hMdDR2dEc0ZFROUGVnXzE2NzY4OTcyNTg6MTY3NjkwMDg1OF9WNA)

![img](https://jxi4fut4kr.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTc4MDQ3OGM0ZDUxNzMzMjdjYzgzZmM3MjQ2YTdlZmJfNURWSGhXMGJzYnRkeU9ZU0EzSmR3aWRrOUF5ZExmZUxfVG9rZW46Ym94Y25SWHQwRzFyVWZaZW5DaG1Qb0tNdkFkXzE2NzY4OTcyNTg6MTY3NjkwMDg1OF9WNA)

### 压力测试

#### 压测环境

| CPU         | 内存 |
| ----------- | ---- |
| M1 Pro 10核 | 16G  |

#### 压测结果

**优化前**

![img](./img/pressure_test2.jpg)

**优化后**

![img](./img/pressure_test1.jpg)

### 接口测试

接口测试以及测试结果我们均保存在了 Postman 中，请访问一下地址查看详细测试内容。

> https://documenter.getpostman.com/view/20584759/2s93CHuuiQ

## 许可证

TikGok 在 GNU General Public 许可证 3.0 版下开源。