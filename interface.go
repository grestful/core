package core

import "time"

type IController interface {
	SetContext(ctx *Context)
	GetContext() *Context
	Decode() IError
	Process() IError
	SetError(err IError)
	getResponse() Response
	GetTrackId() string
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
	GetValue(key string, val *interface{}) error

	SetValue(key string, val interface{}, ex time.Duration) error
	SetList(key string,  val... interface{}) error
	SetString(key, val string, ex time.Duration) error
	SetMap(key string, m map[string]string, ex time.Duration)error
	SetBytes(key string, b []byte, ex time.Duration) error
	Expire(key string, ex time.Duration)

	Command(args... string) error
}

type IModel interface {
	GetBytes() []byte
	GetString() string

	Save() bool
	Create() bool
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
	GetAttrInt64(key string)(int64, error)
	GetAttrFloat(key string) (float32, error)
	GetAttrFloat64(key string)(float64, error)
	GetAttrUInt(key string) (uint, error)
	GetAttrUInt64(key string)(uint64, error)
	GetAttrBool(key string) (bool, error)
	GetAttrString(key string)(string, error)
}


type ISession interface {
	Close () bool
	Destroy(sid string)  bool
	Gc(maxLeftTime int64)  bool
	Open(savePath string)  bool
	Read(sid string) map[string]string
	Write(sid string, data map[string]string)  bool
	Error(sid string) error
}

type IRouter interface {

}