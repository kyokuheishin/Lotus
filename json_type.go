package main

type LoginForm struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type RoomInformation struct {
	IsPasswordRequired bool   `json:"is_password_required" form:"is_password_required" query:"is_password_required" validate:"required"`
	Tittle             string `json:"tittle" form:"tittle" query:"tittle" validate:"required"`
	Description        string `json:"description" form:"description" query:"description" validate:"required"`
	Tag                string `json:"tag" form:"tag" query:"tag" validate:"required"`
	// CreatedTime              int64  `json:"created_time" form:"created_time" query:"created_time" validate:"required"`
	Password string `json:"password"`
	// IsPlaying                bool   `json:"is_playing"`
	AllowOnlookerWhenPlaying bool `json:"allow_onlooker_when_playing" validate:"required"`
}

type CreateRoomForm struct {
	Tittle                   string `json:"tittle" form:"tittle" query:"tittle" validate:"required"`
	Description              string `json:"description" form:"description" query:"description" validate:"required"`
	Tag                      string `json:"tag" form:"tag" query:"tag" validate:"required"`
	CreatedTime              int64  `json:"created_time" form:"created_time" query:"created_time" validate:"required"`
	Password                 string `json:"password"`
	AllowOnlookerWhenPlaying bool   `json:"allow_onlooker_when_playing" validate:"required"`
}

type RoomMember struct {
	RoomID int    `json:"room_id"`
	Role   string `json:"role"`
	UserID int    `json:"user_id"`
	CardID int    `json:"card_id"`
}

type User struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	Role           string `json:"role"`
	CreatedTime    int64  `json:"created_time"`
	LastActiveTime int64  `json:"last_active_time"`
}

type Card struct {
	CardID      int    `json:"card_id"`
	UserID      int    `json:"user_id"`
	Content     string `json:"content"`
	CreatedTime int64  `json:"created_time"`
}

type RoomNo struct {
	RoomNo int `json:"room_no"`
}

type JsonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type WebSocketMessage struct {
	JwtToken string               `json:"jwt_token"`
	Data     WebSocketMessageData `json:"data"`
}

type WebSocketMessageData struct {
	Text string `json:"text"`
	Cmd  string `json:"cmd"`
}

type JwtTokenData struct {
	Token string `json:"token"`
}
