package main

import (
	"fmt"
	"lessonprivatevar/model"
	"reflect"
	"unsafe"
)

func main() {

	// 反射
	p := model.NewPerson("小明", 18)

	ref := reflect.ValueOf(&p).Elem()

	fmt.Println(ref.Type())      // model.Person
	fmt.Println(ref.Kind())      // struct
	fmt.Println(ref.Interface()) // {小明 18}
	fmt.Println(ref.CanSet())    // false
	field := ref.FieldByName("name")
	fmt.Println(field.String())

	// unsafe 指针 操作私有字段
	realField := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()

	realField.SetString("小红")

	fmt.Println(p) // {小红 18}
	fmt.Println(realField.String())

}
