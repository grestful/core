package core

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grestful/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var gGore *Core
var configOne sync.Once
var contextPool *sync.Pool

func GetContext(g *gin.Context) *Context {
	c := contextPool.Get().(*Context)
	c.Context = g
	return c
}

func init() {
	gGore = &Core{
		Gin:    gin.New(),
		Log:    log.New(),
		Config: &goconfig.ConfigFile{},
		Db:     make(map[string]*gorm.DB),
		Redis:  make(map[string]*redis.Client),
	}
	contextPool = &sync.Pool{New: func() interface{} {
		return &Context{}
	}}
}

func InitConfig(path string) {
	configOne.Do(func() {
		var err error
		GetCore().Config, err = goconfig.LoadConfigFile(path)
		if err != nil {
			panic(fmt.Sprintf("无法加载配置文件：%s \n", err))
		}

		initDb()

		initRedis()

		initSessionType()
	})
}

func initSessionType() {
	cf, err := gGore.Config.GetValue("", "SESSION_TYPE")
	if err == nil {
		gGore.SessionType = cf
	}
}

func initDb() {
	cf, err := gGore.Config.GetValue("", "MYSQL_DSN")
	if err == nil {
		db, err := gorm.Open("mysql", cf)
		if err == nil {
			GetCore().Db["default"] = db
		}
	}
}

func initRedis() {
	cf, err := GetCore().Config.GetSection("redis")
	if err == nil {
		c := redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         cf["host"] + ":" + cf["port"],
			Dialer:       nil,
			OnConnect:    nil,
			Password:     cf["auth"],
			DB:           utils.String2Int(cf["db"], 0),
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     20,
			MinIdleConns: 5,
			MaxConnAge:   20,
			PoolTimeout:  3 * time.Second,
			IdleTimeout:  5 * time.Second,
		})
		if c.Ping().Err() == nil {
			GetCore().Redis["default"] = c
		}
	}
}
