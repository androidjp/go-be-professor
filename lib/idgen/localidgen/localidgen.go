package localidgen

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gofrs/uuid"
	"github.com/lithammer/shortuuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

func LocalIDGenerator() {
	fmt.Println("--------------------------------------")
	fmt.Println(`---UUID的 v1版本：采用时间 和 mac地址的算法。`)
	uuidV1()

	fmt.Println("--------------------------------------")
	fmt.Println(`---UUID的 v4版本：纯随机数`)
	uuidV4()

	fmt.Println("--------------------------------------")
	fmt.Println(`---shortUUID，基于uuidV4，定长22字节`)
	shortUUID()

	fmt.Println("--------------------------------------")
	fmt.Println(`---xid：时间戳4位+MAC地址3位+PID2位+有序随机数3位，共12位，长度12字节`)
	xidFunc()

	fmt.Println("--------------------------------------")
	fmt.Println(`---ksuid（时间戳4字节+随机数16字节），长度20`)
	ksuidFunc()

	fmt.Println("--------------------------------------")
	fmt.Println(`---ulid，随机数+时间戳，长度：26`)
	ulidFunc()

	fmt.Println("--------------------------------------")
	fmt.Println(`---snowflake雪花算法(开源)，长度：19`)
	fmt.Println(`+--------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
+--------------------------------------------------------------------------+`)
	fmt.Println("默认每毫秒生成4096个ID")
	snowflakeFunc()

	fmt.Println("--------------------------------------")
	fmt.Println("---nanoID (比uuid更短：21位，且可用于URL中)")
	nanoID()
}

func nanoID() {
	defaultID, _ := gonanoid.New()

	fmt.Println("default:", defaultID, ", length:", len(defaultID))

	customID, _ := gonanoid.Generate("abcde", 54)
	fmt.Println("custom:", customID, ", length:", len(customID))

	customID, _ = gonanoid.Generate("abcde", 10)
	fmt.Println("custom2:", customID, ", length:", len(customID))
}

func snowflakeFunc() {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	fmt.Println("id.String(): ", id.String(), ", length: ", len(id.String()))

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d  (毫秒)\n", id.Time())
	fmt.Printf("ID Time format : %s\n", time.UnixMilli(id.Time()))

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())

}

func ulidFunc() {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	// id: 01G902ZSM96WV5D5DC5WFHF8WY length: 26
	fmt.Println("id:", id.String(), "length:", len(id.String()))
}

func ksuidFunc() {

	id := ksuid.New()
	// id: 2CWvPg766SUvezbiiV9nzrTZsgf length: 20
	fmt.Println("id:", id, "length:", len(id))

	id1 := ksuid.New()
	id2 := ksuid.New()
	// id1:2CTqTLRxCh48y7oUQzQHrgONT2k id2:2CTqTHf07C09CXyRMHdGKXnY5HP
	fmt.Println(id1, id2)

	// 支持ID对比，这个功能比较鸡肋了,目前没想到有用的地方
	compareResult := ksuid.Compare(id1, id2)
	fmt.Println(compareResult) // -1

	// 判断顺序性 (存在误判的情况！！！！)
	isSorted := ksuid.IsSorted([]ksuid.KSUID{id1, id2})
	fmt.Println(isSorted) // true 降序
}

func xidFunc() {
	// hostname+pid+atomic.AddUint32
	id := xid.New()
	containerName := "我们的唯一容器名"
	// 由于xid默认使用可重复ip地址填充4 5 6位 MAC地址位。
	// 实际场景中，服务都是部署在docker中，这里把ip地址位替换成了容器名
	// 这里只取了容器名MD5的前3位，验证会重复，放弃使用
	containerNameID := make([]byte, 3)
	hw := md5.New()
	hw.Write([]byte(containerName))
	copy(containerNameID, hw.Sum(nil))
	id[4] = containerNameID[0]
	id[5] = containerNameID[1]
	id[6] = containerNameID[2]

	// id: cbgjhf89htlrr1955d5g length: 12
	fmt.Println("xid第一次生成-id:", id, "length:", len(id))
	id = xid.New()
	id[4] = containerNameID[0]
	id[5] = containerNameID[1]
	id[6] = containerNameID[2]
	fmt.Println("xid第一次生成-id:", id, "length:", len(id))
	id = xid.New()
	id[4] = containerNameID[0]
	id[5] = containerNameID[1]
	id[6] = containerNameID[2]
	fmt.Println("xid第一次生成-id:", id, "length:", len(id))

	id = xid.NewWithTime(time.Date(2024, 3, 6, 0, 0, 0, 0, time.UTC))
	id[4] = containerNameID[0]
	id[5] = containerNameID[1]
	id[6] = containerNameID[2]
	fmt.Println("xid第一次生成-id:", id, "length:", len(id))

	/**
	特点：
	- 长度短。
	- 有序。
	- 不重复。
	- 时间戳这个随机数原子+1操作，避免了时钟回拨的问题。
	*/
	id = xid.NewWithTime(time.Date(2024, 3, 6, 0, 0, 0, 0, time.UTC))
	id[4] = containerNameID[0]
	id[5] = containerNameID[1]
	id[6] = containerNameID[2]
	fmt.Println("xid第一次生成-id:", id, "length:", len(id))

}

func shortUUID() {
	id := shortuuid.New()
	// id: iDeUtXY5JymyMSGXqsqLYX length: 22
	fmt.Println("id:", id, "length:", len(id))

	// V22s2vag9bQEZCWcyv5SzL 固定不变
	id = shortuuid.NewWithNamespace("http://127.0.0.1.com")
	// id: K7pnGHAp7WLKUSducPeCXq length: 22
	fmt.Println("id:", id, "length:", len(id))
	time.Sleep(50 * time.Millisecond)
	id = shortuuid.NewWithNamespace("http://127.0.0.1.com")
	fmt.Println("id:", id, "length:", len(id))
	time.Sleep(50 * time.Millisecond)
	id = shortuuid.NewWithNamespace("http://127.0.0.1.com")
	fmt.Println("id:", id, "length:", len(id))

	// NewWithAlphabet函数可以用于自定义的基础字符串，字符串要求不重复、定长57
	str := "12345#$%^&*67890qwerty/;'~!@uiopasdfghjklzxcvbnm,.()_+·><"
	id = shortuuid.NewWithAlphabet(str)
	// id: q7!o_+y('@;_&dyhk_in9/ length: 22
	fmt.Println("id:", id, "length:", len(id))
	time.Sleep(50 * time.Millisecond)
	id = shortuuid.NewWithAlphabet(str)
	fmt.Println(time.Millisecond)
	id = shortuuid.NewWithAlphabet(str)
	fmt.Println("id:", id, "length:", len(id))
}

func uuidV4() {
	var id uuid.UUID
	var err error
	// Version 4:是纯随机数,error会在内部报panic
	id, err = uuid.NewV4()
	if err != nil {
		fmt.Printf("uuid NewUUID err:%+v", err)
	}
	// id: 3b4d1268-9150-447c-a0b7-bbf8c271f6a7 length: 36
	fmt.Println("id:", id.String(), "length:", len(id.String()))
}

func uuidV1() {
	// Version 1:时间+Mac地址
	id, err := uuid.NewV1()
	if err != nil {
		fmt.Printf("uuid NewUUID err:%+v", err)
	}
	// id: f0629b9a-0cee-11ed-8d44-784f435f60a4 length: 36
	fmt.Println("id:", id.String(), "length:", len(id.String()))
}
