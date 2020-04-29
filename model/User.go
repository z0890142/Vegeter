package model

type User struct {
	IslandsName string
	UserName    string
	TimeZone    int
	Uuid        string
}

type FriendShip struct {
	User1        string
	User2        string
	Relationship int //1=已發出邀請 2＝朋友
}

type Price struct {
	Uuid        string
	Price       int
	Description string
	isOverTime  bool
	Time        string
}
