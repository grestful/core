package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grestful/session"
	"github.com/grestful/utils"
	"net/http"
)

type Context struct {
	*gin.Context
}

type ProcessFunc func(process *Controller) IError
type Controller struct {
	TrackId    string //访问id
	Session    session.IUserSession
	Ctx        *Context
	Code       string
	error      IError
	Data       interface{} //output  from data raw data parse
	Req        interface{} //input params
	ProcessFun ProcessFunc
}

func Session(process Controller) error {
	typ := gGore.GetSessionType()

	sid := ""
	cook,err := process.GetContext().Request.Cookie("sid")
	if err != nil {
		return err
	}
	if cook.Value == "" {
		sid = process.GetContext().Request.URL.Query().Get("sid")
	}else{
		sid = cook.Value
	}

	if sid == "" {
		return errors.New("sid is empty")
	}

	maxLeftStr, err := gGore.Config.GetValue("session", "MAX_LEFT_TIME")
	if err != nil || maxLeftStr == "" {
		maxLeftStr = "3600"
	}
	maxLeft := utils.String2Int64(maxLeftStr, 3600)
	switch typ {
	case SessionTypeFile:
		path, err := gGore.Config.GetValue("session", "FILE_PATH")
		if err != nil || path == "" {
			path = "/tmp"
		}

		sess := session.GetNewFileSession(path, maxLeft)
		process.Session = session.GetNewUserSession(sid, sess)
		return nil
	case SessionTypeRedis:
		name, err := gGore.Config.GetValue("session", "REDIS_CONNECTION")
		if err != nil || name == "" {
			name = "default"
		}
		c := gGore.GetRedis(name)
		if c == nil {
			return errors.New(fmt.Sprintf("can't find redis connection name %s" , name))
		}
		maxLeft := utils.String2Int64(maxLeftStr, 3600)
		sess := session.GetNewRedisSession(c, maxLeft)
		process.Session = session.GetNewUserSession(sid, sess)
		return nil
	}

	return errors.New("session type error")
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
