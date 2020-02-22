package core

import (
	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/grestful/logs"
	"github.com/jinzhu/gorm"
)

var ServiceName string
type Core struct {
	Gin         *gin.Engine
	Log         log.Logger
	Config      *goconfig.ConfigFile
	Db          map[string]*gorm.DB
	Redis       map[string]*redis.Client
	Cache       map[string]ICache
	SessionType string
}

func GetConfig() *goconfig.ConfigFile  {
	return  GetCore().Config
}

func GetCore() *Core {
	return gGore
}

func GetLog() log.Logger {
	return GetCore().Log
}

func GetDb(name string) *gorm.DB {
	if name == "" {
		name = "default"
	}
	return GetCore().GetDb(name)
}

func (gG *Core) GetSessionType() string {
	return gG.SessionType
}

func GetRouter() IRouter {
	return GetCore().Gin
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
func (gG *Core) GetDefaultCache() ICache {
	if db, ok := gG.Cache["default"]; ok {
		return db
	}
	return nil
}

func (gG *Core) GetCache(name string) ICache {
	if db, ok := gG.Cache[name]; ok {
		return db
	}

	panic(name + " cache not exists")
}

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
	if c, ok := gG.Redis["default"]; ok {
		return c
	}
	return nil
}

func (gG *Core) GetRedis(name string) *redis.Client {
	if c, ok := gG.Redis[name]; ok {
		return c
	}

	return nil
}

func (gG *Core) SetLoggerFormat(format string) {
	gG.Log.GetDefaultFilter().SetFormat(format)
}

