package core

import (
	"encoding/json"
	"errors"
	"github.com/grestful/utils"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

type Model struct {
	Value interface{}
	Info  TableInfo
	Err   error

	db *gorm.DB

	attrMap map[string]interface{}
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
		v, e := m.GetAttrInt64(m.GetKeyName())
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
		b, _ := json.Marshal(m.Value)
		return string(b[:])
	}

	return ""
}

func (m *Model) GetBytes() []byte {
	if m.Value != nil {
		b, _ := json.Marshal(m.Value)
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

	id, err := m.GetAttrInt64(m.GetKeyName())
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
		id, err := m.GetAttrInt64(m.GetKeyName())
		if err != nil {
			m.Err = errors.New("must set Value use orm struct")
			return false
		}

		o := m.db.New()
		o.Table(m.TableName()).Where(m.GetKeyName() + " = " + utils.Int642String(id)).Delete(m.Value)
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

	if v, ok := value.Interface().(IGetterSetter); ok {
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

func (m *Model) GetAttribute(key string) interface{} {
	if v, ok := m.attrMap[key]; ok {
		return v
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return nil
	}

	return reflect.ValueOf(m.Value).FieldByName(key).Interface()
}

func (m *Model) SetAttribute(key string, value interface{}) bool {
	if utils.StructExistsProperty(m.Value, key) {
		vt := reflect.TypeOf(value).String()
		mt := reflect.TypeOf(m.Value).String()

		if vt == mt {
			m.attrMap[key] = value
			reflect.ValueOf(m.Value).FieldByName(key).Set(reflect.ValueOf(value))
			return true
		}
	}

	return false
}
func (m *Model) SetAttributes(mp map[string]interface{}) bool {
	m.attrMap = mp
	return true
}
func (m *Model) GetAttributes() map[string]interface{} {
	return m.attrMap
}
func (m *Model) GetAttrInt(key string) (int, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(int), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return 0, errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(int)
	if ok {
		return v, nil
	}
	return 0, errors.New("type is not int")
}
func (m *Model) GetAttrInt64(key string) (int64, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(int64), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return 0, errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(int64)
	if ok {
		return v, nil
	}
	return 0, errors.New("type is not int64")
}
func (m *Model) GetAttrFloat(key string) (float32, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(float32), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return 0, errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(float32)
	if ok {
		return v, nil
	}
	return 0, errors.New("type is not float32")
}
func (m *Model) GetAttrFloat64(key string) (float64, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(float64), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return 0, errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(float64)
	if ok {
		return v, nil
	}
	return 0, errors.New("type is not float64")
}

func (m *Model) GetAttrUInt(key string) (uint, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(uint), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return 0, errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(uint)
	if ok {
		return v, nil
	}
	return 0, errors.New("type is not uint")
}

func (m *Model) GetAttrUInt64(key string) (uint64, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(uint64), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return 0, errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(uint64)
	if ok {
		return v, nil
	}
	return 0, errors.New("type is not uint")
}

func (m *Model) GetAttrBool(key string) (bool, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(bool), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return false, errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(bool)
	if ok {
		return v, nil
	}
	return false, errors.New("type is not bool")
}

func (m *Model) GetAttrString(key string) (string, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(string), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return "", errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(string)
	if ok {
		return v, nil
	}
	return "", errors.New("type is not string")
}

func (m *Model) GetAttrTime(key string) (time.Time, error) {
	if v, ok := m.attrMap[key]; ok {
		return v.(time.Time), nil
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return time.Now(), errors.New("value is not struct field")
	}

	v, ok := reflect.ValueOf(m.Value).FieldByName(key).Interface().(time.Time)
	if ok {
		return v, nil
	}
	return time.Now(), errors.New("type is not time.Time")
}
