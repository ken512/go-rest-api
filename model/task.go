package model

import "time"
// どのユーザーがタスクを作成したのかをわかるように
type Task struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Title     string      `json:"title" gorm:"not null"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	User      User        `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`// constraint:OnDelete:CASCADEで、ユーザーにを削除したときその削除したユーザーに紐づいているタスクも一緒に削除される
	UserId    uint        `json:"user_id gorm:"not null""`
}

// クライアントからGETメソッドにリクエストがあったとき、クライアント側に返すタスクのデータ構造の定義
type TaskResponse struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Title     string      `json:"title" gorm:"not null"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`	
}