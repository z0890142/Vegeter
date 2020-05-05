package model

type User struct {
	IslandsName string
	UserName    string
	TimeZone    string
	Uuid        string
	GameID      string
}

type GameFriend struct {
	ToUUID       string
	ToUUIDGameID string
	Applicant    string
	Approve      string
}

type FriendShip struct {
	User1        string
	User2        string
	Relationship int //1=已發出邀請 2＝朋友
}

type FirendList struct {
	No int
	User
}

type GmaeFirendList struct {
	No int
	User
	ToUUIDGameID string
}

type Price struct {
	Uuid       string
	Price      int
	isOverTime bool
	Date       string
}

type PriceList struct {
	User
	Price      int
	IsOverTime bool
	Date       string
}
