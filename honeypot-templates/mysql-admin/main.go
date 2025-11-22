package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	payload := os.Getenv("CHM_PAYLOAD")
	srcIP := os.Getenv("CHM_SRC_IP")

	db := GenerateFakeDB()

	// Log to stdout so Aggregator can catch it
	e.Logger.Infof("HONEYPOT STARTED — src=%s payload=%s db=%s", srcIP, payload, db.Name)

	// --- ROUTES ---
	e.GET("/", func(c echo.Context) error {
		SlowDown()
		return c.HTML(http.StatusOK,
			"<h3>MySQL Admin Panel</h3><p>Login required.</p>")
	})

	e.POST("/login", func(c echo.Context) error {
		SlowDown()

		// user := c.FormValue("username")
		pass := c.FormValue("password")

		// Fake brute force lockout increase
		if len(pass) < 3 {
			return c.String(200, "Access denied: Invalid credentials")
		}

		return c.String(200, "Login failed: user does not exist")
	})

	e.POST("/query", func(c echo.Context) error {
		SlowDown()

		q := c.FormValue("q")
		if strings.Contains(strings.ToLower(q), "select") {
			return c.String(500, GenerateSQLError(q))
		}
		return c.String(200, "Query executed (no results)")
	})

	e.GET("/db/info", func(c echo.Context) error {
		SlowDown()

		return c.JSON(200, map[string]interface{}{
			"database": db.Name,
			"tables":   db.Tables,
			"source":   srcIP,
			"payload":  payload,
		})
	})

	e.Start(":8080")
}
