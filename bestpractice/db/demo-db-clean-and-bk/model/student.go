package model

import (
	"fmt"
)

/**

  CREATE TABLE student (
      id bigint AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键',
      name VARCHAR(32) NOT NULL DEFAULT '' COMMENT '姓名',
      age tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '年龄',
      ctime bigint NOT NULL DEFAULT 0 COMMENT '创建时间戳',
      mtime bigint NOT NULL DEFAULT 0 COMMENT '修改时间戳',
      created_by varchar(255) NOT NULL DEFAULT 'system' COMMENT '创建者',
      modified_by varchar(255) NOT NULL DEFAULT 'system' COMMENT '修改者',
      deleted_by varchar(255) NOT NULL DEFAULT 'system' COMMENT '删除者'
  ) COMMENT='学生表';

*/
// Student 结构体表示学生信息
type Student struct {
	// 学生的唯一标识符
	Id int `json:"id" db:"id"`
	// 学生的姓名
	Name string `json:"name" db:"name"`
	// 学生的年龄
	Age int `json:"age" db:"age"`
	// 创建时间
	Ctime int64 `json:"ctime" db:"ctime"`
	// 修改时间
	Mtime int64 `json:"mtime" db:"mtime"`
	// 创建者
	CreatedBy string `json:"created_by" db:"created_by"`
	// 修改者
	ModifiedBy string `json:"modified_by" db:"modified_by"`
	// 删除者（如果已删除）
	DeletedBy string `json:"deleted_by" db:"deleted_by"`
}

func (s *Student) Print() {
	// // 创建一个Student实例
	// student := Student{
	// 	Id:         1,
	// 	Name:       "John",
	// 	Age:        20,
	// 	Ctime:      time.Now(),
	// 	Mtime:      time.Now(),
	// 	CreatedBy:  "admin",
	// 	ModifiedBy: "admin",
	// 	DeletedBy:  "",
	// }
	fmt.Printf("学生ID: %d\n", s.Id)
	fmt.Printf("学生姓名: %s\n", s.Name)
	fmt.Printf("学生年龄: %d\n", s.Age)
	fmt.Printf("创建时间: %v\n", s.Ctime)
	fmt.Printf("修改时间: %v\n", s.Mtime)
	fmt.Printf("创建者: %s\n", s.CreatedBy)
	fmt.Printf("修改者: %s\n", s.ModifiedBy)
	fmt.Printf("删除者: %s\n", s.DeletedBy)
}
