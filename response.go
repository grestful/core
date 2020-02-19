package core

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code string      `json:"retcode"`
	Msg  string      `json:"desc"`
	Data interface{} `json:"biz"`
}

func (res Response) GetBytes() []byte {
	b,_ := json.Marshal(res)
	return b
}

func (res Response) GetString() string {
	b,_ := json.Marshal(res)
	return string(b[:])
}

func getDefaultErrorResponse(err IError) Response {
	var data interface{}

	if err.Error() == "操作失败" {
		data = err.GetDetail()
	} else {
		data = err.Error()
	}

	return Response{
		err.GetCode(), err.GetMsg(), data,
	}
}

func getResponseWithCode(code string, data ...interface{}) Response {
	if code == "" {
		code = SuccessCode
	}
	msg := DefaultCodeMapping.GetCodeInfo(code)
	var r = Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	if data == nil {
		return r
	}

	l := len(data)
	if l > 0 {
		if l == 1 {
			r.Data = data[0]
		} else {
			r.Data = data
		}
	}

	return r
}

func rStr(c *Context, str string) {
	c.String(http.StatusOK, str)
}

func rJson(c *Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
