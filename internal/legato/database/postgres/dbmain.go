package postgres

import (
	"fmt"
	"legato_server/config"
	"legato_server/internal/legato/database"
	"legato_server/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LegatoDB struct {
	db *gorm.DB
}

var legatoDb LegatoDB

var log, _ = logger.NewLogger(logger.Config{})

func NewPostgresDatabase(cfg *config.Config) (database.Database, error) {
	databaseConf := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.DatabaseName,
		cfg.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(databaseConf), &gorm.Config{})
	if err != nil {
		log.Warn("Error in connecting to the postgres database")
		log.Fatal(err)
	}

	legatoDb.db = db

	log.Info("Creating schema...")
	err = createSchema(legatoDb.db)
	if err != nil {
		return nil, err
	}

	return &legatoDb, nil
}

// createSchema creates database schema (tables and ...)
// for all of our models.
func createSchema(db *gorm.DB) error {
	_ = db.AutoMigrate(User{})
	_ = db.AutoMigrate(Connection{})
	_ = db.AutoMigrate(Scenario{})
	_ = db.AutoMigrate(Service{})
	_ = db.AutoMigrate(Webhook{})
	_ = db.AutoMigrate(Http{})
	_ = db.AutoMigrate(Telegram{})
	_ = db.AutoMigrate(Spotify{})
	_ = db.AutoMigrate(Token{})
	_ = db.AutoMigrate(Ssh{})
	_ = db.AutoMigrate(History{})
	_ = db.AutoMigrate(ServiceLog{})
	_ = db.AutoMigrate(LogMessage{})
	_ = db.AutoMigrate(Gmail{})
	_ = db.AutoMigrate(Github{})
	_ = db.AutoMigrate(Discord{})
	_ = db.AutoMigrate(ToolBox{})

	return nil
}
