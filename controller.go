package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grestful/utils"
	"net/http"
)

type Context struct {
	*gin.Context
}

type ProcessFunc func(process *Controller) IError
type Controller struct {
	TrackId    string //访问id
	Session    IUserSession
	Ctx        *Context
	Code       string
	error      IError
	Data       interface{} //output  from data raw data parse
	Req        interface{} //input params
	ProcessFun ProcessFunc
}

func RunProcess(process IController, g *gin.Context) {
	c := GetContext(g)
	process.SetContext(c)
	var ierr IError
	defer func() {
		if ierr != nil {
			process.SetError(ierr)
		}

		var err error
		c.Status(200)
		c.Header("Content-Type", "application/json; charset=utf-8")
		defer func() {
			if err != nil {
				errStr := fmt.Sprintf(`{"code":"%s","msg":"system error: %s","data":null}`, FailInternal, err.Error())
				rStr(c, errStr)
			}
		}()

		res := process.getResponse()

		rJson(c, res)
	}()

	if ierr = getTrackId(process); ierr != nil {
		return
	}
	if ierr = runProcess(process); ierr != nil {
		return
	}
}

func getTrackId(process IController) IError {
	trackId := utils.GetRequestKey(process.GetContext().Request, "track_id")
	if trackId == "" {
		return NewErrorStr("need params track_id")
	}
	process.SetTrackId(trackId)
	return nil
}

//run process
func runProcess(process IController) (err IError) {
	err = process.Decode()
	if err != nil {
		process.SetError(err)
		return
	}

	err = process.Process()
	if err != nil {
		return
	}

	return
}

func (process *Controller) Decode() IError {
	process.Data = nil

	switch process.Ctx.Context.Request.Method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		if process.Ctx.Context.Request.Header.Get("Content-type") == "application/json" {
			bt, err := process.Ctx.GetRawData()
			if err == nil {
				if len(bt) == 0 {
					bt = []byte("{}")
				}
				if err := json.Unmarshal(bt, &process.Req); err != nil {
					return NewErrorCode(FailJsonParse)
				}
			}
		}
	default:

	}

	return nil
}

func (process *Controller) Process() IError {
	if process.ProcessFun != nil {
		return process.ProcessFun(process)
	}
	return nil
}

func (process *Controller) getResponse() Response {
	if process.error != nil {
		return getDefaultErrorResponse(process.error)
	} else {
		return getResponseWithCode(process.Code, process.Data)
	}
}

func (process *Controller) SetError(err IError) {
	GetLog().Infof("set error in %s,track_id: %s, err: %v\n", process.Ctx.Context.Request.URL.Path, process.TrackId, err)
	process.error = err
	return
}

func (process *Controller) GetTrackId() string {
	return process.TrackId
}

func (process *Controller) SetTrackId(id string) {
	process.TrackId = id
}

func (process *Controller) SetContext(c *Context) {
	if process.Ctx == nil {
		process.Ctx = c
	}
}

func (process *Controller) GetContext() *Context {
	return process.Ctx
}
