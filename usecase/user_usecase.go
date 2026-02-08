package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// usecase の interface
// ユーザーのサインアップ/ログインの処理をcontrollerがusecaseをどう呼べるかを定義するinterface
// usecase：アプリの目的（サインアップ/ログイン）を実現する
type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error) // 目的：ユーザー登録
	Login(user model.User) (string, error)              // 目的：ログイン
}

// IUserUsecaseを実際に動かす実装本体。
// repository(IUserRepository)を持ち、DIで具体実装を注入してDB依存を分離する。
type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

// NewUserUsecase は usecase の生成関数（DI用）。
// repository の実装(ur)を外から受け取り、userUsecase に注入して返す。
// usecase がDB実装に直接依存しないようにし、テストではモックに差し替えやすくする。
func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

// サインアップのusecaseの処理
func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10) // パスワードをハッシュ化
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)} // ハッシュ化したパスワードでユーザーをDBに保存
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	// DBに保存した結果を、APIで返せる形で返却する
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

/*
メールアドレスでDBからユーザーを取る

入力パスワードが正しいか（bcryptで）確認する

正しければJWTトークンを作る（有効期限つき）

署名して文字列にして返す（クライアントが以後それを持つ）
*/
func (uu *userUsecase) Login(user model.User) (string, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}
	// 登録済みユーザーを取得
	storeeUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storeeUser, user.Email); err != nil {
		return "", err
	}
	// storeeUser.Password は DBに保存されているハッシュ（bcryptの結果）user.Password は ログインフォームで入力された生パスワード
	err := bcrypt.CompareHashAndPassword([]byte(storeeUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	// JWTトークンを作る
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // claimsはトークンに入れる情報(ペイロード) HS256 は署名方式（秘密鍵で署名するタイプ）
		"user_id": storeeUser.ID,                         // 誰とログインしたか
		"exp":     time.Now().Add(time.Hour * 12).Unix(), // 有効期限(12時間後)
	})
	// 署名を改竄できないよう文字列にするして、改ざんされたら検知する
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
