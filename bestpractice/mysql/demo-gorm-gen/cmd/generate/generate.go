package main

import (
	"demo-gorm-gen/dal/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
		Mode:    gen.WithDefaultQuery,
	})

	g.ApplyBasic(model.Passport{}, model.User{})

	g.ApplyInterface(func(m model.Method) {}, model.User{})
	g.ApplyInterface(func(m model.UserMethod) {}, model.User{})

	g.Execute()
}
