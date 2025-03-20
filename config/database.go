package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB constant
var DB *gorm.DB

func LoadDB() {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", ENV.DbUsername, ENV.DbPassword, ENV.DbURL, ENV.DbPort, ENV.DbDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db
}

//// LoadDB function for load database connection
//func LoadDB() {
//	aes := encryption.NewAes()
//	dbPassword, err := aes.Decrypt(ENV.DbPassword)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	dbDebug := logger.Default.LogMode(logger.Info)
//	if !ENV.DbDebug {
//		dbDebug = nil
//	}
//
//	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", ENV.DbURL, ENV.DbUsername, dbPassword, ENV.DbDatabase, ENV.DbPort)
//	db, err := gorm.Open(postgres.New(postgres.Config{
//		DSN:                  connectionString,
//		PreferSimpleProtocol: true, // disables implicit prepared statement usage
//	}), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			TablePrefix:   fmt.Sprintf("%s.", ENV.DbSchema),
//			SingularTable: false,
//		},
//		Logger: dbDebug,
//	})
//
//	if err != nil {
//		panic("failed to connect database")
//	}
//
//	DB = db
//}
