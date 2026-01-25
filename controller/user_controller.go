package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"
  "fmt"
	"github.com/labstack/echo/v4"
)

// User系のHTTPハンドラを揃えるための型
type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

// controllerに対してusecaseをDIするため
func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

// userController型がIUserControllerを満たすために、IUserController内に定義された３つのメソッド(LogOut,Login,SignUp)を実装する必要がある

// エラー時にHHTPレスポンスを返す処理
func (uc *userController) SignUp(c echo.Context) error {
	// リクエストボディ(JSON等)を受け取って、User構造体に詰める
	user := model.User{}
	// Bind(読み取り/変換)失敗を400で返す
	if err := c.Bind(&user); err != nil {
		// JSONで、リクエストのステータスとエラー内容を返す
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Bindに成功した場合
	userRes, err := uc.uu.SignUp(user)
	// サインアップに失敗した場合
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// サインアップに成功した場合
	// リクエストのステータスと新しく作成したuser情報を返す
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// user_usecaseのログインメソッドを呼び出し、JWTトークンを生成する
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// 成功した場合、取得したJWTトークンをサーバーサイドでクッキーに設定
	// ログイン成功時にJWTをCookieとしてブラウザに保存させる処理
	// 以後、ブラウザは同じ条件のリクエストでそのCookieを自動送信

	/*
	   このCookie設定の役割は：
	   JWTを“ブラウザに覚えさせる”
	   次回以降のAPIリクエストで 自動的に送らせる
	   JSから読ませず（HttpOnly）盗まれにくくする
	*/

	// cookie構造体を新しく生成
	cookie := new(http.Cookie)
	// cookieの名前をtoken
	cookie.Name = "token"
	// cookieの値を、JWTトークンを代入
	cookie.Value = tokenString
	// cookieの有効機銀を24時間に指定
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	// domainは、環境変数で定義しているAPI_DOMAINを割り当てる
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	// クライアント側にtokenの値が読み取れないようにする
	cookie.HttpOnly = true
	// フロントとAPIが別オリジンでもCookie送信を許可する
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	fmt.Println("JWT:", tokenString)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	// cookieの値をクリアしたいので空文字指定
	cookie.Value = ""
	// cookieの有効機銀をすぐに実行
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}
