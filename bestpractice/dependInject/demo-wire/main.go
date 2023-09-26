package main

import (
	"demo-wire/wire"
	"fmt"
)

func main() {
	// v1
	//ginEngine := wire.InitializeApp()
	//ginEngine := wire.InitializeAppV2()
	ginEngine := wire.InitializeAppV3()
	if err := ginEngine.Run(); err != nil {
		fmt.Println(err.Error())
	}

	// v2 , 效果和v1一样
	//wire.InitializeAppV2()
}
