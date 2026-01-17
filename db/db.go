package db

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}// GORMでDB操作するためのオブジェクト（gorm.DB）を受け取って返してる

	// DBに接続するためのURLを作成
	// Sprintfを使うことで、指定したフォーマットに従って文字列が整形・変換をしてくれる。
	// urlを使って、DBに接続
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_PW"), os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{}) // gorm.OpenでDBを開く(postgresをオープン(url), からの構造体(デフォルトの値でDBを起動))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("接続済み")
	return db
} 

// DBを閉じる
func CloseDB(db *gorm.DB) {
	sqDB, _ := db.DB() // *sql.DBを取り出す
	if err := sqDB.Close(); err != nil {
		log.Fatalln(err)
	}
}






/*
NewDB()の役割
・環境変数からDB接続情報を読み取る
・GORMでDB接続
・接続できた*gorm.DBを返す(この戻り値を他の場所で使ってDB操作をする)

os.Getenv("GO_ENV")
・GO_ENVという環境変数を見ている

if ... == "dev"
・開発中は.envにDBパスワードを入れておく
・本番環境は.envを使わず、デプロイ先の環境変数を入れる

godotenv.Load()
・.envファイルを読み込んでいる
・中に書かれたPOSTGRES_USER=...を、Goのos.Getenv()で取れる環境変数として登録をしてくれる

log.Fatalln(err)
・エラーを表示して、そこでプログラムを強制終了させる
*/
