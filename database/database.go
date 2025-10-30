package database

import (
	"blog_system/config"
	"blog_system/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error
	var dialector gorm.Dialector

	// 根据配置选择数据库驱动
	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.Database.GetMySQLDSN())
		log.Printf("Using MySQL database: %s", cfg.Database.DSN)
	case "sqlite":
		fallthrough
	default:
		dialector = sqlite.Open(cfg.Database.DSN)
		log.Printf("Using SQLite database: %s", cfg.Database.DSN)
	}

	// 配置GORM
	gormConfig := &gorm.Config{}
	if cfg.Server.Mode == "debug" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return err
	}

	// 自动迁移
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		return err
	}

	log.Println("Database connected and migrated successfully")
	return nil
}
