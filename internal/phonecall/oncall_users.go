package phonecall

import (
	"fmt"
	"time"

	"github.com/weiqiang333/infra-prometheus-webhook/model"
)

var user model.OncallUser

// GetOncallUser 连接数据库,获取当前值班人员
func GetOncallUser(role string) model.OncallUser {
	db := GetDB()
	defer db.Close()
	table := model.Config.DataBase["table"]
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
			fmt.Print(err)
		}
	}
	return user
}