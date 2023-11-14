package main

import (
	"demo-protobuf-struct/proto/demo"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func main() {
	// 结构体 转 proto 二进制
	attrStructObj, err := structpb.NewStruct(map[string]interface{}{
		"A": "B",
		"B": 1231.43,
	})

	fmt.Println(err)

	extraInfo := &demo.Location{Country: "tyy"}

	anyObj, _ := anypb.New(extraInfo)

	protoUser := &demo.User{
		Name:    "小明",
		Age:     18,
		Sex:     demo.SexType_Male,
		Balance: -19,
		LessonSourceMap: map[string]int64{
			"语文": 90,
			"数学": 90,
		},
		Friends: []*demo.User_FriendInfo{
			{
				Name: "Mike",
			},
		},
		Loc:          &demo.Location{Country: "China"},
		Sign:         []byte{},
		OtherContact: &demo.User_Qq{Qq: "1231231232"},
		Attrs:        attrStructObj,
		Extra:        anyObj,
	}

	bs, _ := proto.Marshal(protoUser)

	// proto 二进制 转回来
	resUser := &demo.User{}
	_ = proto.Unmarshal(bs, resUser)

	pmsg, err := anypb.UnmarshalNew(resUser.Extra, proto.UnmarshalOptions{})
	loc := pmsg.(*demo.Location)
	fmt.Println("loc.Country = ", loc.Country)

	jsonBS, _ := json.Marshal(resUser)
	fmt.Println(string(jsonBS))
}
