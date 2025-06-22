package controllers

import (
	"github.com/gin-gonic/gin"
	"online-shop/models"
	"gorm.io/gorm"
	"net/http"
)

// Получить все товары
func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		if err := db.Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve products"})
			return
		}
		c.JSON(200, products)
	}
}

// Функция для поиска товаров по ID или имени
func SearchProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем параметры поиска из query-параметров
		id := c.DefaultQuery("id", "")   
		name := c.DefaultQuery("name", "") 

		var products []models.Product

		// Если передан ID, ищем по ID
		if id != "" {
			if err := db.Where("id = ?", id).Find(&products).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to retrieve products by ID"})
				return
			}
		} else if name != "" {
			// Если передано имя, ищем по имени
			if err := db.Where("name ILIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to retrieve products by name"})
				return
			}
		} else {
			// Если ничего не передано, выводим все товары
			if err := db.Find(&products).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to retrieve products"})
				return
			}
		}

		c.JSON(200, products)
	}
}

// Создать новый продукт
func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&product).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(200, product)
	}
}

// Обновить продукт
func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&product).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update product"})
			return
		}

		c.JSON(200, product)
	}
}

// Удалить продукт
func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Product not found"})
			return
		}

		if err := db.Delete(&product).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete product"})
			return
		}

		c.JSON(200, gin.H{"message": "Product deleted"})
	}
}
