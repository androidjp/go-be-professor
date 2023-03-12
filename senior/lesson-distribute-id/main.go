package main

import (
	"fmt"
	"github.com/gofrs/uuid"
)

func main() {
	// Version 1:时间+Mac地址
	id, err := uuid.NewV1()
	if err != nil {
		fmt.Printf("uuid NewUUID err:%+v", err)
	}
	// id: f0629b9a-0cee-11ed-8d44-784f435f60a4 length: 36
	fmt.Println("id:", id.String(), "length:", len(id.String()))

	// Version 4:是纯随机数,error会在内部报panic
	id, err = uuid.NewV4()
	if err != nil {
		fmt.Printf("uuid NewUUID err:%+v", err)
	}
	// id: 3b4d1268-9150-447c-a0b7-bbf8c271f6a7 length: 36
	fmt.Println("id:", id.String(), "length:", len(id.String()))
}
