package model

import "time"

type User struct {
	ID         uint       `json:"id" gorm: "primaryKey"` // primaryKeyで主キーの役割
	Email      string     `json:"email" gorm: "unique"`  // uniqueで重複しないようにする
	Password   string     `json:"password"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
} // 構造体のオブジェクトをjson形式に変換したときに、フィールドの名前を自動的に小文字へ変換されるようになる


// クライアントにレスポンスするデータ型
type UserResponse struct {
	ID         uint       `json:"id" gorm: "primaryKey"`
	Email      string     `json:"email" gorm: "unique"`
}