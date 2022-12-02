package phonecall

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func GetDB() *sql.DB {
	dbHost := viper.GetString("DataBase.host")
	dbPort := viper.GetString("DataBase.port")
	dbUsername := viper.GetString("DataBase.username")
	dbPassword := viper.GetString("DataBase.password")
	dbDatabase := viper.GetString("DataBase.database")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbHost, dbPort, dbUsername, dbPassword, dbDatabase)
	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	err = database.Ping()
	if err != nil {
		fmt.Print(err)
		return nil
	}

	return database
}
