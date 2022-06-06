package phonecall

import (
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/weiqiang333/infra-prometheus-webhook/model"
)

func GetDB() *sql.DB {
	db := model.Config.DataBase
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		db["host"], db["port"], db["username"], db["password"], db["database"])
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
