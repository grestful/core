package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grestful/session"
	"github.com/grestful/utils"
	"net/http"
	"time"
)

type Context struct {
	*gin.Context
}

type ProcessFunc func(controller*Controller) IError
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

func Session(controller *Controller) error {
	typ := gGore.GetSessionType()

	sid := ""
	cookName,err := gGore.Config.GetValue("", "SESSION_NAME")
	if cookName == "" {
		cookName = "sid"
	}
	cook,err := controller.GetContext().Request.Cookie(cookName)
	if err != nil {
		return err
	}
	if cook.Value == "" {
		sid = controller.GetContext().Request.URL.Query().Get(cookName)
	}else{
		sid = cook.Value
	}

	if sid == "" {
		return errors.New("sid is empty")
	}

	maxLeftStr, err := gGore.Config.GetValue("session", "max_left_time")
	if err != nil || maxLeftStr == "" {
		maxLeftStr = "3600"
	}
	maxLeft := utils.String2Int64(maxLeftStr, 3600)
	switch typ {
	case SessionTypeFile:
		path, err := gGore.Config.GetValue("session", "file_path")
		if err != nil || path == "" {
			path = "/tmp"
		}

		sess := session.GetNewFileSession(path, maxLeft)
		controller.Session = session.GetNewUserSession(sid, sess)
		return nil
	case SessionTypeRedis:
		name, err := gGore.Config.GetValue("session", "redis_name")
		if err != nil || name == "" {
			name = "default"
		}
		c := gGore.GetRedis(name)
		if c == nil {
			return errors.New(fmt.Sprintf("can't find redis connection name %s" , name))
		}
		maxLeft := utils.String2Int64(maxLeftStr, 3600)
		sess := session.GetNewRedisSession(c, maxLeft)
		controller.Session = session.GetNewUserSession(sid, sess)
		return nil
	}

	return errors.New("session type error")
}

func RunProcess(controller IController, g *gin.Context) {
	c := GetContext(g)
	controller.SetContext(c)
	var ierr IError
	defer func() {
		if ierr != nil {
			controller.SetError(ierr)
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

		res := controller.getResponse()

		rJson(c, res)
	}()

	if ierr = getTrackId(controller); ierr != nil {
		return
	}
	if ierr = runProcess(controller); ierr != nil {
		return
	}
}

func getTrackId(controller IController) IError {
	trackId := utils.GetRequestKey(controller.GetContext().Request, "track_id")
	if trackId == "" {
		//return NewErrorStr("need params track_id")
		trackId = utils.Int642String(time.Now().Unix() * 1000)
	}
	controller.SetTrackId(trackId)
	return nil
}

//run controller process
func runProcess(controller IController) (err IError) {
	err = controller.Decode()
	if err != nil {
		controller.SetError(err)
		return
	}

	err = controller.Process()
	if err != nil {
		return
	}

	return
}

func (controller*Controller) Decode() IError {
	controller.Data = nil

	switch controller.Ctx.Context.Request.Method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		if controller.Ctx.Context.Request.Header.Get("Content-type") == "application/json" {
			bt, err := controller.Ctx.GetRawData()
			if err == nil {
				if len(bt) == 0 {
					bt = []byte("{}")
				}
				if err := json.Unmarshal(bt, &controller.Req); err != nil {
					return NewErrorCode(FailJsonParse)
				}
			}
		}
	default:

	}

	return nil
}

func (controller*Controller) Process() IError {
	if controller.ProcessFun != nil {
		return controller.ProcessFun(controller)
	}
	return nil
}

func (controller*Controller) getResponse() Response {
	if controller.error != nil {
		return getDefaultErrorResponse(controller.error)
	} else {
		return getResponseWithCode(controller.Code, controller.Data)
	}
}

func (controller*Controller) SetError(err IError) {
	GetLog().Info("set error in %s,track_id: %s, err: %s\n", controller.Ctx.Context.Request.URL.Path, controller.TrackId, err.GetMsg())
	controller.error = err
	return
}

func (controller*Controller) GetTrackId() string {
	return controller.TrackId
}

func (controller*Controller) SetTrackId(id string) {
	controller.TrackId = id
}

func (controller*Controller) SetContext(c *Context) {
	if controller.Ctx == nil {
		controller.Ctx = c
	}
}

func (controller*Controller) GetContext() *Context {
	return controller.Ctx
}
