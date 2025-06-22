package main

import (
	"fmt"
	"log"
	"online-shop/controllers"
	"online-shop/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectToDatabase() *gorm.DB {
	// Подключение к базе данных
	dsn := "host=db user=admin password=123 dbname=PO2SEM port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to PostgreSQL!")
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	//Роуты для пользователей
	r.POST("/users/register", controllers.RegisterUser(db))
	r.POST("/users/login", controllers.LoginUser(db))
	r.GET("/users", controllers.GetAllUsers(db))
	r.GET("/users/:id", controllers.GetUserByID(db))
	r.PUT("/users/:id", controllers.UpdateUser(db))
	r.DELETE("/users/:id", controllers.DeleteUser(db))

	// Роуты для категорий
	r.GET("/categories", controllers.GetCategories(db))         // Получить все категории
	r.GET("/categories/:id", controllers.GetCategoryByID(db))   // Получить категорию по ID
	r.POST("/categories", controllers.CreateCategory(db))       // Создать новую категорию
	r.PUT("/categories/:id", controllers.UpdateCategory(db))    // Обновить категорию
	r.DELETE("/categories/:id", controllers.DeleteCategory(db)) // Удалить категорию

	// Роуты для продуктов
	r.GET("/products", controllers.GetProducts(db))          // Получить все товары
	r.GET("/products/:id", controllers.SearchProducts(db))   // Получить продукт по ID или имени
	r.POST("/products", controllers.CreateProduct(db))       // Создать новый продукт
	r.PUT("/products/:id", controllers.UpdateProduct(db))    // Обновить продукт
	r.DELETE("/products/:id", controllers.DeleteProduct(db)) // Удалить продукт

	// Роуты для заказов
	r.POST("/orders", controllers.CreateOrder(db))          // Создать новый заказ
	r.GET("/orders", controllers.GetOrders(db))             // Получить все заказы
	r.GET("/orders/:id", controllers.GetOrderByID(db))      // Получить заказ по ID
	r.PUT("/orders/:id", controllers.UpdateOrderStatus(db)) // Обновить статус заказа
	r.DELETE("/orders/:id", controllers.DeleteOrder(db))    // Удалить заказ

	// Роуты для корзины
	r.POST("/cart", controllers.AddToCart(db))                  // Добавить товар в корзину
	r.GET("/cart/:user_id", controllers.GetCartItems(db))       // Получить все товары в корзине
	r.DELETE("/cart/:id", controllers.RemoveFromCart(db))       // Удалить товар из корзины
	r.DELETE("/cart/clear/:user_id", controllers.ClearCart(db)) // Очистить корзину

	// Роуты для отзывов
	r.POST("/reviews", controllers.AddReview(db))                        // Добавить отзыв
	r.GET("/reviews/:product_id", controllers.GetReviewsByProductID(db)) // Получить отзывы по ID товара
	r.DELETE("/reviews/:id", controllers.DeleteReview(db))               // Удалить отзыв

	for _, ri := range r.Routes() {
		fmt.Println(ri.Method, ri.Path)
	}

	return r
}

func main() {
	// Подключаемся к базе данных
	db := connectToDatabase()

	// Автоматическая миграция моделей
	err := db.AutoMigrate(
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
		log.Fatal("Failed to migrate models:", err)
	}

	// Запуск сервера
	r := setupRouter(db)
	r.Run(":8080")
}
