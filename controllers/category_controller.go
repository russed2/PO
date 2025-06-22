package controllers

import (
	"github.com/gin-gonic/gin"
	"online-shop/models"
	"gorm.io/gorm"
	"net/http"
)

// Получить все категории
func GetCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category
		if err := db.Find(&categories).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve categories"})
			return
		}
		c.JSON(200, categories)
	}
}

// Получить категорию по ID
func GetCategoryByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(200, category)
	}
}

// Создать новую категорию
func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category models.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&category).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create category"})
			return
		}

		c.JSON(200, category)
	}
}

// Обновить категорию
func UpdateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Category not found"})
			return
		}

		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&category).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update category"})
			return
		}

		c.JSON(200, category)
	}
}

// Удалить категорию
func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Category not found"})
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete category"})
			return
		}

		c.JSON(200, gin.H{"message": "Category deleted"})
	}
}
