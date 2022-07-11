package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"test/test_app/app/service/logger"
	"time"

	"bitbucket.org/liamstask/goose/lib/goose"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

// Init : Initializes the database migrations
func Init(ctx context.Context) (db *gorm.DB, err error) {
	log := logger.Logger(ctx)

	dbUserName := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	maxIdleConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTION"))
	maxOpenConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTION"))
	connectionMaxLifetime, _ := strconv.Atoi(os.Getenv("DB_CONNECTION_MAX_LIFETIME"))

	var dbURI string
	var dialect goose.SqlDialect

	dbURI = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUserName, dbPassword, dbName)

	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Println("failed to connect.", dbURI, err)
		log.Fatalf("Failed to connect to DB", dbURI, err.Error())
		os.Exit(1)
	}
	dialect = &goose.PostgresDialect{}

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Not able to fetch the working directory")
		log.Fatalf("Not able to fetch the working directory")
		os.Exit(1)
	}

	if strings.Contains(workingDir, "core_service_api") {
		tempDir := strings.Split(workingDir, "/core_service_api")
		if len(tempDir) > 1 {
			workingDir = tempDir[0]
		}
	}
	db.DB().SetMaxIdleConns(maxIdleConnections)
	db.DB().SetMaxOpenConns(maxOpenConnections)
	db.DB().SetConnMaxLifetime(time.Hour * time.Duration(connectionMaxLifetime))
	db.SingularTable(true)

	workingDir = workingDir + "/test_app/app/db/migrations"
	migrateConf := &goose.DBConf{
		MigrationsDir: workingDir,
		Driver: goose.DBDriver{
			Name:    "postgres",
			OpenStr: dbURI,
			Import:  "github.com/lib/pq",
			Dialect: dialect,
		},
	}
	log.Infof("Fetching the most recent DB version")
	latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)
	if err != nil {
		log.Errorf("Unable to get recent goose db version", err)
	}
	fmt.Println(" Most recent DB version ", latest)
	log.Infof("Running the migrations on db", workingDir)
	err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, db.DB())
	if err != nil {
		log.Fatalf("Error while running migrations", err)
		os.Exit(1)
	}
	return
}

func New(dbConn *gorm.DB) *DBService {
	return &DBService{
		DB: dbConn,
	}
}

type DBService struct {
	DB *gorm.DB
}

// GetDB : Get an instance of DB to connect to the database connection pool
func (d DBService) GetDB() *gorm.DB {
	return d.DB
}
