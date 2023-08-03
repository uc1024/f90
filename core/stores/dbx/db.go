package dbx

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type (
	Server struct {
		Host string `json:"host"`
		Port int32  `json:"port"`
		User string `json:"user"`
		Pwd  string `json:"pwd"`
	}
	Options struct {
		Debug   bool     `json:"debug"`
		Name    string   `json:"name"`
		Idle    int      `json:"idle"`
		Open    int      `json:"open"`
		Source  []Server `json:"source"`
		Replica []Server `json:"replica"`
		Logger  logger.Interface
	}

	DBOptions = func(*Options)
)

func NewDB(opts ...DBOptions) (*gorm.DB, error) {

	options := Options{}

	for _, v := range opts {
		v(&options)
	}

	var open = func(tag string, name string, list []Server) []gorm.Dialector {
		dialectors := make([]gorm.Dialector, 0, len(list))
		for i, v := range list {
			var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", []interface{}{
				v.User,
				v.Pwd,
				v.Host,
				v.Port,
				name,
				"charset=utf8mb4&parseTime=True&loc=Local",
			}...)
			_ = i
			dialectors = append(dialectors, mysql.Open(dsn))
		}
		return dialectors
	}
	var data = dbresolver.Config{
		Sources:  open("source", options.Name, options.Source),
		Replicas: open("replica", options.Name, options.Replica),
	}

	db, err := gorm.Open(data.Sources[0], &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Use(dbresolver.Register(data))
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if options.Idle > 0 {
		sqlDB.SetMaxOpenConns(options.Idle)
	}
	if options.Open > 0 {
		sqlDB.SetMaxOpenConns(options.Open)
	}

	lv := logger.Error

	if options.Debug {
		lv = logger.Info
	}

	if options.Logger != nil {
		db.Logger = options.Logger.LogMode(lv)
	} else {
		db.Logger = logger.Default.LogMode(lv)
	}

	return db, nil
}
