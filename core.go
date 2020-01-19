package core

import (
	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io"
)

type Core struct {
	Gin    *gin.Engine
	Log    *log.Logger
	Config *goconfig.ConfigFile
	Db     map[string]*gorm.DB
	Redis  map[string]*redis.Client
	//Cache
}

func GetCore() *Core {
	return gGore
}

func GetLog() *log.Logger {
	return GetCore().Log
}

func (gG *Core) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return gG.Gin.Group(path, handlers...)
}

func (gG *Core) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return gG.Gin.Use(middleware...)
}

func (gG *Core) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return gG.Gin.Handle(httpMethod, relativePath, handlers...)
}

//func (gG *Core) Db(name string) {
//	return
//}

func (gG *Core) GetDefaultDb() *gorm.DB {
	if db, ok := gG.Db["default"]; ok {
		return db
	}
	return nil
}

func (gG *Core) GetDb(name string) *gorm.DB {
	if db, ok := gG.Db[name]; ok {
		return db
	}

	panic(name + " db not exists")
}

func (gG *Core) GetDefaultRedis() *redis.Client {
	if c, ok := gG.Redis["redis"]; ok {
		return c
	}
	return nil
}

func (gG *Core) GetRedis(name string) *gorm.DB {
	if db, ok := gG.Db[name]; ok {
		return db
	}

	panic(name + " db not exists")
}

func (gG *Core) SetLoggerFormat(format log.Formatter) {
	gG.Log.SetFormatter(format)
}

func (gG *Core) SetLoggerLevel(level log.Level) {
	gG.Log.SetLevel(level)
}

func (gG *Core) SetOutput(out io.Writer) {
	gG.Log.SetOutput(out)
}

func (gG *Core) SetReportCaller(reportCaller bool) {
	gG.Log.SetReportCaller(reportCaller)
}
