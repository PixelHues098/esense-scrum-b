package main

import (
	"esense/controller"
	"esense/database"
	"esense/middleware"
	"esense/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	serveApplication()
}

func loadDatabase() {
	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	database.Database.AutoMigrate(&model.Project{})
	database.Database.AutoMigrate(&model.Issue{})
	database.Database.AutoMigrate(&model.Sprint{})
	database.Database.AutoMigrate(&model.Swimlane{})
	database.Database.AutoMigrate(&model.Epic{})
}

func loadEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Ping(context *gin.Context) {
	context.JSON(http.StatusCreated, gin.H{"msg": "pong"})
}

func serveApplication() {
	router := gin.Default()
	router.GET("/ping", Ping)

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	protectedRoutes := router.Group("/api")

	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	routeProject(router)
	routeSprint(router)
	routeIssue(router)
	routeSwimlane(router)
	routeUser(router)
	routeEpic(router)

	router.Run(":8000")
	fmt.Println("Server running on port 8000")
}

func routeProject(router *gin.Engine) {
	projectRoutes := router.Group("/project")
	projectRoutes.Use(middleware.JWTAuthMiddleware())

	projectRoutes.POST("/add", controller.AddProject)
	projectRoutes.POST("/add-user", controller.AddUserToProject)
	projectRoutes.POST("/user-owned", controller.GetOwnProjects)
	projectRoutes.POST("/user-joined", controller.GetJoinedProjects)
	projectRoutes.POST("/get-by-id", controller.GetProject)
	projectRoutes.POST("/update", controller.UpdateProjectInfo)
}

func routeSprint(router *gin.Engine) {
	projectRoutes := router.Group("/sprint")
	projectRoutes.Use(middleware.JWTAuthMiddleware())

	projectRoutes.POST("/add", controller.AddSprint)
	projectRoutes.POST("/start", controller.StartSprint)
	projectRoutes.POST("/end", controller.EndSprint)
}

func routeIssue(router *gin.Engine) {
	projectRoutes := router.Group("/issue")
	projectRoutes.Use(middleware.JWTAuthMiddleware())

	projectRoutes.POST("/add", controller.AddIssue)
	projectRoutes.POST("/delete", controller.DeleteIssue)
	projectRoutes.POST("/update", controller.UpdateIssue)
	projectRoutes.POST("/move-sprint", controller.MoveIssueSprint)
	projectRoutes.POST("/move-swimlane", controller.MoveIssueSwimlane)
}

func routeSwimlane(router *gin.Engine) {
	projectRoutes := router.Group("/swimlane")
	projectRoutes.Use(middleware.JWTAuthMiddleware())

	projectRoutes.POST("/add", controller.AddSwimlane)
}

func routeUser(router *gin.Engine) {
	projectRoutes := router.Group("/user")
	projectRoutes.Use(middleware.JWTAuthMiddleware())

	projectRoutes.POST("/get", controller.GetUser)
	projectRoutes.POST("/update", controller.UpdateUserInfo)
	projectRoutes.POST("/change-password", controller.ChangePassword)
}

func routeEpic(router *gin.Engine) {
	projectRoutes := router.Group("/epic")
	projectRoutes.Use(middleware.JWTAuthMiddleware())

	projectRoutes.POST("/add", controller.AddEpic)
}
