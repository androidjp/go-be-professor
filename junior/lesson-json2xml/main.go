package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// 定义一个结构体来表示JSON数据
type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Body  string `json:"body"`
}

func main() {
	// 示例JSON数据
	jsonData := `{"name":"John Doe","email":"john@example.com"}`

	// 解析JSON数据到结构体
	var person Person
	err := json.Unmarshal([]byte(jsonData), &person)
	if err != nil {
		fmt.Println("解析JSON数据失败:", err)
		return
	}

	person.Body = jsonData

	// 将结构体转换为XML
	xmlData, err := xml.Marshal(person)
	if err != nil {
		fmt.Println("转换为XML失败:", err)
		return
	}

	// 输出转换后的XML内容
	fmt.Println(string(xmlData))
	// <Person><Name>John Doe</Name><Email>john@example.com</Email><Body>{&#34;name&#34;:&#34;John Doe&#34;,&#34;email&#34;:&#34;john@example.com&#34;}</Body></Person>

	// 特殊地，针对单引号和双引号，Golang使用的是十进制的字符串 进行替换，目的是让整个字符串更短，但需要考虑具体业务使用场景
	//escQuot = []byte("&#34;") // shorter than "&quot;"
	//	escApos = []byte("&#39;") // shorter than "&apos;"
	// 为避免对接方无法解析xml内容，如果要使用原来的 原始字符集，详情需要自己重写xml对应的方法，去将 这里的 &#34; 变回 &quot;

}
