package core

import "strconv"

type IUserSession interface {
	ISession
	SetData(data map[string]interface{}) error
	GetData() (data map[string]interface{}, err error)
	GetUserId() int64
	GetProperty(name string) interface{}
	GetAuthName() string
	SetSessionHandler(session ISession) error
}

type UserSession struct {
	UserId		int64	`json:"user_id"`
	Property    map[string]string `json:"property"`
	Sid         string  `json:"sid"`
}

func (s *UserSession) SetData(data map[string]string) error {
	if id,ok := data["user_id"]; ok {
		//s.UserId = base.String2Int64(id, 0)
		s.UserId,_ = strconv.ParseInt(id, 10, 64)
	}
	s.Property = data
	return nil
}

func (s *UserSession) GetData() (data map[string]string, err error) {
	return s.Property,nil
}

func (s *UserSession) GetUserId() int64 {
	return s.UserId
}

func (s *UserSession) GetProperty(name string) interface{} {
	if v,ok := s.Property[name]; ok {
		return v
	}
	return nil
}

func (s *UserSession) GetAuthName() string {
	return "cookie"
}
