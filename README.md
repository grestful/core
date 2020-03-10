# core
framework core instance

# controller

controller提供了基础类和实现方式

简单控制器实现：
```
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
```
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

# model(可用配套mode生成器生成)

示例：
```
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
```

# session （1.1.4）

quote from github.com/grest/session


# cache

开发中
