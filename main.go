package main

import (
	"database/sql"
	"fmt"
	"lawoffice/config"
	"lawoffice/operations"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	db       *sql.DB
	dbconfig config.Database
	server   config.Server
)

func InitDB() (*sql.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",dbconfig.Username,dbconfig.Password,dbconfig.Host,dbconfig.Port,dbconfig.Dbname)
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Initialize configuration by calling the functions
	dbconfig = config.DbConfig()   // Call the function to get the database configuration
	server = config.ServerConfig() // Call the function to get the server configuration
	e := echo.New()
	e.Logger.SetOutput(log.Writer())
	e.Logger.SetLevel(3)
	e.Use(middleware.Logger())
	// Initialize the database
	var err error
	db, err = InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	// Pass the database connection to the operations package
	operations.SetDB(db)
	operations.DashboardRouter(e)
	operations.ClientsRouter(e)
	operations.CasesRouter(e)
	operations.AppointRouter(e)
	operations.PowerRouter(e)
	operations.SessionsRouter(e)
	operations.InovicesRouter(e)
	operations.LawyersRouter(e)
	// Server configuration
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", server.Host, server.Port)))
}
