package phonecall

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type OncallUser struct {
	Mobile   string
	UserName string
}

// GetOncallUser 连接数据库,获取当前值班人员
func GetOncallUser(role string) OncallUser {
	var user OncallUser
	db := GetDB()
	defer db.Close()
	table := viper.GetString("DataBase.table")
	date := time.Now().UTC().Format("2006-01-02")
	sql := `
		SELECT username, mobile
		FROM ` + table +
		` WHERE role=$1 AND started_at<=$2 
		ORDER BY started_at DESC LIMIT 1
		`
	rows, err := db.Query(sql, role, date)
	if err != nil {
		fmt.Print(err)
	}
	// 判断是否查询到了数据
	// 如果没有查询到数据，就直接进入同比数据查询
	if rows.Next() {
		err := rows.Scan(&user.UserName, &user.Mobile)
		if err != nil {
			log.Println("Failed GetOncallUser error: ", err.Error())
		}
	}
	return user
}
