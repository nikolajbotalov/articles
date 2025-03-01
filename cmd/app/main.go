package main

import (
	// _ "PersonalBlog/docs"
	"PersonalBlog/internal/app"
)

// @title Personal Blog API
// @version 1.0
// @description API для управления статьями в личном блоге
// @host localhost:8080
// @BasePath /
func main() {
	// инициализация приложения
	application, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer application.Close()

	// запуск сервера
	application.Server.Run()
}
