package core

type CodeMapping map[string]string

func (cm CodeMapping) GetCodeInfo(code string) string {

	if v, ok := cm[code]; ok {
		return v
	}

	return ""
}

const (
	SuccessCode     = "000000"
	FailCode        = "500"
	FailUnknownCode = "-1"
	FailJsonParse   = "-2"
	FailThirdLogin  = "-3"
	FailSecret      = "-4"
	FailParamsError = "-5"
	FailPostMax     = "-6"
	FailInternal	= "400"

)


var DefaultCodeMapping = CodeMapping{
	FailCode:              "操作失败",
	FailUnknownCode:       "未知接口",
	FailJsonParse:         "JSON解析错误或者接口不存在",
	FailThirdLogin:        "第三方登录异常",
	FailSecret:            "加密参数错误",
	SuccessCode:           "操作成功",
	FailInternal:          "内部错误",
	FailParamsError:       "参数不合法",
	FailPostMax:           "数据超过限制的大小",
}
