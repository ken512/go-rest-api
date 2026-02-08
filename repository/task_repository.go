package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// repositoryがGORM/SQLの詳細を担当

// Taskのinterface
// Taskの永続化(取得/作成/更新/削除)の窓口。
// usecaseがGORM/SQLに依存しないようにinterfaceで抽象化する。
type ITaskRepository interface {
	GetAllsTasks(tasks *[]model.Task, userId uint) error // ユーザーのtaskの一覧を取得
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

// taskRepositoryの構造体
// repository実装（GORMで実際にDBを触る）
type taskRepository struct {
	db *gorm.DB
}

// DB接続を受け取る
// 返り値をinterfaceにすることで、usecaseは実装詳細を知らなくてよくなる
func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

// taskの作成日時が一番新しいものが末尾に来る順番でデータを取得する
func (tr *taskRepository) GetAllsTasks(tasks *[]model.Task, userId uint) error {
	if err := tr.db.Joins("User").Where("user_id=?", userId).Order("created_at").Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

// ログイン中ユーザーのタスクだけをtaskIdで1件取得し、見つかったtaskに詰める
func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return  nil
}

func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

// 指定した1件のtaskのtitleだけを更新する
func (tr *taskRepository)	UpdateTask(task *model.Task, userId uint, taskId uint) error {
	// 指定した taskId のタスクを、ログイン中ユーザー(userId)のものに限定して検索し、
	// 見つかった1件の title だけを task.Title に更新する（他人のタスクは更新できない）
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
	// DB更新でエラーが起きたらそのまま返す（接続/SQLなど）
	if result.Error != nil {
		return result.Error
	}
	// 更新対象が0件なら、存在しない or 自分のタスクではない（条件に一致しない）扱い
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
		// 指定した taskId のタスクを、ログイン中ユーザー(userId)のものに限定して削除する
	// → 他人のタスクは削除できないようにするための条件
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}