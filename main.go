package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
)

// アプリ起動するために、格レイヤー層を順番に組み立てて接続する
func main() {
	db := db.NewDB()
	userRepository := repository.NewUserRepository(db)
	// DIを追加
	taskRepository := repository.NewTaskRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	taskUsecase := usecase.NewTaskRepository(taskRepository)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	e := router.NewRouter(userController, taskController)
	// サーバーを起動
	e.Logger.Fatal(e.Start(":8080"))
}

