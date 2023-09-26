//go:build wireinject

package entity

import "github.com/google/wire"

type Name string

func NewName() Name {
	return "AAA"
}

type PublicAccount string

func NewPublicAccount() PublicAccount {
	return "ggggg"
}

type User struct {
	MyName          Name
	MyPublicAccount PublicAccount
}

/*
上述代码中，首先定义了自定义类型 Name 和 PublicAccount 以及结构体类型 User，
并分别提供了 Name 和 PublicAccount 的初始化函数（providers）。
然后定义一个注入器（injectors）InitializeUser，用于构造连接提供者并构造 *User 实例。
*/
func InitializeUser() *User {
	wire.Build(
		NewName,
		NewPublicAccount,
		wire.Struct(new(User), "MyName", "MyPublicAccount"),
	)
	return &User{}
}

// 绑定值
// 我们可以在注入器中通过 值表达式 给一个类型进行赋值，而不是依赖提供者（providers）
// 要注意的是，值表达式将被复制到生成的代码文件中。
func InjectUser() User {
	wire.Build(wire.Value(User{
		MyName:          "yongwww",
		MyPublicAccount: "ggg",
	}))
	return User{}
}
