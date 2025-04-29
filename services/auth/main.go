package main

import (
	"os"

	"roomsbackend/libs/jwt"

	"github.com/joho/godotenv"

	"roomsbackend/services/auth/handler"
	"roomsbackend/services/auth/repository"

	authpb "roomsbackend/proto/generated/auth"
	userpb "roomsbackend/proto/generated/user"
	grpcclient "roomsbackend/services/common/grpc"
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

	jwt.Init()

	// Создание репозитория
	repo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)

	conn := grpcclient.New("localhost:50052")
	defer conn.Close()
	userClient := userpb.NewUserServiceClient(conn)

	// gRPC хендлер с внедрением репозитория
	authHandler := &handler.AuthHandler{
		Repo:        repo,
		RefreshRepo: refreshRepo,
		UserGRPC:    userClient,
	}

	go commonserver.StartREST(":8080", "localhost:50051", authpb.RegisterAuthServiceHandlerFromEndpoint)
	commonserver.StartGRPC[authpb.AuthServiceServer](":50051", authpb.RegisterAuthServiceServer, authHandler)
}
