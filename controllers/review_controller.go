package controllers

import (
	"github.com/gin-gonic/gin"
	"online-shop/models"
	"gorm.io/gorm"
)

// Добавить отзыв о товаре
func AddReview(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var review models.Review
		if err := c.ShouldBindJSON(&review); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Добавляем новый отзыв
		if err := db.Create(&review).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to add review"})
			return
		}

		c.JSON(200, review)
	}
}

// Получить все отзывы для товара
func GetReviewsByProductID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("product_id")
		var reviews []models.Review
		if err := db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve reviews"})
			return
		}
		c.JSON(200, reviews)
	}
}

// Удалить отзыв
func DeleteReview(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var review models.Review
		if err := db.First(&review, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Review not found"})
			return
		}

		if err := db.Delete(&review).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete review"})
			return
		}

		c.JSON(200, gin.H{"message": "Review deleted"})
	}
}
