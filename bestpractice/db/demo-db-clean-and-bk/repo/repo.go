package repo

import (
	"context"
	"demodbclient/model"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo() *Repo {
	// 数据库连接字符串
	dsn := "root:@tcp(127.0.0.1:4000)/test?charset=utf8mb4&parseTime=True&loc=Local"
	// 使用sqlx连接数据库
	db, err := sqlx.Connect("mysql", dsn)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		panic(err)
	}
	// 测试数据库连接是否正常
	err = db.PingContext(context.Background())
	if err != nil {
		fmt.Printf("数据库连接不可用: %v\n", err)
		panic(err)
	}
	fmt.Println("成功连接到TiDB数据库")
	return &Repo{
		db: db,
	}
}

func (r *Repo) Stop() {
	// 在这里可以进行后续的数据库操作，如查询、插入等
	// 操作完成后关闭数据库连接
	if r.db != nil {
		r.db.Close()
	}
}

func (r *Repo) InsertStudent(stu *model.Student) error {
	// 插入数据
	_, err := r.db.NamedExecContext(context.Background(), "INSERT INTO student (name, age, ctime, mtime) VALUES (:name, :age, :ctime, :mtime)", stu)
	if err != nil {
		fmt.Printf("插入数据失败: %v\n", err)
		return err
	}
	return nil
}

func (r *Repo) BatchRead() ([]*model.Student, error) {
	var stus []*model.Student
	err := r.db.SelectContext(context.Background(), &stus, "select * from student")
	if err != nil {
		fmt.Printf("查询数据失败: %v\n", err)
		return nil, err
	}
	return stus, nil
}

func (r *Repo) DeleteOutdated() error {
	// 找到所有过期数据
	res, err := r.db.ExecContext(context.Background(), "DELETE FROM student WHERE mtime <?", time.Now().Add(-time.Hour).Unix())

	// 删除过期数据

}
