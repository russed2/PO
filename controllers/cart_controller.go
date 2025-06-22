package controllers

import (
	"github.com/gin-gonic/gin"
	"online-shop/models"
	"gorm.io/gorm"
)

// Добавить товар в корзину
func AddToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cartItem models.CartItem
		if err := c.ShouldBindJSON(&cartItem); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Проверяем, есть ли товар в корзине
		var existingCartItem models.CartItem
		if err := db.Where("user_id = ? AND product_id = ?", cartItem.UserID, cartItem.ProductID).First(&existingCartItem).Error; err == nil {
			// Если товар уже есть, обновляем количество
			existingCartItem.Quantity += cartItem.Quantity
			if err := db.Save(&existingCartItem).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to update cart"})
				return
			}
			c.JSON(200, existingCartItem)
			return
		}

		// Добавляем новый товар в корзину
		if err := db.Create(&cartItem).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to add to cart"})
			return
		}

		c.JSON(200, cartItem)
	}
}

// Получить все товары в корзине пользователя
func GetCartItems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		var cartItems []models.CartItem
		if err := db.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve cart items"})
			return
		}
		c.JSON(200, cartItems)
	}
}

// Удалить товар из корзины
func RemoveFromCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var cartItem models.CartItem
		if err := db.First(&cartItem, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Cart item not found"})
			return
		}

		if err := db.Delete(&cartItem).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete cart item"})
			return
		}

		c.JSON(200, gin.H{"message": "Cart item deleted"})
	}
}

// Очистить корзину пользователя
func ClearCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if err := db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to clear cart"})
			return
		}
		c.JSON(200, gin.H{"message": "Cart cleared"})
	}
}
