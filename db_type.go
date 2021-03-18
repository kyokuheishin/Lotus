package main

type RoomMembers struct {
	Id     int    `xorm:"not null pk autoincr INTEGER"`
	RoomId int    `xorm:"not null INTEGER"`
	Role   string `xorm:"not null NULL VARCHAR"`
	UserId int    `xorm:"not null INTEGER"`
	CardId int    `xorm:"not null INTEGER"`
}

type Cards struct {
	Id          int    `xorm:"not null pk autoincr INTEGER"`
	UserId      int    `xorm:"not null INTEGER"`
	Content     []byte `xorm:"BLOB"`
	CreatedTime int64  `xorm:"not null BIGINT"`
}

type Rooms struct {
	Id                       int    `xorm:"not null pk autoincr INTEGER"`
	HostId                   int    `xorm:"not null INTEGER"`
	IsPasswordRequired       bool   `xorm:"not null default '(0)' BOOLEAN"`
	Tittle                   string `xorm:"VARCHAR"`
	Description              string `xorm:"VARCHAR"`
	Tag                      string `xorm:"VARCHAR"`
	CreatedTime              int64  `xorm:"not null BIGINT"`
	LastActiveTime           int64  `xorm:"not null BIGINT"`
	Password                 string `xorm:"STRING"`
	AllowOnlookerWhenPlaying bool   `xorm:"BOOLEAN"`
}

type Users struct {
	Id       int    `xorm:"not null pk autoincr INTEGER"`
	Username string `xorm:"not null NULL VARCHAR"`

	Password string `xorm:"not null NULL VARCHAR"`

	CreatedTime    int64  `xorm:"not null BIGINT"`
	LastActiveTime int64  `xorm:"not null BIGINT"`
	Role           string `xorm:"not null NULL VARCHAR"`
}
