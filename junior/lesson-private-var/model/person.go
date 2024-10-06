package model

// 结构体定义，字段 age 是私有的
type Person struct {
	name string
	age  int
}

// 包内函数，能够访问私有字段
func NewPerson(name string, age int) Person {
	return Person{name: name, age: age}
}

func GetPersonAge(p Person) int {
	return p.age
}
