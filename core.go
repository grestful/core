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

// return instance config
func GetConfig() *goconfig.ConfigFile {
	return GetCore().Config
}

// return instance
func GetCore() *Core {
	return gGore
}

// return instance gin core
func GetGin() *gin.Engine {
	return GetCore().Gin
}

// return instance logger
func GetLog() log.Logger {
	return GetCore().Log
}

// get db if exists
func GetDb(name string) *gorm.DB {
	if name == "" {
		name = "default"
	}
	return GetCore().GetDb(name)
}

//return instance session
func (gG *Core) GetSessionType() string {
	return gG.SessionType
}

// alias return  instance gin core
func GetRouter() IRouter {
	return GetCore().Gin
}

// alias gin group
func (gG *Core) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return gG.Gin.Group(path, handlers...)
}

// alias gin Use
func (gG *Core) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return gG.Gin.Use(middleware...)
}

// alias gin Handle
func (gG *Core) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return gG.Gin.Handle(httpMethod, relativePath, handlers...)
}

// get default cache, if config set MYSQL_DSN
func (gG *Core) GetDefaultCache() ICache {
	if db, ok := gG.Cache["default"]; ok {
		return db
	}
	return nil
}

// get cache with name
// if not exists, fatal err
func (gG *Core) GetCache(name string) ICache {
	if db, ok := gG.Cache[name]; ok {
		return db
	}

	panic(name + " cache not exists")
}

// get default db, if config set MYSQL_DSN
func (gG *Core) GetDefaultDb() *gorm.DB {
	if db, ok := gG.Db["default"]; ok {
		return db
	}
	return nil
}

// get cache with name
// now only support gorm
// if not exists, fatal err
func (gG *Core) GetDb(name string) *gorm.DB {
	if db, ok := gG.Db[name]; ok {
		return db
	}

	panic(name + " db not exists")
}

// get default redis
func (gG *Core) GetDefaultRedis() *redis.Client {
	if c, ok := gG.Redis["default"]; ok {
		return c
	}
	return nil
}

// get cache with name
// now only support gorm
func (gG *Core) GetRedis(name string) *redis.Client {
	if c, ok := gG.Redis[name]; ok {
		return c
	}

	return nil
}

// set logger format
// params:
// %A - Time (2006-01-02T15:04:05.000Z)  means all
// %T - Time (15:04:05 MST)
// %t - Time (15:04)
// %D - Date (2006/01/02)
// %d - Date (01/02/06)
// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
// %S - Source
// %M - Message
// Ignores unknown formats
// Recommended: "[%A] [%L] (%S) %M"
func (gG *Core) SetLoggerFormat(format string) {
	gG.Log.GetDefaultFilter().SetFormat(format)
}
