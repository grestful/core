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
	Value IModelStruct
	Err   error

	db *gorm.DB

	attrMap map[string]interface{}
}

func GetModel(modelStruct IModelStruct) *Model{
	return &Model{
		Value:   modelStruct,
		Err:     nil,
		db:      GetDb("default"),
		attrMap: make(map[string]interface{}),
	}
}

//set model orm with gorm
func (m *Model) SetGorm(db *gorm.DB) {
	m.db = db
}

//alias get model name
func (m *Model) GetKeyName() string {
	return m.Value.GetKeyName()
}

// get model primary key
func (m *Model) GetKey() int64 {
	if !m.checkValueSet() {
		return 0
	}

	v, e := m.GetAttrInt64(m.GetKeyName())
	if e != nil {
		m.Err = e
		return 0
	}

	return v
}

//alias get model name
func (m *Model) GetName() string {
	return m.Value.TableName()
}

//get model table name
func (m *Model) TableName() string {
	return m.Value.TableName()
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

	id, err := m.GetAttrInt64(m.GetKeyName())
	if err != nil {
		m.Err = errors.New("must set Value use orm struct")
		return false
	}

	if m.attrMap != nil {
		db := m.db.New()
		err = db.Model(m.Value).Table(m.TableName()).
			Where(m.GetKeyName() + " = " + utils.Int642String(id)).Update(m.attrMap).Error

		if err != nil {
			m.Err = err
			return false
		}

		if db.RowsAffected != 1 {
			return false
		}
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

	id, err := m.GetAttrInt64(m.GetKeyName())
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
	id, err := m.GetAttrInt64(m.GetKeyName())
	if err != nil {
		m.Err = errors.New("must set Value use orm struct")
		return false
	}

	o := m.db.New()
	o.Model(m.Value).Where(m.GetKeyName() + " = " + utils.Int642String(id)).Delete(m.Value)
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

	o := m.db.New()
	slice := utils.SliceBuildWithInterface(m.Value)
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

func (m *Model) GetAttribute(key string) interface{} {
	if !m.checkValueSet() {
		return nil
	}

	if v, ok := m.attrMap[key]; ok {
		return v
	}
	if !utils.StructExistsProperty(m.Value, key) {
		return nil
	}

	return reflect.ValueOf(m.Value).FieldByName(key).Interface()
}

func (m *Model) SetAttribute(key string, value interface{}) bool {
	if !m.checkValueSet() {
		return false
	}

	if utils.StructExistsProperty(m.Value, key) {
		if utils.StructSetFieldValue(m.Value, key, value) {
			m.attrMap[key] = value
			return true
		}
	}

	return false
}

func (m *Model) SetAttributes(mp map[string]interface{}) bool {
	if !m.checkValueSet() {
		return false
	}
	//for key,value := range mp {
	//	if utils.StructExistsProperty(m.Value, key) {
	//		vt := reflect.TypeOf(value).String()
	//		mt := reflect.TypeOf(m.Value).String()
	//
	//		if vt == mt {
	//			m.attrMap[key] = value
	//			reflect.ValueOf(m.Value).FieldByName(key).Set(reflect.ValueOf(value))
	//		}
	//	}
	//}
	m.attrMap = mp
	return true
}

func (m *Model) GetAttributes() map[string]interface{} {
	if !m.checkValueSet() {
		return nil
	}
	return m.attrMap
}

func (m *Model) GetAttrInt(key string) (int, error) {
	if !m.checkValueSet() {
		return 0, m.Err
	}
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
	if !m.checkValueSet() {
		return 0, m.Err
	}
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
	if !m.checkValueSet() {
		return 0, m.Err
	}
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
	if !m.checkValueSet() {
		return 0, m.Err
	}
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
	if !m.checkValueSet() {
		return 0, m.Err
	}
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
	if !m.checkValueSet() {
		return 0, m.Err
	}
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
	if !m.checkValueSet() {
		return false, m.Err
	}
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
	if !m.checkValueSet() {
		return "", m.Err
	}
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
	if !m.checkValueSet() {
		return time.Now(), m.Err
	}
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

//check value is nil
func (m *Model) checkValueSet() bool {
	if m.Value == nil {
		m.Err = errors.New("please set value as impl IGetterSetter")
		return false
	}

	return true
}
