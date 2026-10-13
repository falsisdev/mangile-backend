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
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("cmd/api/.env"); err != nil {
			log.Printf("[UYARI]: .env yuklenemedi: %v", err)
		}
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
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Logto-Id"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "[✅]: Sunucu başarıyla başlatıldı.")
	})

	registerReadRoutes := func(group *echo.Group) {
		group.GET("/titles", titleHandler.GetTitleDetails)
		group.GET("/scan/:id", handlers.GetScanHandler)
		group.GET("/list/:id", handlers.GetListHandler)
		group.GET("/lists", handlers.GetPublicListsHandler)
		group.GET("/user/:id", handlers.GetUserHandler)
		group.GET("/user/:id/lists", handlers.GetUserListsHandler)
		group.GET("/user/:id/liked-lists", handlers.GetLikedListsHandler)
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

	draftAPI := e.Group("/api/draft", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request().WithContext(services.WithDrafts(c.Request().Context()))
			c.SetRequest(request)
			return next(c)
		}
	})
	registerReadRoutes(draftAPI)

	e.POST("/api/user", handlers.CreateUserHandler)

	api.POST("/list", handlers.CreateListHandler)
	api.PATCH("/list/:id", handlers.UpdateListHandler)
	api.DELETE("/list/:id", handlers.DeleteListHandler)
	api.POST("/list/:id/items", handlers.AddListItemHandler)
	api.DELETE("/list/:id/items/:key", handlers.RemoveListItemHandler)
	api.PUT("/list/:id/items", handlers.UpdateListItemsHandler)
	api.POST("/list/:id/like", handlers.ToggleLikeListHandler)

	trackerSyncService := services.NewTrackerSyncService()
	malService := services.NewMALService()
	aniListTrackerService := services.NewAniListTrackerService()
	mangaDexService := services.NewMangaDexTrackerService()
	trackerHandler := handlers.NewTrackerHandler(trackerSyncService, malService, aniListTrackerService, mangaDexService)

	api.GET("/user/:id/history", handlers.GetHistoryHandler)
	api.POST("/user/:id/history", handlers.UpsertHistoryHandler)
	api.DELETE("/user/:id/history/:seriesId", handlers.DeleteHistoryItemHandler)
	api.DELETE("/user/:id/history", handlers.ClearHistoryHandler)
	api.POST("/user/:id/history/sync", handlers.SyncHistoryHandler)

	api.GET("/user/:id/trackers", trackerHandler.GetTrackersHandler)
	api.GET("/user/:id/trackers/:seriesId/status", trackerHandler.GetSeriesTrackerStatusHandler)
	api.POST("/user/:id/trackers/bind", trackerHandler.BindTrackerHandler)
	api.DELETE("/user/:id/trackers/:seriesId/:provider", trackerHandler.UnbindTrackerHandler)
	api.POST("/user/:id/trackers/update", trackerHandler.UpdateTrackerStatusHandler)
	api.POST("/user/:id/trackers/progress", trackerHandler.SyncTrackerProgressHandler)
	api.GET("/tracker/search", trackerHandler.SearchTrackerTitlesHandler)

	api.GET("/sync/oauth/:provider/url", trackerHandler.OAuthInitiateHandler)
	api.GET("/sync/oauth/:provider/callback", trackerHandler.OAuthCallbackHandler)
	api.DELETE("/user/:id/integrations/:provider", trackerHandler.DisconnectIntegrationHandler)

	e.Logger.Fatal(e.Start(":2611"))
}
