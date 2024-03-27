package datastore

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const (
	sslmode  = " sslmode="
	user     = "user="
	password = " password="
	dbname   = " dbname="
)

func CreateNewSearchRDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func CreateNewSelectionRDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	return rdb
}

func CreateNewDBForSurvey() (*sql.DB, error) {
	connStr := fmt.Sprintf("%s%s%s%s%s%s%s%s", user, viper.GetString("postgre.user"),
		password, viper.GetString("postgre.password"), dbname, viper.GetString("postgre.dbname"),
		sslmode, viper.GetString("postgre.sslmode"))

	dtb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error from `Open` function, package `sql`: %#v", err)
	}
	return dtb, nil
}

func CreateNewDBForVehicles() (*sql.DB, error) {
	connStr := fmt.Sprintf("%s%s%s%s%s%s%s%s", user, viper.GetString("postgre.user"),
		password, viper.GetString("postgre.password"), dbname, viper.GetString("postgre.dbname1"),
		sslmode, viper.GetString("postgre.sslmode"))

	dtb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error from `Open` function, package `sql`: %#v", err)
	}
	return dtb, nil
}
