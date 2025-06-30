package controllers

import (
	"net/http"
	"online-shop/models"
	"online-shop/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Регистрация пользователя
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось захешировать пароль"})
			return
		}
		input.Password = hashedPassword

		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации пользователя"})
			return
		}

		c.JSON(http.StatusOK, input)
	}
}

// Вход пользователя
func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Поиск пользователя по email
		var user models.User
		if err := db.Where("email_user = ?", credentials.Email).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "Неверные учетные данные"})
			return
		}

		// Проверка пароля
		if !utils.CheckPasswordHash(credentials.Password, user.Password) {
			c.JSON(401, gin.H{"error": "Неверные учетные данные"})
			return
		}

		// Создаем JWT токен с ролью админа
		token, err := utils.GenerateToken(user.ID, user.Name, user.Admin)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ошибка создания токена"})
			return
		}

		// Возвращаем токен и данные пользователя
		c.JSON(200, gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"admin": user.Admin, // Клиент узнает свою роль
			},
		})
	}
}

// Получить всех пользователей
func GetAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to get users"})
			return
		}
		c.JSON(200, users)
	}
}

// Получить пользователя по ID
func GetUserByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// Получить профиль текущего пользователя
func GetUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем данные пользователя из JWT токена
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
			return
		}

		// Получаем полную информацию о пользователе из БД
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		// Возвращаем данные профиля (без пароля!)
		c.JSON(http.StatusOK, gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"admin":      user.Admin,
			"created_at": user.CreatedAt,
			"message":    "Профиль загружен успешно",
		})
	}
}

// Обновить профиль текущего пользователя (для авторизованного пользователя)
func UpdateUserProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID из токена
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
			return
		}

		// Находим пользователя
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		// Структура для получения обновлений
		var updateData struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Обновляем поля если они переданы
		if updateData.Name != "" {
			user.Name = updateData.Name
		}
		if updateData.Email != "" {
			user.Email = updateData.Email
		}
		if updateData.Password != "" {
			hashedPassword, err := utils.HashPassword(updateData.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
				return
			}
			user.Password = hashedPassword
		}

		// Сохраняем изменения
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления профиля"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"admin":   user.Admin,
			"message": "Профиль обновлен успешно",
		})
	}
}

// Обновить пользователя (для админов)
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		var input models.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Password != "" {
			hashed, _ := utils.HashPassword(input.Password)
			input.Password = hashed
		} else {
			input.Password = user.Password
		}

		input.ID = user.ID
		if err := db.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления пользователя"})
			return
		}

		c.JSON(http.StatusOK, input)
	}
}

// Удалить пользователя
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления пользователя"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Пользователь удалён"})
	}
}
