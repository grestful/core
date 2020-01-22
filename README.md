# core
framework core instance

# controller

controller提供了基础类和实现方式

简单控制器实现：
```go
func Index(c *gin.Context) {
	p := &Controller{}

	p.ProcessFun = func(c *base.Controller) base.IError {
		c.Data = "hello world"
		return nil
	}

	base.RunProcess(p, c)
}
```

自定义控制器实现:
```go
type ExampleService struct {
	*core.Controller
}

func (ex *ExampleService) Process() core.IError {
	m :=  models.Order{}
	core.GetCore().GetDb("default").
		First(&m, "id=?",ex.Ctx.Param("id"))
	return nil
}
func Example(c *gin.Context) {
    base.RunProcess(&ExampleService{}, c)
}
```

# model

本框架内置使用gorm作为model，模型需要支持实现接口：
```go
type IModelStruct interface {
	TableName() string
	GetKeyName() string
}
```

示例：
```go
type Order struct {
	Id         uint64  `gorm:"column:id;type:bigint(20) unsigned" json:"id"`
	OrderId    uint64  `gorm:"column:order_id;type:bigint(20) unsigned" json:"order_id"`
	UserId     uint64  `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`
	Name       string  `gorm:"column:name;type:varchar(255)" json:"name"` // 订单内容的缩写或者主题描述
}

//get real primary key name
func (order *Order) GetKeyName() string {
	return "id"
}

//get real primary key name
func (order *Order) TableName() string {
	return "order"
}

m := GetModel(&Order{
    OrderId:1,
    UserId:2,
    Name:"ojbk",
})
id := m.Create()
fmt.Println(id)
ok := m.Update(map[string]interface{"Name":"ok","UserId":1})
if !ok {
    fmt.Println(m.Err)
}


list := m.List("UserId = ?", 1)
for _,v := range list {
    id, err := v.GetAttrUInt64("id", 0)
    fmt.Println(id, err)
}

ok = m.Delete()
if !ok {
    fmt.Println(m.Err)
}
```

# session

# cache

开发中
