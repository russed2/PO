package controllers

import (
	"github.com/gin-gonic/gin"
	"online-shop/models"
	"gorm.io/gorm"
	"net/http"
)

// Создание нового заказа
func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Создание заказа в базе данных
		if err := db.Create(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create order"})
			return
		}

		c.JSON(200, order)
	}
}

// Получить заказ по ID
func GetOrderByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id") // Получаем ID из параметров запроса
		var order models.Order
		// Ищем заказ по ID
		if err := db.First(&order, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}
		// Возвращаем найденный заказ
		c.JSON(200, order)
	}
}

// Удалить заказ
func DeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var order models.Order
		if err := db.First(&order, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}

		if err := db.Delete(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete order"})
			return
		}

		c.JSON(200, gin.H{"message": "Order deleted"})
	}
}

// Обновить статус заказа
func UpdateOrderStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var order models.Order
		if err := db.First(&order, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}

		// Проверяем, что статус допустим
		status := c.DefaultQuery("status", "")
		if status != "оплачен" && status != "ожидает оплаты" {
			c.JSON(400, gin.H{"error": "Invalid status"})
			return
		}

		// Обновляем статус заказа
		order.Status = status
		if err := db.Save(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update order status"})
			return
		}

		c.JSON(200, order)
	}
}

// Получение всех заказов
func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orders []models.Order
		if err := db.Find(&orders).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve orders"})
			return
		}
		c.JSON(200, orders)
	}
}
