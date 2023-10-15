package database

import (
	"fmt"
	"shadowguard/pkg/config"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	Insert(v interface{}) (int64, error)
}

type MockDatabase struct {
	elements map[string]interface{}
}

func (m *MockDatabase) Insert(v interface{}) (int64, error) {
	id := uuid.NewString()
	m.elements[id] = v
	return 0, nil
}

func NewMock() *MockDatabase {
	return &MockDatabase{
		elements: map[string]interface{}{},
	}
}

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
