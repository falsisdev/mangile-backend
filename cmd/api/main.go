package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/falsisdev/mangile-backend/internal/handlers"
)

// .env'den bilgi çekmek için: os.Getenv("SANITY_TOKEN") (os paketi gerekiyor)
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[HATA]: .env yüklenemedi: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "[✅]: Sunucu başarıyla başlatıldı.")
	})

	e.GET("/api/manga/:id", handlers.GetMangaHandler)
	e.GET("/api/lightNovel/:id", handlers.GetLightNovelHandler)
	e.GET("/api/scan/:id", handlers.GetScanHandler)
	e.GET("/api/list/:id", handlers.GetListHandler)
	e.GET("/api/user/:id", handlers.GetUserHandler)
	e.GET("/api/article/:slug", handlers.GetArticleHandler)

	e.Logger.Fatal(e.Start(":3000"))
}
