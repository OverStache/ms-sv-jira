package main

import (
	"log"
	"ms-sv-jira/dbconn"
	"ms-sv-jira/helper/logger"
	"ms-sv-jira/middleware"
	handler "ms-sv-jira/module/delivery/http"
	"ms-sv-jira/module/repository"
	"ms-sv-jira/module/usecase"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConn := dbconn.DB()
	l := logger.L
	appPort := os.Getenv("PORT")
	timeout, _ := strconv.Atoi(os.Getenv("APP_TIMEOUT"))
	timeoutContext := time.Duration(timeout) * time.Second
	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.Log)
	// e.Use(middL.Custom)
	e.Use(_echoMiddleware.CORSWithConfig(_echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentDisposition, echo.HeaderAccessControlRequestHeaders, "ngrok-skip-browser-warning"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	repo := repository.NewRepository(dbConn, l)
	usecase := usecase.NewUsecase(repo, timeoutContext, l)
	handler.NewHandler(e, middL, usecase, l)

	log.Fatal(e.Start(":" + appPort))
}
