package main

import (
	"PersonalBlog/internal/app"
)

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
