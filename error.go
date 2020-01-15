package core

type Error struct {
	Code string
	Msg  string
	Err  interface{}
}

func (err Error) Error() string {
	if err.Msg == "" {
		return DefaultCodeMapping.GetCodeInfo(err.GetCode())
	}

	return err.Msg
}

func (err Error) GetMsg() string {
	return err.Error()
}

func (err Error) GetCode() string {
	if err.Code == "" {
		return "500"
	}

	return err.Code
}

func (err Error) IsNil() bool {
	return err.Msg == "" && err.Code == ""
}

func (err Error) GetDetail() interface{} {
	if err.IsNil() {
		return err.Err
	}

	return nil
}

func NewError(err error) Error {
	return Error{Err: err}
}
func NewErrorCode(code string) Error {
	return Error{Code: code}
}

func NewErrorStr(msg string) Error {
	return Error{Code: "500", Msg: msg}
}
