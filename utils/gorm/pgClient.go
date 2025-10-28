package mygorm

import (
	"asset-pulse-api/entities"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresCredential struct {
	Host        string
	Port        string
	DBName      string
	Username    string
	Password    string
	SSLMode     string
	SSLRootCert string
	TimeZone    string
}

func (p *PostgresCredential) Parse() string {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		p.Host, p.Port, p.Username, p.Password, p.DBName, p.SSLMode, p.TimeZone,
	)

	if p.SSLRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", p.SSLRootCert)
	}

	return dsn
}

func NewPostgres(dsn string, config *gorm.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err = db.AutoMigrate(&entities.SoftwareLicense{}, &entities.RawCurrentGroupedSoftware{}); err != nil {
		fmt.Sprintf("Failed to migrate database: %v", err)
	}

	return db
}
