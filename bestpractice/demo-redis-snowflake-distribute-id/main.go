package main

import (
	"demo-redis-snowflake-distribute-id/idgen"
	"fmt"
)

func main() {

	fmt.Println("--------------------------------------")
	id := idgen.Get().GenSnowID()
	fmt.Println("snowID:", id)
	fmt.Println("--------------------------------------")
	id = idgen.Get().GenSnowID()
	fmt.Println("snowID:", id)
	fmt.Println("--------------------------------------")
	id = idgen.Get().GenSnowID()
	fmt.Println("snowID:", id)
	fmt.Println("--------------------------------------")
	idgen.Get().SnowIdToTime(id)
}
