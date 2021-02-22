package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	rooms = make(map[string]map[*websocket.Conn]bool)
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return echo.NewHTTPError(http.StatusInternalServerError, cv.validator.Struct(i).Error())
}

func initEcho() *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}
	// Routes
	e.GET("/", hello)
	e.POST("/login", login)
	e.Group("/room", middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("kokorowotokihanate"),
	}))
	// e.GET("/room/new", newRoom, middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: []byte("kokorowotokihanate"),
	// }))

	// e.POST("/room/new", newRoomPost, middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: []byte("kokorowotokihanate"),
	// }))
	e.POST("/new", newRoomPost)
	e.GET("/room/:id", enterRoom, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("kokorowotokihanate"),
	}))
	// Start server

	return e
}

func main() {
	go h.run()
	e := initEcho()
	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Lotus")
}

func login(c echo.Context) (err error) {

	l := new(LoginForm)

	if err = c.Bind(l); err != nil {
		return c.String(http.StatusInternalServerError, "Unknown Error")
	}
	// fmt.Printf("114514")
	// fmt.Printf(l.username)
	if l.Username != "kyoku" || l.Password != "heishin" {
		return c.String(http.StatusUnauthorized, "Incorrect username or password")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "kyoku"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("kokorowotokihanate"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})

}

func newRoomPost(c echo.Context) error {
	newRoomInformation := new(RoomInformationForm)

	if err := c.Bind(newRoomInformation); err != nil {

		response := &JsonResponse{
			Code: 502,
			Msg:  "Invalid request",
			Data: struct{}{},
		}

		return c.JSONPretty(http.StatusBadRequest, response, " ")
	}

	if err := c.Validate(newRoomInformation); err != nil {
		response := &JsonResponse{
			Code: 502,
			Msg:  "Invalid request",
			Data: struct{}{},
		}
		return c.JSONPretty(http.StatusBadRequest, response, " ")
	}

	data := &RoomNo{
		RoomNo: 0423,
	}
	response := &JsonResponse{
		Code: 201,
		Msg:  "Room created",
		Data: data,
	}

	return c.JSONPretty(http.StatusOK, response, "")
}

func enterRoom(c echo.Context) error {
	roomID := c.Param("id")
	if roomID == "1" {
		return c.String(http.StatusOK, "OK")

	}

	return c.String(http.StatusNotFound, "404")
}
