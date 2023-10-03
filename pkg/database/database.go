package database

import (
	"fmt"
	"shadowguard/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	conn *gorm.DB
}

func (db *Database) Insert(v interface{}) (int64, error) {
	transaction := db.conn.Create(v)
	return transaction.RowsAffected, transaction.Error
}

func New(dbConfig config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
	)

	dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	dbConnection.AutoMigrate(&Request{})

	return &Database{
		conn: dbConnection,
	}, nil
}
