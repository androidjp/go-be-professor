package main

import (
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("this is a viper demo")
	// 设置读取前缀为"APP"的系统环境变量
	viper.SetEnvPrefix("APP")
	// 绑定环境变量
	viper.AutomaticEnv()
	// 读取环境变量中APP_NAME的值
	viper.Get("APP_NAME")
	// 将环境变量APP_NAME的值绑定到app.name
	err := viper.BindEnv("app.name", "APP_NAME")
	if err != nil {
		fmt.Println(err)
	}
	// 默认值（在没有指定配置文件的情况下，使用默认值）
	viper.SetDefault("server.port", 8080)
	// 设置配置文件名称（不包含文件类型后缀）
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("toml")
	// 设置配置文件所在路径
	viper.AddConfigPath("./")
	// 启动一个监控项,当配置文件发生改动后自动加载新的配置项目，需要注意的是，在使用时需提前设置配置文件路径信息
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 定义配置
	viper.Set("data.newitem", "test111")
	viper.Set("data.newitem2", "test222")
	// 将上面定义的配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径（它先获取原始配置文件中的配置，通过上面定义的配置，对读取的配置进行修改/新增）,这里需要注意的时此方法会覆盖原始配置文件中的内容，并不是追加！！！
	err = viper.WriteConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 此方法和"viper.WriteConfig()"作用基本一致，区别为他不会覆盖原始配置文件中的内容，如果配置文件不存在，则会触发异常
	err = viper.SafeWriteConfig()
	if err != nil {
		fmt.Println(err)
	}
	// 将上面定义的配置写入指定路径的配置文件（它先获取原始配置文件中的配置，通过上面定义的配置，对读取的配置进行修改/新增）,这里需要注意的时此方法会覆盖原始配置文件中的内容，并不是追加！！！
	err = viper.WriteConfigAs("go_code/viper_demo/conf/config.toml")
	if err != nil {
		fmt.Println(err)
	}
	// 此方法和"viper.WriteConfigAs()"作用基本一致，区别为他不会覆盖原始配置文件中的内容,若配置文件不存在，它则会创建一个新的配置文件,若配置文件已存在，则会出发异常！！！
	err = viper.SafeWriteConfigAs("go_code/viper_demo/conf/config.toml")
	if err != nil {
		fmt.Println(err)
	}
	// 读取配置信息
	serverport := viper.Get("server.port")
	fmt.Println("port:", serverport)
	// 读取data标签中info配置中的name字段
	name := viper.Get("data.info.name")
	fmt.Println("name:", name)
	// 读取data标签中numbers配置
	numbers := viper.Get("data.numbers")
	fmt.Println("numbers:", numbers)
	// 读取data标签中prod配置,如果prod的值是true，则返回true，否则返回false
	prod := viper.GetBool("data.prod")
	fmt.Println("prod:", prod)
	// 读取data标签中info配置中的age字段,并将age的值转换为int类型
	age := viper.GetInt("data.info.age")
	fmt.Println("age:", age)
	// viper.GetDuration() 获取配置文件中的配置项，并转换为time.Duration类型
	// viper.GetFloat64() 获取配置文件中的配置项，并转换为float64类型
	// viper.AllKeys() 获取配置文件中的所有key
	// viper.Sub() 获取配置文件中的子配置项
	// viper.AllSettings() 获取配置文件中的所有配置项
	http.ListenAndServe(":8080", nil)
}
