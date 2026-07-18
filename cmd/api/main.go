package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/falsisdev/mangile-backend/internal/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[HATA]: .env yüklenemedi: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://mangile.vercel.app"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "[✅]: Sunucu başarıyla başlatıldı.")
	})

	e.GET("/api/manga/:id", handlers.GetMangaHandler)
	e.GET("/api/manga/:id/recommendations", handlers.GetMangaRecommendationsHandler)
	e.GET("/api/scan/:id", handlers.GetScanHandler)
	e.GET("/api/list/:id", handlers.GetListHandler)
	e.GET("/api/user/:id", handlers.GetUserHandler)
	e.GET("/api/article/:slug", handlers.GetArticleHandler)
	e.GET("/api/mangaList", handlers.GetMangaListHandler)
	e.GET("/api/sanityList", handlers.GetSanityListHandler)

	e.Logger.Fatal(e.Start(":3001"))
}
