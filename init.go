package core

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/grestful/logs"
	"github.com/grestful/utils"
	"github.com/jinzhu/gorm"
	"os"
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
		Log:    log.Global,
		Config: &goconfig.ConfigFile{},
		Db:     make(map[string]*gorm.DB),
		Redis:  make(map[string]*redis.Client),
	}
	contextPool = &sync.Pool{New: func() interface{} {
		return &Context{}
	}}
}

func initLog() {
	log.Project = ServiceName
	logConfig := gin.LoggerConfig{
		Formatter: gin.LogFormatter(func(param gin.LogFormatterParams) string {
			lg := &log.LogRecord{
				Level:   1,
				Created: param.TimeStamp,
				Source:  "",
				Message: fmt.Sprintf("ip: %s, method: %s, path: %s, code: %d, agent: %s, error %s",
					param.ClientIP,
					param.Method,
					param.Path,
					param.StatusCode,
					param.Request.UserAgent(),
					param.ErrorMessage),
				Category: "default",
			}
			// fmt.Fprintln(os.Stdout, log.FormatLogRecord(log.FORMAT, lg))
			return log.FormatLogRecord(log.FORMAT, lg)
		}),
		Output:    os.Stdout,
		SkipPaths: nil,
	}
	log.Project = ServiceName
	typ, err := gGore.Config.GetValue("log", "type")
	if typ == "" || err != nil {
		typ = "console"
	}
	switch typ {
	case "file":
		path, err := gGore.Config.GetValue("log", "path")
		if path == "" || err != nil {
			typ = "console"
		}
		level, err := gGore.Config.GetValue("log", "level")
		if level == "" || err != nil {
			level = "DEBUG"
		}
		log.SetFile(log.FileConfig{
			Enable:   true,
			Category: "file",
			Level:    level,
			Filename: path+"/cast.log",
			Rotate:   true,
			Daily:    true,
			Sanitize: false,
		})
		writer := log.GetLogger("file")
		log.SetDefaultLog(writer)
		logConfig.Output = writer
		GetCore().Gin.Use(gin.LoggerWithConfig(logConfig))
	case "conn":
		proto, err := gGore.Config.GetValue("log", "net")
		if typ == "" || err != nil {
			typ = "console"
		}
		addr, err := gGore.Config.GetValue("log", "addr")
		if typ == "" || err != nil {
			typ = "console"
		}
		level, err := gGore.Config.GetValue("log", "level")
		if level == "" || err != nil {
			level = "console"
		}
		conn := log.SocketConfig{
			Enable:   true,
			Category: "SOCKET",
			Level:    "DEBUG",
			Addr:     addr,
			Protocol: proto,
		}
		log.SetConn(conn)
		writer := log.GetLogger("socket")
		logConfig.Output = writer
		log.SetDefaultLog(writer)
		GetCore().Gin.Use(gin.LoggerWithConfig(logConfig))
	case "console":
		fallthrough
	default:
		log.SetConsole(log.ConsoleConfig{
			Enable: true,
			Level:  "DEBUG",
		})
		writer := log.GetLogger("stdout")
		logConfig.Output = writer
		log.SetDefaultLog(writer)
		GetCore().Gin.Use(gin.LoggerWithConfig(logConfig))
	}
}

func InitConfig(path string) {
	configOne.Do(func() {
		var err error
		GetCore().Config, err = goconfig.LoadConfigFile(path)
		if err != nil {
			panic(fmt.Sprintf("无法加载配置文件：%s \n", err))
		}
		ServiceName, _ = gGore.Config.GetValue("", "SERVICE_NAME")
		initLog()

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

	if err != nil {
		log.Error("db init error: %s", err.Error())
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

	if err != nil {
		log.Error("redis init error: %s", err.Error())
	}
}
