package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	APPConfig *Config)

type Config struct {
	APPPort string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	JWTSecret string
	JWTExpireMinutes string
	JWTRefreshToken string
	JWTExpire string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	APPConfig = &Config{
		APPPort: getEnv("PORT", "3030"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName: getEnv("DB_NAME", "wpu_goreact"),
		JWTSecret: getEnv("JWT_SECRET", "supersecret"),
		JWTExpireMinutes: getEnv("JWT_EXPIRY_MINUTES", "6000"),
		JWTRefreshToken: getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
		JWTExpire: getEnv("JWT_EXPIRED", "1j"),
	}
}

func getEnv(key string, fallback string) string {
	value, exist := os.LookupEnv(key)

	if exist {
		return value
	} else {
		return fallback
	}
}

func connectToDB(){
	cfg := APPConfig
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s ssl=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal("Failed to get database instance", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}