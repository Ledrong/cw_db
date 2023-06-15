package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"time"
)

type setting struct {
	ServerHost string
	ServerPort string
	PgHost     string
	PgPort     string
	PgUser     string
	PgPassword string
	PgBase     string
}

type Server struct {
	http.Server
}

func NewServer(e *echo.Echo) *Server {
	return &Server{
		http.Server{
			Addr:              ":8080",
			Handler:           e,
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	log.Println("start serving in :8080")
	return s.ListenAndServe()
}

var cfg setting

func init() {
	//открыть файл конфигурации
	file, e := os.Open("apiii\\server\\setting.cfg")
	if e != nil {
		fmt.Println(e.Error())

		panic("не удалось открыть файл конфигурации")
	}

	defer file.Close()

	stat, e := file.Stat()
	if e != nil {
		fmt.Println(e.Error())
		panic("Не удалось прочитать инфу о файле конфигурации")
	}

	readByte := make([]byte, stat.Size())

	_, e = file.Read(readByte)
	if e != nil {
		fmt.Println(e.Error())

		panic("Не удалось прочитать файл конфигурации")
	}

	e = json.Unmarshal(readByte, &cfg)
	if e != nil {
		fmt.Println(e.Error())

		panic("Не удалось считать данные файла конфигурации")
	}

}

//func AuthMiddleware(next http.Handler, roles ...string) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("AuthMiddleware", r.URL.Path)
//		cookie, err := r.Cookie("jwt-token")
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusUnauthorized)
//			return
//		}
//		tokenString := cookie.Value
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//			}
//			return tokenPkg.SECRET, nil
//		})
//		if err != nil || !token.Valid {
//			http.Error(w, err.Error(), http.StatusUnauthorized)
//			return
//		}
//		claims, ok := token.Claims.(jwt.MapClaims)
//		if !ok {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//		if len(roles) > 0 {
//			roleMatch := false
//			for _, role := range roles {
//				if claims["role"] == role {
//					roleMatch = true
//					break
//				}
//			}
//			if !roleMatch {
//				http.Error(w, "Forbidden", http.StatusForbidden)
//				return
//			}
//		}
//		next.ServeHTTP(w, r)
//	})
//}

//func AuthMiddleware(next echo.HandlerFunc, roles ...string) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		cookie, err := c.Cookie("oreo")
//		if err != nil {
//			return c.String(http.StatusUnauthorized, err.Error())
//		}
//		cookieString := cookie.Value
//		token, err := jwt.Parse(cookieString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, errors.New("unexpected signing method")
//			}
//			return tokenPkg.SECRET, nil
//		})
//		if err != nil || !token.Valid {
//			return c.String(http.StatusUnauthorized, err.Error())
//		}
//		claims, ok := token.Claims.(jwt.MapClaims)
//		if !ok {
//			return c.String(http.StatusUnauthorized, "Unauthorized")
//		}
//		if len(roles) > 0 {
//			roleMatch := false
//			for _, role := range roles {
//				if claims["role"] == role {
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
