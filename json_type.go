package main

type LoginForm struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type RoomInformationForm struct {
	IsPasswordRequired bool   `json:"is_password_required" form:"is_password_required" query:"is_password_required" validate:"required"`
	Tittle             string `json:"tittle" form:"tittle" query:"tittle" validate:"required"`
	Description        string `json:"description" form:"description" query:"description" validate:"required"`
	Tag                string `json:"tag" form:"tag" query:"tag" validate:"required"`
	CreatedTime        int    `json:"created_time" form:"created_time" query:"created_time" validate:"required"`
}

type RoomNo struct {
	RoomNo int `json:"room_no"`
}

type JsonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
