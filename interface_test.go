package core

import (
	"time"
)

type My struct {
	Name string
	Id   int
}

//func Test_ModelInterface(t *testing.T) {
//	my := &My{}
//	myType := reflect.TypeOf(my)
//	fmt.Println(myType)
//	//v := reflect.New(myType).Elem().Interface()
//	a:=reflect.MakeSlice(myType, 0, 0)
//	// I want to make array  with My
//	//a := make([](myType.(type),0)  //can compile
//	fmt.Println(a)
//}

type Order struct {
	Id         uint64  `gorm:"column:id;type:bigint(20) unsigned" json:"id"`
	OrderId    uint64  `gorm:"column:order_id;type:bigint(20) unsigned" json:"order_id"`
	UserId     uint64  `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`
	Name       string  `gorm:"column:name;type:varchar(255)" json:"name"` // 订单内容的缩写或者主题描述
	Cost       float64 `gorm:"column:cost;type:decimal(10,2) unsigned;default:'0.00'" json:"cost"`
	Price      float64 `gorm:"column:price;type:decimal(10,2) unsigned;default:'0.00'" json:"price"`
	PayType    uint8   `gorm:"column:pay_type;type:tinyint(3) unsigned" json:"pay_type"` // 订单类型：0-未知，1-微信，2-支付宝，3-ios，4-果冻
	Status     uint8   `gorm:"column:status;type:tinyint(1) unsigned" json:"status"`     // 订单状态：0-待支付，1-支付成功，2-支付失败
	CreateTime string  `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime string  `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP" json:"update_time"`
}

func (o Order) GetAttribute(key string) interface{} {
	return nil
}
func (o Order) SetAttribute(key string, value interface{}) bool {
	return false
}
func (o Order) SetAttributes(mp map[string]interface{}) bool {
	return false
}
func (o Order) GetAttributes() map[string]interface{} {
	return nil
}
func (o Order) GetAttrInt(key string) (int, error) {
	return 0, nil
}
func (o Order) GetAttrInt64(key string) (int64, error) {
	return 0, nil
}
func (o Order) GetAttrFloat(key string) (float32, error) {
	return 0, nil
}
func (o Order) GetAttrFloat64(key string) (float64, error) {
	return 0, nil
}

func (o Order) GetAttrUInt(key string) (uint, error) {
	return 0, nil
}
func (o Order) GetAttrUInt64(key string) (uint64, error) {
	return 0, nil
}

func (o Order) GetAttrBool(key string) (bool, error) {
	return false, nil
}
func (o Order) GetAttrString(key string) (string, error) {
	return "", nil
}
func (o Order) GetAttrTime(key string) (time.Time, error) {
	return time.Now(), nil
}

//func Test_OrmModel(t *testing.T) {
//
//
//	InitConfig("D:\\work\\pay2\\config\\config.ini")
//
//	m := &Model{
//		Value: nil,
//		Err:   Error{},
//		db:    nil,
//	}
//
//	m.db = GetDb("default")
//	result := m.List("status = ?", 1)
//	fmt.Println(result, "result is null")
//
//
//	//values := reflect.New(t)
//	//fmt.Println(values)
//	//values := reflect.ArrayOf(0, t)
//	//m.db.Table("order").Where("status = ?", "1").Scan(o)
//	//fff(m.db)
//}
//
//func fff(db *gorm.DB) {
//	o:=&Order{}
//	t := reflect.TypeOf(o)
//	values := reflect.SliceOf(t)
//	slice := reflect.MakeSlice(values, 0, 0)
//	s := reflect.New(slice.Type())
//	s.Elem().Set(slice)
//	fmt.Println(slice)
//	db.Table("order").Where("status = ?", "1").Scan(s.Interface())
//
//	fmt.Println("values", s.Elem().Len())
//	for i:=0; i<s.Elem().Len();i++ {
//		fmt.Println(s.Elem().Index(i).Interface().(*Order))
//	}
//}
