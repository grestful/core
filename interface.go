package core

import "time"

// It's a controller with a request process
type IController interface {

	//set context(gin)
	SetContext(ctx *Context)
	//get context(gin)
	GetContext() *Context
	//self decode request obj
	//example:
	//get: ?a=b&c=d  => struct{A:b string `url:"a"`,C:d string `url:"c"`}
	//post(json): {"a":"b","c":"d"} struct{A:b string `url:"a"`,C:d string `url:"c"`}
	Decode() IError

	// processing business
	Process() IError

	//defer set error
	SetError(err IError)

	//get real response
	getResponse() Response

	// Used to track distributed service link logs
	// recommend set it in query string
	GetTrackId() string

	// set trackId
	SetTrackId(id string)
}

type IError interface {
	GetCode() string
	GetMsg() string
	Error() string
	IsNil() bool
	GetDetail() interface{}
}

type IResponse interface {
	GetBytes() []byte
	GetString() string
}

type ICache interface {
	GetBytes(key string) []byte
	GetString(key string) string
	GetMap(key string) map[string]string
	GetList(key string) []string
	GetValue(key string) (string, error)

	SetValue(key string, val interface{}, ex time.Duration) error
	SetList(key string, val ...interface{}) error
	SetString(key, val string, ex time.Duration) error
	SetMap(key string, m map[string]string, ex time.Duration) error
	SetBytes(key string, b []byte, ex time.Duration) error
	Expire(key string, ex time.Duration)

	Command(args ...string) error
}

type IModel interface {
	GetBytes() []byte
	GetString() string

	Update() bool
	Create() int64
	Delete() bool
	One(sql string, args ...interface{}) IGetterSetter
	List(sql string, args ...interface{}) []IGetterSetter

	TableName() string
	GetKeyName() string
	GetKey() int64
	GetName() string
}

type IGetterSetter interface {
	GetAttribute(key string) interface{}
	SetAttribute(key string, value interface{}) bool
	SetAttributes(mp map[string]interface{}) bool
	GetAttributes() map[string]interface{}

	GetAttrInt(key string) (int, error)
	GetAttrInt64(key string) (int64, error)
	GetAttrFloat(key string) (float32, error)
	GetAttrFloat64(key string) (float64, error)
	GetAttrUInt(key string) (uint, error)
	GetAttrUInt64(key string) (uint64, error)
	GetAttrBool(key string) (bool, error)
	GetAttrString(key string) (string, error)

	GetAttrTime(key string) (time.Time, error)
}

type IModelStruct interface {
	TableName() string
	GetKeyName() string
}

type IRouter interface {
}
