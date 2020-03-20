package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grestful/logs"
	"github.com/grestful/session"
	"github.com/grestful/utils"
	"net/http"
	"strings"
	"time"
)

type Context struct {
	*gin.Context
	Session session.IUserSession
}

type ProcessFunc func(controller *Controller) IError
type Controller struct {
	//log trackId
	TrackId string //访问id
	//the gin context
	Ctx *Context
	//response code, default success
	Code string
	//controller error
	error IError
	//response data
	Data interface{} //output  from data raw data parse
	Req  interface{} //input params

	// middleware functions
	Middleware []ProcessFunc
	// logic functions
	ProcessFun ProcessFunc
}

//func Session(controller *Controller) error {
//	typ := gGore.GetSessionType()
//
//	sid := ""
//	cookName, err := gGore.Config.GetValue("", "SESSION_NAME")
//	if cookName == "" {
//		cookName = "sid"
//	}
//	cook, err := controller.GetContext().Request.Cookie(cookName)
//	if err != nil {
//		return err
//	}
//	if cook.Value == "" {
//		sid = controller.GetContext().Request.URL.Query().Get(cookName)
//	} else {
//		sid = cook.Value
//	}
//
//	if sid == "" {
//		return errors.New("sid is empty")
//	}
//
//	maxLeftStr, err := gGore.Config.GetValue("session", "max_left_time")
//	if err != nil || maxLeftStr == "" {
//		maxLeftStr = "3600"
//	}
//	maxLeft := utils.String2Int64(maxLeftStr, 3600)
//	switch typ {
//	case SessionTypeFile:
//		path, err := gGore.Config.GetValue("session", "file_path")
//		if err != nil || path == "" {
//			path = "/tmp"
//		}
//
//		sess := session.GetNewFileSession(path, maxLeft)
//		controller.Session = session.GetNewUserSession(sid, sess)
//		return nil
//	case SessionTypeRedis:
//		name, err := gGore.Config.GetValue("session", "redis_name")
//		if err != nil || name == "" {
//			name = "default"
//		}
//		c := gGore.GetRedis(name)
//		if c == nil {
//			return errors.New(fmt.Sprintf("can't find redis connection name %s", name))
//		}
//		maxLeft := utils.String2Int64(maxLeftStr, 3600)
//		sess := session.GetNewRedisSession(c, maxLeft)
//		controller.Session = session.GetNewUserSession(sid, sess)
//		return nil
//	}
//
//	return errors.New("session type error")
//}

// get new Controller
// demo:
//		c := GetNewController(g, &UserInfo{UserId:1})
//		c.ProcessFunc = func(con *Controller) IError {
//			u := con.Req.(*Req)
//			con.Data = u
//		}
func GetNewController(g *gin.Context, req interface{}) *Controller {
	return &Controller{
		Ctx: &Context{
			Context: g,
		},
		Middleware: make([]ProcessFunc, 0),
		Code:       SuccessCode,
		Req:        req,
	}
}

//run
func RunProcess(controller IController, g *gin.Context) {
	c := GetContext(g)
	controller.SetContext(c)
	var ierr IError
	defer func() {
		if x := recover(); x != nil {
			logs.Error(" panic :", x)
			ierr = Error{
				Code: "500",
				Msg:  "运行时内部错误",
				Err:  x,
			}
		}
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

func (controller *Controller) Use(fn func(controller *Controller) IError) {
	if controller.Middleware == nil {
		controller.Middleware = make([]ProcessFunc, 1)
	}
	controller.Middleware = append(controller.Middleware, fn)
}

// controller default Decode
func (controller *Controller) Decode() IError {
	controller.Data = nil

	switch controller.Ctx.Context.Request.Method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		ct := controller.Ctx.Context.Request.Header.Get("Content-Type")
		if strings.Contains(ct, "json") {
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

// controller default Process
func (controller *Controller) Process() IError {
	if controller.Middleware != nil {
		if len(controller.Middleware) > 0 {
			for _, m := range controller.Middleware {
				err := m(controller)
				if err != nil {
					return err
				}
			}
		}
	}
	if controller.ProcessFun != nil {
		return controller.ProcessFun(controller)
	}
	return nil
}

//ger ready response
func (controller *Controller) getResponse() Response {
	if controller.error != nil {
		return getDefaultErrorResponse(controller.error)
	} else {
		return getResponseWithCode(controller.Code, controller.Data)
	}
}

//set error
func (controller *Controller) SetError(err IError) {
	logs.Info("set error in %s, track_id: %s, err: %s\n", controller.Ctx.Context.Request.URL.Path, controller.GetTrackId(), err.GetMsg())
	controller.error = err
	return
}

//get trackId
func (controller *Controller) GetTrackId() string {
	return controller.TrackId
}

//set trackId
func (controller *Controller) SetTrackId(id string) {
	controller.TrackId = id
}

//set Context
func (controller *Controller) SetContext(c *Context) {
	if controller.Ctx == nil {
		controller.Ctx = c
	}
}

//get gin Content
func (controller *Controller) GetContext() *Context {
	return controller.Ctx
}

//query things with gin
func (controller *Controller) Query(key string) string {
	return controller.Ctx.Query(key)
}

//router params with gin
func (controller *Controller) Param(key string) string {
	return controller.Ctx.Param(key)
}
