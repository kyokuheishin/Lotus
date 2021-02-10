package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/login", login)

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
