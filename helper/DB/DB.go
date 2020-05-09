package DB

import (
	"Vegeter/model"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-sql-driver/mysql" //前面加 _ 是為了只讓他執行init

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var db *sqlx.DB

func SetTls(Log *logrus.Entry) error {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("./config/db/client-cert.pem")
	if err != nil {
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("ioutil.ReadFile error")
		return err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}
	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair("./config/db/client-cert.pem", "./config/db/client-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	clientCert = append(clientCert, certs)
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
	})

	return nil
}

func CreateDbConn(driveName string, dataSourceName string, Log *logrus.Entry) error {
	var err error
	// SetTls(Log)
	db, err = sqlx.Open(driveName, dataSourceName)
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("CreateDbConn error")
		return err
	}

	Log.Info("connnect success")

	return err
}
func CheckUserEsixt(user model.User) (bool, string) {
	var tmpCount int
	var uuid string

	sqlString := "select count(UserName) as isExist,UUID from User where IslandsName=? and UserName=? group by UUID"
	db.QueryRow(sqlString, user.IslandsName, user.UserName).Scan(&tmpCount, &uuid)
	if tmpCount > 0 {
		return true, uuid
	}
	return false, uuid
}
func InsertUser(user model.User) error {
	sqlString := "insert into User(IslandsName,UserName,TimeZone,UUID,GameID) values(?,?,?,?,?)"

	_, err := db.Exec(sqlString, user.IslandsName, user.UserName, user.TimeZone, user.Uuid, user.GameID)
	if err != nil {
		return fmt.Errorf("Insert User : %v", err)
	}
	return err
}

func InsertFriendShip(friendShip model.FriendShip) error {
	sqlString := "insert into FriendShip(User1,User2,Relationship) values(?,?,1)"

	_, err := db.Exec(sqlString, friendShip.User1, friendShip.User2)
	if err != nil {
		return fmt.Errorf("Insert FirendShip : %v", err)
	}
	return err
}

func UpdateFriendShip(friendShip model.FriendShip) error {
	sqlString := "update FriendShip set Relationship=2 where User1=? and User2=?"

	_, err := db.Exec(sqlString, friendShip.User1, friendShip.User2)
	if err != nil {
		return fmt.Errorf("update FirendShip : %v", err)
	}
	return err
}

func InsertGameFriend(gameFriend model.GameFriend) error {
	sqlString := "insert into GameFriend(Applicant,ToUUID) values(?,?)"

	_, err := db.Exec(sqlString, gameFriend.Applicant, gameFriend.ToUUID)
	if err != nil {
		return fmt.Errorf("Insert GameFriend : %v", err)
	}
	return err
}

func InsertPrice(price model.Price) error {
	sqlString := "insert into Price(UUID,Price,Date,isOverTime) values(?,?,?,0)"

	_, err := db.Exec(sqlString, price.Uuid, price.Price, price.Date)
	if err != nil {
		return fmt.Errorf("Insert Price : %v", err)
	}
	return err
}
func UpdatePrice(price model.Price) error {
	sqlString := "update Price set Price=? where UUID=? and isOverTime=0"

	_, err := db.Exec(sqlString, price.Price, price.Uuid)
	if err != nil {
		return fmt.Errorf("Update Price : %v", err)
	}
	return err
}

func GetCurrentPrice(uuid string) (price int, err error) {
	sqlString := "select Price from Price where UUID=? and isOverTime=0 order by Price desc"
	err = db.QueryRow(sqlString, uuid).Scan(&price)
	if err != nil {
		err = fmt.Errorf("Get CurrentPrice : %v", err)
		return
	}
	return
}

