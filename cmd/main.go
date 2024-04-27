package main

import (
	"database/sql"
	"fmt"

	"encoding/json"
	"log"

	_ "github.com/lib/pq"

	"github.com/go-redis/redis"
)

const (
	username = "postgres"
	password = "Imangali2004"
	hostname = "localhost"
	port     = 5432
	db       = "postgres"
)

var (
	redisClient *redis.Client
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func initDB() *sql.DB {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, db)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func SetProductCache(product Product) error {

	productJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product to JSON: %v", err)
	}

	err = redisClient.HSet("products", product.ID, productJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to set product in cache: %v", err)
	}

	return nil
}

func GetProductFromCache(productID string) (*Product, error) {
	productJSON, err := redisClient.HGet("products", productID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get product from cache: %v", err)
	}

	var product Product
	err = json.Unmarshal([]byte(productJSON), &product)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product JSON: %v", err)
	}

	return &product, nil
}

func main() {
	initRedis()
	db := initDB()

	defer db.Close()
	product := Product{
		ID:    "1",
		Name:  "Example Product",
		Price: 100,
	}

	err := SetProductCache(product)
	if err != nil {
		log.Fatalf("Failed to set product in cache: %v", err)
	}

	retrievedProduct, err := GetProductFromCache("1")
	if err != nil {
		log.Fatalf("Failed to get product from cache: %v", err)
	}

	log.Printf("Retrieved product from cache: %+v\n", *retrievedProduct)
}
