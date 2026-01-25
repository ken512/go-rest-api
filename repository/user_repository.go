package repository

import (
	"go-rest-api/model"
	"gorm.io/gorm"
)

// repositoryのinterface
// ユーザーをDBから取得/作成する処理を、usecaseから切り離すためのinterface(窓口)
// repository：DB操作の抽象（保存/検索など）
type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error // ログイン時のユーザーのメールチェック・新規登録時のEメール重複チェック
	CreateUser(user *model.User) error                   // userをDBに新規保存(サインアップでユーザー生成)
}
// userRepository は「ユーザー関連のDB操作」をまとめる実装
// 目的：usecase（処理の流れ）からDB操作（GORM）を分離する
type userRepository struct {
	db *gorm.DB
}
// コンストラクタ：DB接続を受け取って repository を作る
// 返り値を interface にすることで、usecaseは実装詳細（GORM）を知らなくてよくなる
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}
// email でユーザーを1件取得する（ログイン/重複チェックで使用）
// user は取得結果を書き込む入れ物なのでポインタで受け取る
func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	// WHERE email = ? で検索し、最初の1件を user に詰める
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}
// CreateUser は「ユーザーをDBに保存する処理」
// 保存した結果（特に ID）が user に反映される可能性があるので、ポインタで受け取る
func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
