package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Секретный ключ для подписи токенов
var jwtSecret = []byte("secret-key")

// AuthMiddleware - проверяет, что пользователь прислал правильный токен
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization из HTTP запроса
		authHeader := c.GetHeader("Authorization")

		// Если заголовка нет - возвращаем ошибку
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Нет токена авторизации"})
			c.Abort() // Останавливаем выполнение, не пускаем к контроллеру
			return
		}

		// Убираем слово "Bearer " из начала, чтобы остался только токен
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Проверяем токен: настоящий ли он и не истек ли срок
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильный токен"})
			c.Abort()
			return
		}

		// Извлекаем данные пользователя из токена
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Сохраняем данные пользователя, чтобы контроллеры могли их использовать
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			c.Set("admin", claims["admin"])
		}

		c.Next()
	}
}

// AdminMiddleware - проверяет, что пользователь - администратор
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем роль пользователя
		admin, exists := c.Get("admin")
		if !exists || admin != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "Нужны права администратора"})
			c.Abort()
			return
		}

		c.Next()
	}
}
