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

//set model orm with gorm
func (m *Model) SetGorm(db *gorm.DB) {
	m.db = db
}

//alias get model name
func (m *Model) GetKeyName() string {
	return m.Info.Table
}

// get model primary key
func (m *Model) GetKey() int64 {
	if !m.checkValueSet() {
		return 0
	}

	v, e := m.Value.GetAttrInt64(m.GetKeyName())
	if e != nil {
		m.Err = e
		return 0
	}

	return v
}

//alias get model name
func (m *Model) GetName() string {
	return m.Info.Table
}

//get model table name
func (m *Model) TableName() string {
	return m.Info.Table
}

//get real value string
func (m *Model) GetString() string {
	if !m.checkValueSet() {
		return ""
	}

	b, _ := json.Marshal(m.Value)
	return string(b[:])
}

//get real value bytes
func (m *Model) GetBytes() []byte {
	if !m.checkValueSet() {
		return nil
	}

	b, _ := json.Marshal(m.Value)
	return b
}

//update model m.Value value.primary key
func (m *Model) Update() bool {
	if !m.checkValueSet() {
		return false
	}

	id, err := m.Value.GetAttrInt64(m.GetKeyName())
	if err != nil {
		m.Err = errors.New("must set Value use orm struct")
		return false
	}

	err = m.db.Model(m.Value).Table(m.TableName()).
		Where(m.GetKeyName() + " = " + utils.Int642String(id)).UpdateColumns(m.Value).Error
	if err != nil {
		m.Err = err
		return false
	}

	if m.db.RowsAffected != 1 {
		return false
	}

	return true
}

//create model with value
func (m *Model) Create() int64 {
	if !m.checkValueSet() {
		return 0
	}

	o := m.db.New()
	err := o.Table(m.TableName()).Create(m.Value).Error
	if err != nil {
		m.Err = err
		return 0
	}

	id, err := m.Value.GetAttrInt64(m.GetKeyName())
	if err != nil {
		m.Err = err
		return 0
	}

	if id < 1 {
		m.Err = errors.New("insert fail")
		return 0
	}

	return id
}

//delete model with value.primary key
func (m *Model) Delete() bool {
	if !m.checkValueSet() {
		return false
	}

	id, err := m.Value.GetAttrInt64(m.GetKeyName())
	if err != nil {
		m.Err = errors.New("must set Value use orm struct")
		return false
	}

	o := m.db.New()
	o.Table(m.TableName()).Where(m.GetKeyName() + " = " + utils.Int642String(id)).Delete(m.Value)
	if o.RowsAffected == 1 {
		return true
	}

	m.Err = errors.New("rows affected zero")

	return false
}

//get One Value impl IGetterSetter
func (m *Model) One(sql string, args ...interface{}) IGetterSetter {
	if !m.checkValueSet() {
		return nil
	}

	typ := reflect.TypeOf(m.Value)
	value := reflect.New(typ)

	o := m.db.New()
	err := o.Table(m.TableName()).Where(sql, args).Scan(value.Interface()).Error
	if err != nil {
		m.Err = err
		return nil
	}

	if v, ok := value.Interface().(IGetterSetter); ok {
		return v
	}

	m.Err = errors.New("value not impl IGetterSetter")
	return nil
}

//get List Value impl IGetterSetter
func (m *Model) List(sql string, args ...interface{}) []IGetterSetter {
	if !m.checkValueSet() {
		return nil
	}

	typ := reflect.TypeOf(m.Value)

	modelType := reflect.SliceOf(typ)
	sliceValue := reflect.MakeSlice(modelType, 0, 0)

	slice := reflect.New(sliceValue.Type())
	slice.Elem().Set(sliceValue)

	o := m.db.New()
	err := o.Table(m.TableName()).Where(sql, args).Scan(slice.Interface()).Error
	if err != nil {
		m.Err = err
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

//check value is nil
func (m *Model) checkValueSet() bool {
	if m.Value == nil {
		m.Err = errors.New("please set value as impl IGetterSetter")
		return false
	}

	return true
}
