package main

import (
	"encoding/json"
	"fmt"
)

type ResultObj struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ReqObj struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
	addr string `json:"addr"`
}

func main() {
	mm := map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": ReqObj{
			Name: "xiaom",
			Age:  18,
			addr: "街道xxxx",
		},
	}

	bbb, err := json.Marshal(mm)
	if err != nil {
		return
	}

	var res ResultObj

	fmt.Println(string(bbb)) // {"code":200,"data":{"name":"xiaom","age":18},"msg":"success"}

	if err := json.Unmarshal(bbb, &res); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res.Code) // 200
	fmt.Println(res.Msg)  // success
	fmt.Println(res.Data) // map[age:18 name:xiaom]
}
