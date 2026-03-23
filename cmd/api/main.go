package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"task-manager-api/internal/config"
	"task-manager-api/internal/handler"
	"task-manager-api/internal/middleware"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/service"
	"task-manager-api/internal/worker"

	_ "task-manager-api/docs"
)

// @title           Task Manager API
// @version         1.0
// @description     RESTful API untuk manajemen tugas dengan implementasi Clean Architecture dan AI Assistant.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.email  support@example.com
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @host      localhost:8080
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.ConnectDB()

	// Dependency Injection: Auth
	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Dependency Injection: Task
	taskRepo := repository.NewTaskRepository(config.DB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// Dependency Injection: AI Assistant (Baru)
	aiService := service.NewAIService()
	aiHandler := handler.NewAIHandler(aiService)

	// Start Task Worker di background
	worker.StartTaskWorker(taskRepo)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Mengizinkan akses HANYA dari frontend Next.js Anda
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Caching aturan CORS selama 12 jam agar lebih cepat
	}))

	router.Static("/uploads", "./uploads")

	// Swagger Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "API is running!"})
		})

		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Rute yang butuh Login (JWT)
		protected := v1.Group("")
		protected.Use(middleware.RequireAuth())
		{
			// Endpoint Task
			taskRoutes := protected.Group("/tasks")
			{
				taskRoutes.POST("", taskHandler.CreateTask)
				taskRoutes.GET("", taskHandler.GetTasks)
				taskRoutes.GET("/:id", taskHandler.GetTaskByID)
				taskRoutes.PUT("/:id", taskHandler.UpdateTask)
				taskRoutes.DELETE("/:id", taskHandler.DeleteTask)
				taskRoutes.POST("/:id/subtasks", taskHandler.AddSubTasks)
				taskRoutes.POST("/:id/upload", taskHandler.UploadAttachment)
				taskRoutes.POST("/:id/collaborators", taskHandler.AddCollaborator)
			}

			// Endpoint AI Assistant (Baru)
			aiRoutes := protected.Group("/ai")
			{
				aiRoutes.POST("/generate-tasks", aiHandler.GenerateTaskBreakdown)
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Mulai menjalankan server di port :%s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}