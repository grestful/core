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