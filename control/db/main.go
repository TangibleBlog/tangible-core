package db

import (
	"tangible-core/public/common"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Dbc *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("./documents/history.db"), &gorm.Config{})
	common.HandleError(err)
	err = db.AutoMigrate(&History{})
	common.HandleError(err)
	Dbc = db
}
