package router

import (
	"go-rest-api/controller"
	"os"
  echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// エンドポイントを定義
func NewRouter(uc controller.IUserController,tc controller.ITaskController) *echo.Echo {
	e := echo.New()
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	t := e.Group("/tasks")
	// /tasks グループ配下の全エンドポイントにJWT認証ミドルウェアを適用し、CookieのtokenからJWTを取得してSECRETで検証することで、ログイン済みユーザーだけアクセス可能にする。
	// tasksのGroupにミドルウェアを適応できるようにする
	// Useを使うことで、エンドポイントにミドルウェアを追加することができる
	// JWT認証ミドルウェア
	/*
	リクエストからJWTを探す
  署名が正しいか検証する（改ざんチェック）
  OKなら「このリクエストはログイン済み」と判断して次へ
  NGなら 401 Unauthorized で止める

	✅ SigningKey: []byte(os.Getenv("SECRET"))
    JWTは「秘密鍵（SECRET）」で署名して発行しています。
    なので受け取ったJWTが
    本当にサーバーが発行したものか？
    改ざんされてないか？
    を確かめるために 同じSECRETで検証します。

		JWTを どこから探すか を指定しています。
    "cookie:token" は
    「Cookie の token という名前の値からJWTを取ってね」という意味
	*/

	t.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	// taskのエンドポイントを追加
	/*
	使い分け
GET /tasks：タスク一覧を取る（複数件）
POST /tasks：新しいタスクを作る（1件追加）

	*/
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)
  return e
}