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

	registerReadRoutes := func(group *echo.Group) {
		group.GET("/titles", titleHandler.GetTitleDetails)
		group.GET("/scan/:id", handlers.GetScanHandler)
		group.GET("/list/:id", handlers.GetListHandler)
		group.GET("/user/:id", handlers.GetUserHandler)
		group.GET("/article/:slug", handlers.GetArticleHandler)
		group.GET("/mangaList", handlers.GetMangaListHandler)
		group.GET("/latestTitles", handlers.GetLatestTitlesHandler)
		group.GET("/titlesByTag", handlers.GetTitlesByTagHandler)
		group.GET("/chapter", handlers.GetChapterHandler)
		group.GET("/latestChapters", handlers.GetLatestChaptersHandler)
		group.GET("/localTitles", handlers.GetLocalTitlesHandler)
	}

	api := e.Group("/api")
	registerReadRoutes(api)

	// Draft endpoint'leri normal endpoint'lerin ayni sorgularini kullanir;
	// middleware sadece Sanity perspektifini drafts olarak degistirir.
	draftAPI := e.Group("/api/draft", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request().WithContext(services.WithDrafts(c.Request().Context()))
			c.SetRequest(request)
			return next(c)
		}
	})
	registerReadRoutes(draftAPI)

	e.POST("/api/user", handlers.CreateUserHandler)

	e.Logger.Fatal(e.Start(":2611"))
}
