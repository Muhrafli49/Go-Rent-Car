package config

import (
	"fmt"
	"log"
	"github.com/joho/godotenv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

var DB *sqlx.DB

func InitDB() {
	// Memuat file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Mengambil nilai dari variabel lingkungan
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Konfigurasi koneksi ke PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	// Membuka koneksi ke database
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Cek koneksi untuk memastikan berhasil
	if err = DB.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}

	fmt.Println("Database connected successfully!")
}
