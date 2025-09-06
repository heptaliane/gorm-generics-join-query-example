package main

type User struct {
	Id   uint
	Name string
}

type Product struct {
	Id    uint
	Name  string
	Price uint
}

type Order struct {
	Id        uint
	UserId    uint
	ProductId uint
}

type JoinedOrder struct {
	Order

	User    *User    `gorm:"foreginkey:UserId;references:Id"`
	Product *Product `gorm:"foreginkey:ProductId;references:Id"`
}
