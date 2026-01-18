package main

import (
	"fmt"
	"go-rest-api/db"
	"go-rest-api/model"
)

/*
マイグレーション
データベースのスキーマ変更(テーブル定義の追加、更新、削除など)を管理し、
自動実行するプロセスなどシステム開発におけるデータ構造の変更を効率よくする
*/ 

// NewDB関数を呼び出して、DBの操作の内容を定義
func main() {
  dbConn := db.NewDB() // DB接続
	defer fmt.Println("移行に成功しました") // DB閉じる
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Task{}) // User/Taskテーブルを最新化(追加・作成)
}