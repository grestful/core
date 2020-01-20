package core

import (
	"encoding/json"
	"errors"
	"github.com/grestful/utils"
	"github.com/jinzhu/gorm"
	"reflect"
)

type Model struct {
	Value IGetterSetter
	Info  TableInfo
	Err   error

	db *gorm.DB
}

type TableInfo struct {
	Table string `json:"table"`
	Key   string `json:"key"`
}

func (m *Model) SetGorm(db *gorm.DB) {
	m.db = db
}

func (m *Model) GetKeyName() string {
	return m.Info.Table
}

func (m *Model) GetKey() int64 {
	if m.Value != nil {
		v, e := m.Value.GetAttrInt64(m.GetKeyName())
		if e == nil {
			return v
		}
	}

	return 0
}
func (m *Model) GetName() string {
	return m.Info.Table
}

func (m *Model) TableName() string {
	return m.Info.Table
}

func (m *Model) GetString() string {
	if m.Value != nil {
		b,_ := json.Marshal(m.Value)
		return string(b[:])
	}

	return ""
}

func (m *Model) GetBytes() []byte {
	if m.Value != nil {
		b,_ := json.Marshal(m.Value)
		return b
	}

	return nil
}

func (m *Model) Update() bool {
	if m.Value == nil {
		return false
	}

	err := m.db.Model(m.Value).Table(m.TableName()).UpdateColumns(m.Value).Error
	if err != nil {
		return false
	}

	if m.db.RowsAffected != 1 {
		return false
	}

	return true
}

func (m *Model) Create() int64 {
	o := m.db.New()
	err := o.Table(m.TableName()).Create(m.Value).Error
	if err != nil {
		return 0
	}

	id,err := m.Value.GetAttrInt64(m.GetKeyName())
	if err != nil {
		return 0
	}

	if id < 1 {
		return 0
	}

	return id
}

func (m *Model) Delete() bool {
	if m.Value != nil {
		id,err := m.Value.GetAttrInt64(m.GetKeyName())
		if err != nil {
			m.Err = errors.New("must set Value use orm struct")
			return false
		}

		o := m.db.New()
		o.Table(m.TableName()).Where(m.GetKeyName()+" = "+utils.Int642String(id)).Delete(m.Value)
		if o.RowsAffected == 1 {
			return true
		}
	}

	return false
}

func (m *Model) One(sql string, args ...interface{}) IGetterSetter {
	typ := reflect.TypeOf(m.Value)
	value := reflect.New(typ)

	o := m.db.New()
	err := o.Table(m.TableName()).Where(sql, args).Scan(value.Interface()).Error
	if err != nil {
		return nil
	}

	if v,ok := value.Interface().(IGetterSetter); ok {
		return v
	}

	return nil
}

func (m *Model) List(sql string, args ...interface{}) []IGetterSetter {
	typ := reflect.TypeOf(m.Value)

	modelType := reflect.SliceOf(typ)
	sliceValue := reflect.MakeSlice(modelType, 0, 0)

	slice := reflect.New(sliceValue.Type())
	slice.Elem().Set(sliceValue)

	o := m.db.New()
	err := o.Table(m.TableName()).Where(sql, args).Scan(slice.Interface()).Error
	if err != nil {
		return nil
	}

	l := slice.Elem().Len()
	if l > 0 {

		result := make([]IGetterSetter, 0)
		for i := 0; i < slice.Elem().Len(); i++ {
			v, ok := slice.Elem().Index(i).Interface().(IGetterSetter)
			if ok {
				result = append(result, v)
			}
		}

		return result
	}

	return nil
}
