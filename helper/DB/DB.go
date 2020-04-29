package DB

import (
	"Vegeter/model"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var db *sqlx.DB

func CreateDbConn(driveName string, dataSourceName string, Log *logrus.Entry) error {
	var err error
	db, err = sqlx.Open(driveName, dataSourceName)
	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil {
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("CreateDbConn error")
		return err
	}

	Log.Info("connnect success")

	return err
}

func InsertUser(user model.User) error {
	sqlString := "insert into User('IslandsName','UserName','TimeZone','UUID') values(?,?,?,?)"

	_, err := db.Exec(sqlString, user.IslandsName, user.UserName, user.TimeZone, user.Uuid)
	if err != nil {
		return fmt.Errorf("Insert User : %v", err)
	}
	return err
}

func InsertFirendShip(friendShip model.FriendShip) error {
	sqlString := "insert into User('User1','User2','Relationship') values(?,?,1)"

	_, err := db.Exec(sqlString, friendShip.User1, friendShip.User2, friendShip.Relationship)
	if err != nil {
		return fmt.Errorf("Insert FirendShip : %v", err)
	}
	return err
}

func InsertPrice(price model.Price) error {
	sqlString := "insert into Price('UUID','Price','Description','Time','isOverTime') values(?,?,?,?,1)"

	_, err := db.Exec(sqlString, price.Uuid, price.Price, price.Description, price.Time)
	if err != nil {
		return fmt.Errorf("Insert Price : %v", err)
	}
	return err
}
func UpdatePrice(price model.Price) error {
	sqlString := "update Price set Price=? and Description=? where UUID=? and isOverTime=0"

	_, err := db.Exec(sqlString, price.Price, price.Description, price.Uuid)
	if err != nil {
		return fmt.Errorf("Update Price : %v", err)
	}
	return err
}

func GetCurrentPrice(uuid string) (price model.Price, err error) {
	sqlString := "select UUID,Price,Description,Time from Price where UUID=? and isOverTime=0"
	err = db.QueryRow(sqlString, uuid).Scan(&price)
	if err != nil {
		err = fmt.Errorf("Get CurrentPrice : %v", err)
		return
	}
	return
}

func GetPriceRecord(uuid string) (priceList []model.Price, err error) {
	sqlString := "select Price,Description,Time from Price where UUID=?"
	rows, err := db.Query(sqlString, uuid)
	if err != nil {
		err = fmt.Errorf("Get CurrentPrice : %v", err)
		return
	}
	for rows.Next() {
		var tmpRows model.Price
		rows.Scan(&tmpRows)
		priceList = append(priceList, tmpRows)
	}
	return
}
