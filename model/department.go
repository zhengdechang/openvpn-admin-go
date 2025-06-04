package model

import (
   "time"

   "github.com/google/uuid"
   "gorm.io/gorm"
)

// Department 部门模型
type Department struct {
   ID        string    `gorm:"primaryKey;size:36"`
   Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
   // HeadID 部门负责人用户ID
   HeadID    string    `gorm:"size:36" json:"headId,omitempty"`
   // Head 部门负责人信息
   Head      *User     `gorm:"foreignKey:HeadID" json:"head,omitempty"`
   // ParentID 上级部门ID
   ParentID  string       `gorm:"size:36" json:"parentId,omitempty"`
   // Parent 上级部门信息
   Parent    *Department  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
   // Children 子部门列表
   Children  []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
   CreatedAt time.Time
   UpdatedAt time.Time
   Users     []User    `gorm:"foreignKey:DepartmentID" json:"-"`
}

// BeforeCreate 在创建记录前生成 UUID
func (d *Department) BeforeCreate(tx *gorm.DB) (err error) {
   d.ID = uuid.NewString()
   return
}