package main

import (
	"log"
	"net/http"

	"github.com/falsisdev/mangile-backend/internal/handlers"
	"github.com/falsisdev/mangile-backend/internal/services"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[HATA]: .env yüklenemedi: %v", err)
	}

	jikanService := services.NewJikanService()
	aniListService := services.NewAniListService()

	titleWorker := services.NewTitleDetailsWorker(jikanService, aniListService)
	titleHandler := handlers.NewTitleHandler(titleWorker)

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

	e.GET("/api/titles", titleHandler.GetTitleDetails)
	e.GET("/api/scan/:id", handlers.GetScanHandler)
	e.GET("/api/list/:id", handlers.GetListHandler)
	e.GET("/api/user/:id", handlers.GetUserHandler)
	e.GET("/api/article/:slug", handlers.GetArticleHandler)
	e.GET("/api/mangaList", handlers.GetMangaListHandler)
	e.GET("/api/latestTitles", handlers.GetLatestTitlesHandler)
	e.GET("/api/titlesByTag", handlers.GetTitlesByTagHandler)
	e.GET("/api/chapter", handlers.GetChapterHandler)
	e.GET("/api/latestChapters", handlers.GetLatestChaptersHandler)

	e.Logger.Fatal(e.Start(":2611"))
}
