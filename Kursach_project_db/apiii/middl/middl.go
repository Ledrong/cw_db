package middl

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

//func AuthMiddleware(next echo.HandlerFunc, roles ...string) echo.MiddlewareFunc {
//	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			session, err := c.Cookie("oreo")
//			if err != nil {
//				return c.String(http.StatusUnauthorized, err.Error())
//			}
//			sessionString := session.Value
//			claims := make(map[string]interface{})
//			err = json.Unmarshal([]byte(sessionString), &claims)
//			if err != nil {
//				return c.String(http.StatusUnauthorized, err.Error())
//			}
//			if len(roles) > 0 {
//				roleMatch := false
//				for _, role := range roles {
//					if claims["role"] == role {
//						roleMatch = true
//						break
//					}
//				}
//				if !roleMatch {
//					return c.String(http.StatusForbidden, "Forbidden")
//				}
//			}
//			return next(c)
//		}
//	}
//}

//func AuthMiddleware(next echo.HandlerFunc, roles ...string) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		str := c.Request().Header.Get("Authorization")
//		if len(str) == 0 {
//			return c.String(http.StatusUnauthorized, "Unauthorized")
//		}
//		// Проверка токена может быть выполнена другими способами (например, база данных),
//		// здесь просто пример с хранением токена в переменной окружения.
//		if str != os.Getenv("JWT_TOKEN") {
//			return c.String(http.StatusUnauthorized, "Unauthorized")
//		}
//		claims := make(map[string]interface{})
//		role, ok := claims["role"].(string)
//		if !ok {
//			return c.String(http.StatusUnauthorized, "Unauthorized")
//		}
//		if len(roles) > 0 {
//			roleMatch := false
//			for _, r := range roles {
//				if role == r {
//					roleMatch = true
//					break
//				}
//			}
//			if !roleMatch {
//				return c.String(http.StatusForbidden, "Forbidden")
//			}
//		}
//		return next(c)
//	}
//}

func GetIdFromCookie(c echo.Context) (int64, error) {
	cookie, err := c.Cookie("oreo")
	if err != nil {
		return 0, c.JSON(http.StatusBadRequest, err.Error())
	}
	if cookie == nil {
		return 0, c.String(http.StatusUnauthorized, "Unauthorized")
	}
	var cookierole string
	var id int64
	cookiestring := cookie.Value
	fmt.Sscanf(cookiestring, "%s %d", &cookierole, &id)

	return id, nil
}

func AuthMiddleware(next echo.HandlerFunc, roles ...string) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("oreo")
		if err != nil {
			return c.JSON(http.StatusBadRequest, "no roots")
		}
		if cookie == nil {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		var cookierole string
		var id int64
		cookiestring := cookie.Value
		fmt.Sscanf(cookiestring, "%s %d", &cookierole, &id)

		if len(roles) != 0 {
			flag := false
			for _, role := range roles {
				if cookierole == role {
					flag = true
					break
				}
			}
			if !flag {
				return c.String(http.StatusForbidden, "Forbidden")
			}
		}
		return next(c)
	}
}

//func AuthMiddleware(next echo.HandlerFunc, roles ...string) echo.MiddlewareFunc {
//	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			session, err := c.Cookie("oreo")
//			if err != nil {
//				return c.String(http.StatusUnauthorized, err.Error())
//			}
//			sessionString := session.Value
//			claims := make(map[string]interface{})
//			err = json.Unmarshal([]byte(sessionString), &claims)
//			if err != nil {
//				return c.String(http.StatusUnauthorized, err.Error())
//			}
//			if len(roles) > 0 {
//				roleMatch := false
//				for _, role := range roles {
//					if claims["role"] == role {
//						roleMatch = true
//						break
//					}
//				}
//				if !roleMatch {
//					return c.String(http.StatusForbidden, "Forbidden")
//				}
//			}
//			return next(echo.HandlerFunc(handlerFunc))(c)
//		}
//	}
//}
