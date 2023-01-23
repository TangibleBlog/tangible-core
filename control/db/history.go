package db

type History struct {
	UUID       string `gorm:"primary key;"`
	MetaData   string
	Content    string
	CreateTime int64
}
