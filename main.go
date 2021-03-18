package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

var (
	rooms  = make(map[string]map[*websocket.Conn]bool)
	Salt   = []byte("Keyblade")
	engine = &xorm.Engine{}
)

const DrvierName = "sqlite3"
const DataSourceName = "lotusdb.db"

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

// Validate validate
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func newEngine(dn string, dsn string) *xorm.Engine {
	engine, err := xorm.NewEngine(dn, dsn)
	if err != nil {
		log.Fatal("newEngine", err)
		return nil
	}

	engine.ShowSQL(true)
	return engine
}

func initEcho() *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = NewValidator()
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
	e.POST("/room/new", newRoomPost, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("kokorowotokihanate"),
	}))

	// e.POST("/new", newRoomPost)
	e.GET("/room/:id", enterRoom, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("kokorowotokihanate"),
	}))
	// Start server

	return e
}

func main() {
	go h.run()
	e := initEcho()
	engine = newEngine(DrvierName, DataSourceName)
	defer engine.Close()

	//Init Xorm engine
	err := engine.Sync2(new(Cards), new(RoomMembers), new(Rooms), new(Users))
	if err != nil {
		// log.Fatal("newEngine", err)
		fmt.Printf("Sync db struct failed.")
		return
	}

	if err != nil {
		fmt.Println("Open db failed")
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Lotus")
}

func login(c echo.Context) (err error) {

	L := new(LoginForm)

	if err = c.Bind(L); err != nil {
		return c.String(http.StatusInternalServerError, "Unknown Error")
	}
	// fmt.Printf("114514")
	log.Printf(L.Username)

	User := &Users{Username: L.Username}

	ok, err := engine.Get(User)

	if !ok {
		log.Fatalf(err.Error())
		return c.String(http.StatusUnauthorized, "Incorrect username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(L.Password)); err != nil {
		return c.String(http.StatusUnauthorized, "Incorrect username or password")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = L.Username
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
	log.Println(c.Get("user"))
	NewRoomInformation := new(RoomInformation)
	fmt.Println(c)
	if err := c.Bind(&NewRoomInformation); err != nil {
		log.Printf(NewRoomInformation.Tittle)
		fmt.Printf(err.Error())
		response := &JsonResponse{
			Code: 502,
			Msg:  "Invalid request 1",
			Data: struct{}{},
		}

		return c.JSONPretty(http.StatusBadRequest, response, " ")
	}

	log.Printf(NewRoomInformation.Tittle)

	Room := new(Rooms)

	Room.Tittle = NewRoomInformation.Tittle
	Room.AllowOnlookerWhenPlaying = NewRoomInformation.AllowOnlookerWhenPlaying
	Room.Description = NewRoomInformation.Description
	Room.Tag = NewRoomInformation.Tag
	Room.IsPasswordRequired = NewRoomInformation.IsPasswordRequired
	log.Println(c.Request().Header.Get(echo.HeaderAuthorization))
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	if len(authHeader) <= 7 {

		response := &JsonResponse{
			Code: 401,
			Msg:  "Invalid jwt token",
			Data: struct{}{},
		}
		return c.JSONPretty(http.StatusUnauthorized, response, "")
	}
	log.Print("authHeader" + strconv.Itoa(len(authHeader)))
	tokenStr := authHeader[7:]
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("kokorowotokihanate"), nil
	})

	log.Println(token)

	if err != nil {
		log.Println(err.Error())
		response := &JsonResponse{
			Code: 401,
			Msg:  "Invalid jwt token",
			Data: struct{}{},
		}
		return c.JSONPretty(http.StatusUnauthorized, response, "")
	}

	name := claims["name"].(string)

	User := &Users{Username: name}

	ok, err := engine.Get(User)

	if !ok {
		// log.Fatalf(err.Error())
		response := &JsonResponse{
			Code: 400,
			Msg:  "Invalid request: user not found",
			Data: struct{}{},
		}
		return c.JSONPretty(http.StatusBadRequest, response, "")
	}

	Room.HostId = User.Id
	Room.CreatedTime = time.Now().Unix()
	Room.LastActiveTime = Room.CreatedTime
	Room.Password = NewRoomInformation.Password

	affected, err := engine.Insert(Room)
	if err != nil {
		log.Fatalf(err.Error())
		response := &JsonResponse{
			Code: 500,
			Msg:  "Internal error",
			Data: struct{}{},
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "")
	}

	log.Println("affect=", affected)

	data := &RoomNo{
		RoomNo: Room.Id,
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
