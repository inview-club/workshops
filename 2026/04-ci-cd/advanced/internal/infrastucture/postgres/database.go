package postgres

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"products/internal/entities"
	"products/internal/infrastucture/postgres/model"
)

type dbClient struct {
	client *gorm.DB
}

type Config struct {
	IP       string
	Port     string
	User     string
	Password string
	Database string
}

func New(cnfg *Config) (*dbClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s  password=%s  dbname=%s  port=%s  sslmode=disable", cnfg.IP, cnfg.User, cnfg.Password, cnfg.Database, cnfg.Port)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	return &dbClient{
		client: conn,
	}, nil
}

func (db *dbClient) QueryAll() ([]entities.Product, error) {
	var products []entities.Product
	err := db.client.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (db *dbClient) QueryByID(id string) (entities.Product, error) {
	var product entities.Product
	err := db.client.Find(&product, "id = ?", id).Error
	if err != nil {
		return entities.Product{ID: -1}, err
	}
	return product, nil
}

func (db *dbClient) Insert(product entities.Product) (string, error) {
	nProduct := model.Product{Name: product.Name, Price: product.Price, Status: int(product.Status)}
	result := db.client.Create(&nProduct)
	return strconv.Itoa(nProduct.ID), result.Error
}

func (db *dbClient) Ping() error {
	sqlDB, err := db.client.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
