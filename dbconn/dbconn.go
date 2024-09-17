package dbconn

import (
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/apache/calcite-avatica-go/v5"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// singleton instance of database connection.
var dbInstance *gorm.DB
var dbOnce sync.Once

// DB creates a new instance of gorm.DB if a connection is not established.
// return singleton instance.
func DB() *gorm.DB {
	if dbInstance == nil {
		dbOnce.Do(openDB)
	}
	return dbInstance
}

// openDB initialize gorm DB.
func openDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	gormDB, err := gorm.Open(
		mysql.Open(os.Getenv("DB_CONNECTION_URL")),
		&gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 newLogger,
			QueryFields:            true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)

	if err != nil {
		panic("dbconn openDB PortalBRIBrain: cannot open database")
	}
	dbInstance = gormDB
	db, err := dbInstance.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(200)
}
