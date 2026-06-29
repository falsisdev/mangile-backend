package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// .env'den bilgi çekmek için: os.Getenv("SANITY_TOKEN") (os paketi gerekiyor)
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[HATA]: .env dosyası yüklenirken bir hata oluştu.")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "[✅]: Sunucu başarıyla başlatıldı.")
	})

	e.Logger.Fatal(e.Start(":3000"))
}
