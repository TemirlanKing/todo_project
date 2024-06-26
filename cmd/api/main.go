package main

import (
	"log"

	"github.com/damirbeybitov/todo_project/internal/handlers"
	"github.com/damirbeybitov/todo_project/internal/models"
	"github.com/damirbeybitov/todo_project/internal/repository"
	"github.com/damirbeybitov/todo_project/internal/service"
	pbAuth "github.com/damirbeybitov/todo_project/proto/auth"
	pbTask "github.com/damirbeybitov/todo_project/proto/task"
	pbUser "github.com/damirbeybitov/todo_project/proto/user"
	"google.golang.org/grpc"
)

func main() {
	// Подключение к серверу микросервиса пользователей
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer userConn.Close()

	// Создание клиентского объекта
	userClient := pbUser.NewUserServiceClient(userConn)

	authConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer authConn.Close()

	authClient := pbAuth.NewAuthServiceClient(authConn)

	taskConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer taskConn.Close()

	taskClient := pbTask.NewTaskServiceClient(taskConn)

	microServiceClients := models.MicroServiceClients{
		UserClient: userClient,
		AuthClient: authClient,
		TaskClient: taskClient,
	}

	repo := repository.NewRepository(microServiceClients)

	handler := handlers.NewHandler(repo)

	service := service.NewService(handler)
	service.LaunchServer()

	// Вызов метода RegisterUser
	// registerResponse, err := microServiceClients.UserClient.RegisterUser(context.Background(), &pbUser.RegisterUserRequest{
	// 	Username: "damir",
	// 	Email:    "damir@example.com",
	// 	Password: "password",
	// })
	// if err != nil {
	// 	log.Fatalf("could not register user: %v", err)
	// }
	// log.Printf("Registered user with ID: %s", registerResponse.Id)

	// Вызов метода GetUserProfile
	// getUserProfileResponse, err := microserviceClients.UserClient.GetUserProfile(context.Background(), &pbUser.GetUserProfileRequest{
	// 	Id: "1", // Предполагается, что это ID только что зарегистрированного пользователя
	// })
	// if err != nil {
	// 	log.Fatalf("could not get user profile: %v", err)
	// }
	// log.Printf("Got user profile: %+v", getUserProfileResponse.User)

	// Вызов метода DeleteUser
	// deleteUserResponse, err := microserviceClients.UserClient.DeleteUser(context.Background(), &pbUser.DeleteUserRequest{
	// 	Username: "damir",
	// 	Password: "password",
	// })
	// if err != nil {
	// 	log.Fatalf("could not delete user: %v", err)
	// }
	// log.Printf(deleteUserResponse.Message)
}