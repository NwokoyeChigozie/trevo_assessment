package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gregoflash05/trove/routes"
	"github.com/gregoflash05/trove/utils"
	"github.com/joho/godotenv"

	"github.com/rs/cors"
)

type App struct {
	Port string
}

func (app *App) Run() error {

	if err := utils.ConnectToDB(os.Getenv("CLUSTER_URL")); err != nil {
		return errors.New("Could not connect to MongoDB")
	}

	// transporter
	handler := routes.NewHandler()
	handler.SetupRoutes()

	c := cors.AllowAll()

	srv := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, c.Handler(handler.Router)),
		Addr:         ":" + app.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Trove Evaluation App running on port ", app.Port)

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {

	// load .env file if it exists
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	fmt.Println("Environment variables successfully loaded. Starting application...")

	// get PORT from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}

	app := App{Port: port}

	if err := app.Run(); err != nil {
		fmt.Println("Error occur while starting the Zuri Chat API.")
		log.Fatal(err)
	}
}