func GetPriceRecord(uuid string) (priceList []model.PriceList, err error) {
	sqlString := "select P.Price,P.Date,u.IslandsName,u.UserName,P.isOverTime " +
		"from Price as P join User as u on P.UUID=u.UUID " +
		"where (P.UUID in (select User1 from FriendShip where User2=? and Relationship=2)) or " +
		"(P.UUID in (select User2 from FriendShip where User1=? and Relationship=2)) or P.UUID=? order by Price desc"
	rows, err := db.Query(sqlString, uuid, uuid, uuid)
	if err != nil {
		err = fmt.Errorf("Get CurrentPrice : %v", err)
		return
	}
	for rows.Next() {
		var tmpRows model.PriceList
		rows.Scan(&tmpRows.Price, &tmpRows.Date, &tmpRows.IslandsName, &tmpRows.UserName, &tmpRows.IsOverTime)
		priceList = append(priceList, tmpRows)
	}
	return
}

func GetUserInfo(uuid string) (user model.User) {
	sqlString := "select IslandsName,UserName,TimeZone,UUID,GameID from User where UUID=?"
	db.QueryRow(sqlString, uuid).Scan(&user.IslandsName, &user.UserName, &user.TimeZone, &user.Uuid, &user.GameID)
	return
}

func GetCheckFriend(uuid string) ([]model.FirendList, error) {
	var friendShipList []model.FirendList
	sqlString := "select @no:=@no+1,UserName,IslandsName,UUID " +
		"from (select @no:=0) as n join FriendShip as f join User as u on u.UUID=f.User1 where User2=? and Relationship=1"
	rows, err := db.Query(sqlString, uuid)
	if err != nil {
		return friendShipList, fmt.Errorf("GetCheckFriend : %v", err)
	}

	for rows.Next() {
		var tmpFriendShip model.FirendList
		rows.Scan(&tmpFriendShip.No, &tmpFriendShip.UserName, &tmpFriendShip.IslandsName, &tmpFriendShip.Uuid)
		friendShipList = append(friendShipList, tmpFriendShip)
	}
	return friendShipList, err
}

func GetFriend(uuid string) ([]model.FirendList, error) {
	var friendShipList []model.FirendList
	sqlString1 := "select @no:=@no+1,UserName,IslandsName,UUID,GameID " +
		"from (select @no:=0) as n join FriendShip as f join User as u on u.UUID=f.User1 where User2=? and Relationship=2"
	sqlString2 := "select @no:=@no+1,UserName,IslandsName,UUID,GameID " +
		"from (select @no:=0) as n join FriendShip as f join User as u on u.UUID=f.User2 where User1=? and Relationship=2"
	rows, err := db.Query(sqlString1, uuid)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return friendShipList, fmt.Errorf("GetFriend : %v", err)
	}
	for rows.Next() {
		var tmpFriendShip model.FirendList
		rows.Scan(&tmpFriendShip.No, &tmpFriendShip.UserName, &tmpFriendShip.IslandsName, &tmpFriendShip.Uuid, &tmpFriendShip.GameID)
		friendShipList = append(friendShipList, tmpFriendShip)
	}
	rows, err = db.Query(sqlString2, uuid)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return friendShipList, fmt.Errorf("GetFriend : %v", err)
	}
	for rows.Next() {
		var tmpFriendShip model.FirendList
		rows.Scan(&tmpFriendShip.No, &tmpFriendShip.UserName, &tmpFriendShip.IslandsName, &tmpFriendShip.Uuid, &tmpFriendShip.GameID)
		friendShipList = append(friendShipList, tmpFriendShip)
	}

	return friendShipList, err
}

func GetAllPriceRecord() (priceList []model.PriceList, err error) {
	sqlString := "select P.Price,P.Date,u.IslandsName,u.UserName,P.isOverTime,P.UUID " +
		"from Price as P join User as u on P.UUID=u.UUID " +
		"where P.isOverTime=0 order by Price desc"
	rows, err := db.Query(sqlString)
	if err != nil {
		err = fmt.Errorf("Get AllPriceRecord : %v", err)
		return
	}
	for rows.Next() {
		var tmpRows model.PriceList
		rows.Scan(&tmpRows.Price, &tmpRows.Date, &tmpRows.IslandsName, &tmpRows.UserName, &tmpRows.IsOverTime, &tmpRows.Uuid)
		priceList = append(priceList, tmpRows)
	}
	return
}
