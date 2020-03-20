# middleware

the middleware must impl  ProcessFunc (func(controller *Controller) IError)

demo:
```
func SessionMiddleware(c *core.Controller) core.IError {
	sid := c.GetContext().Query("uid")
	if sid == "" {
		return core.NewErrorCode(common.CodeLoginSidNotExists)
	}

	c.Ctx.Session = GetNewUserSessionStore(sid, Session)
	return nil
}
```


controller demo:
```
func LoginLog(c *gin.Context) {
	p := core.GetNewController(c, &user.UserLoginLog{})
    p.Use(models.SessionMiddleware)
    	p.ProcessFun = func(controller *core.Controller) core.IError {
    		var err error
    		controller.Data,err = controller.Ctx.Session.GetData()
    		if err != nil {
    			return core.NewError(err)
    		}
    		return nil
	}
	core.RunProcess(p, c)
}
```