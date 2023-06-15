package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func connect() error {
	var e error
	db, e = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPassword, cfg.PgBase))
	if e != nil {
		return e
	}

	return nil
}

//func AccessLogMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("accessLogMiddleware", r.URL.Path)
//		start := time.Now()
//		next.ServeHTTP(w, r)
//		fmt.Printf("[%s] %s, %s %s\n",
//			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
//	})
//}
