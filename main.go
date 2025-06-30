package main

import (
	"fmt"
	"online-shop/config"
	"online-shop/controllers"
	"online-shop/logger"
	"online-shop/middleware"
	"online-shop/models"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDatabase() *gorm.DB {
	// Подключение к базе данных
	dsn := "host=db user=admin password=123 dbname=PO2SEM port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to database:", err)
	}
	logger.Log.Info("Connected to PostgreSQL!")
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Добавим endpoint для проверки балансировки
	r.GET("/health", func(c *gin.Context) {
		serverID := os.Getenv("SERVER_ID")
		if serverID == "" {
			serverID = "app-main"
		}
		c.JSON(200, gin.H{
			"status":    "OK",
			"server_id": serverID,
			"message":   "Server is healthy",
		})
	})

	// Роуты без токена
	// Регистрация и логин
	r.POST("/users/register", controllers.RegisterUser(db))
	r.POST("/users/login", controllers.LoginUser(db))

	// Просмотр товаров и категорий
	r.GET("/products", controllers.GetProducts(db))
	r.GET("/products/search", controllers.SearchProducts(db))
	r.GET("/categories", controllers.GetCategories(db))
	r.GET("/categories/:id", controllers.GetCategoryByID(db))
	r.GET("/reviews/:product_id", controllers.GetReviewsByProductID(db))

	// Роуты с токеном
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Профиль пользователя
		auth.GET("/profile", controllers.GetUserProfile(db))
		auth.PUT("/profile", controllers.UpdateUserProfile(db))

		// Корзина
		auth.POST("/cart", controllers.AddToCart(db))
		auth.GET("/cart/:user_id", controllers.GetCartItems(db))
		auth.DELETE("/cart/:id", controllers.RemoveFromCart(db))
		auth.DELETE("/cart/clear/:user_id", controllers.ClearCart(db))

		// Заказы
		auth.POST("/orders", controllers.CreateOrder(db))
		auth.GET("/orders", controllers.GetOrders(db))
		auth.GET("/orders/:id", controllers.GetOrderByID(db))

		// Отзывы
		auth.POST("/reviews", controllers.AddReview(db))
		auth.DELETE("/reviews/:id", controllers.DeleteReview(db))
	}

	// Роуты для админов
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		// Управление пользователями
		admin.GET("/users", controllers.GetAllUsers(db))
		admin.GET("/users/:id", controllers.GetUserByID(db))
		admin.PUT("/users/:id", controllers.UpdateUser(db))
		admin.DELETE("/users/:id", controllers.DeleteUser(db))

		// Управление категориями
		admin.POST("/categories", controllers.CreateCategory(db))
		admin.PUT("/categories/:id", controllers.UpdateCategory(db))
		admin.DELETE("/categories/:id", controllers.DeleteCategory(db))

		// Управление товарами
		admin.POST("/products", controllers.CreateProduct(db))
		admin.PUT("/products/:id", controllers.UpdateProduct(db))
		admin.DELETE("/products/:id", controllers.DeleteProduct(db))

		// Управление заказами
		admin.PUT("/orders/:id", controllers.UpdateOrderStatus(db))
		admin.DELETE("/orders/:id", controllers.DeleteOrder(db))
	}

	// Логируем все роуты
	logger.Log.Info("Registered routes:")
	for _, ri := range r.Routes() {
		logger.Log.Infof("%s %s", ri.Method, ri.Path)
	}

	return r
}

func main() {
	// Загружаем конфигурацию из .env файла
	config.LoadConfig()
	logger.Init()

	logger.Log.Info("Application starting...")

	// Демонстрация загрузки YAML конфигурации (Практика №5, часть 1)
	yamlConfig, err := config.LoadYAMLConfig("config.yaml")
	if err != nil {
		logger.Log.Warn("Could not load config.yaml:", err)
		logger.Log.Info("Using default configuration...")
	} else {
		logger.Log.Infof("Loaded YAML config for app: %s v%s", yamlConfig.App.Name, yamlConfig.App.Version)
	}

	// Подключаемся к базе данных
	db := connectToDatabase()

	// Автоматическая миграция моделей
	logger.Log.Info("Starting database migration...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.CartItem{},
		&models.Review{},
		&models.UserToken{},
		&models.Session{},
		&models.Category{},
	)
	if err != nil {
		logger.Log.Fatal("Failed to migrate models:", err)
	}
	logger.Log.Info("Database migration completed successfully")

	// Выводим какой сервер сейчас работает
	serverID := os.Getenv("SERVER_ID")
	if serverID == "" {
		serverID = "app-main"
	}
	fmt.Printf("Starting server: %s\n", serverID)

	// Запуск сервера
	r := setupRouter(db)
	r.Run(":8080")
}
