package main

import (
	"fmt"
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDo(t *testing.T) {
	Do()
}

func TestDo2(t *testing.T) {
	Convey("aaaa", t, func() {
		Convey("bbbb", func() {
			So(Do(), ShouldEqual, "yes")
		})
	})
}

func TestDo3(t *testing.T) {
	Convey("函数打桩", t, func() {
		patches := ApplyFunc(Do, func() string {
			return "no"
		})
		//defer func() {
		//	fmt.Println("函数打桩defer")
		//	patches.Reset()
		//}()
		defer patches.Reset()

		So(Do(), ShouldEqual, "no")
		fmt.Println("函数打桩")
	})
}

func TestDo4(t *testing.T) {
	Convey("两个函数打桩", t, func() {
		patches := ApplyFunc(Do, func() string {
			return "no"
		})
		defer func() {
			fmt.Println("两个函数打桩defer")
			patches.Reset()
		}()

		So(Do(), ShouldEqual, "no")

		p2 := patches.ApplyFunc(Do2, func() string {
			return "gg"
		})
		defer func() {
			p2.Reset()
		}()

		//defer patches2.Reset()
		//time.Sleep(2 * time.Second)

		So(Do2(), ShouldEqual, "gg")
		fmt.Println("两个函数打桩")
	})
}

func Test05(t *testing.T) {
	Convey("函数指定第一次和第二次输出值", t, func() {
		output1 := "no1"
		output2 := "no2"
		output3 := "no3"
		patches := ApplyFuncSeq(Do, []OutputCell{
			{
				Values: Params{output1},
				Times:  0, // 1次
			},
			{
				Values: Params{output2},
				Times:  2, // 2次
			},
			{
				Values: Params{output3},
				Times:  1, // 1次
			},
		})
		defer patches.Reset()

		So(Do(), ShouldEqual, "no1")
		So(Do(), ShouldEqual, "no2")
		So(Do(), ShouldEqual, "no2")
		So(Do(), ShouldEqual, "no3")
	})
}
