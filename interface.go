package core

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