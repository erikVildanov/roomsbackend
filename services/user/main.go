package main

import (
	"os"

	"github.com/joho/godotenv"

	"roomsbackend/services/user/handler"
	"roomsbackend/services/user/repository"

	userpb "roomsbackend/proto/generated/user"
	commonserver "roomsbackend/services/common/server"
	commonstorage "roomsbackend/services/common/storage"
)

func main() {
	// Загружаем переменные окружения из .env
	_ = godotenv.Load("../../.env.jwt")
	_ = godotenv.Load()

	// Инициализация БД
	cfg := commonstorage.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}
	db := commonstorage.NewPostgres(cfg)
	defer db.Close()

	// Инициализируем репозиторий профилей
	repo := repository.NewProfileRepository(db)

	// Создаём gRPC-обработчик UserService
	userHandler := &handler.UserHandler{
		Repo: repo,
	}

	go commonserver.StartREST(":8081", "localhost:50052", userpb.RegisterUserServiceHandlerFromEndpoint)
	commonserver.StartGRPC[userpb.UserServiceServer](":50052", userpb.RegisterUserServiceServer, userHandler)
}
