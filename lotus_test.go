package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	mockDB   = map[string]string{}
	Jwtkey   string
	fixtures *testfixtures.Loader
	e        = initEcho()
)

func TestRoomPermission(t *testing.T) {

	tests := []struct {
		name       string
		user       UserLegacy
		permission Permission
		room       string
		want       bool
	}{
		{"Test enter room with permission", UserLegacy{"kyoku", map[Permission]map[string]bool{{"enter_room"}: {"114514": true}}, player}, Permission{"enter_room"}, "114514", true},
		{"Test enter room without permission", UserLegacy{"heishin", map[Permission]map[string]bool{}, player}, Permission{"enter_room"}, "114514", false},
		{"Test enter room without corresponding permision", UserLegacy{"kyoku", map[Permission]map[string]bool{{"enter_room"}: {"1919810": true}}, player}, Permission{"enter_room"}, "114514", false},
		{"Test delete room with identity of admin", UserLegacy{"sis", map[Permission]map[string]bool{}, admin}, Permission{"delete_room"}, "114514", true},
		{"Test delete room with identity of super admin", UserLegacy{"con", map[Permission]map[string]bool{}, superAdmin}, Permission{"delete_room"}, "114514", true},
		{"Test with identity of banned", UserLegacy{"troll", map[Permission]map[string]bool{}, banned}, Permission{"delete_room"}, "114514", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := checkRoomPermission(test.user, test.permission, test.room)

			if res != test.want {
				t.Errorf("Expected %t but got %t", test.want, res)
			}
		})
	}

}

func TestLogin(t *testing.T) {
	// e := echo.New()
	t.Run("Test login with correct username and password", func(t *testing.T) {
		loginJSON := `{"username":"kyoku","password":"heishin"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var jtd JwtTokenData
			if err := json.Unmarshal(rec.Body.Bytes(), &jtd); err != nil {
				fmt.Println(err.Error())
			}

			Jwtkey = "Bearer " + jtd.Token

		}
	})

	t.Run("Test login with incorrect username and password", func(t *testing.T) {
		loginJSON := `{"username":"kyoku","password":"baka"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, login(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)

		}
	})

}

func TestNewRoom(t *testing.T) {

	fmt.Println(Jwtkey)
	// _ := RoomInformation{
	// 	IsPasswordRequired:       false,
	// 	Tittle:                   "Another test",
	// 	Description:              "A test",
	// 	Tag:                      "Test",
	// 	Password:                 "",
	// 	AllowOnlookerWhenPlaying: true,
	// }
	Payload := `{"is_password_required":false,"tittle":"Another test","description":"A test","tag":"Test","password":"","allow_onlooker_when_playing":true}`
	fmt.Printf(string(Payload))
	t.Run("Test create room with JWT", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/room/new", strings.NewReader(Payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, Jwtkey)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		// e.ServeHTTP(rec, req)
		if assert.NoError(t, newRoomPost(c)) {
			// fmt.Printf(rec.Body.String())
			assert.Equal(t, http.StatusOK, rec.Code)
		}

	})

	t.Run("Test create room without JWT", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/room/new", strings.NewReader(""))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		e.ServeHTTP(rec, req)
		if assert.NoError(t, newRoomPost(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Test create room with invalid JWT", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/room/new", strings.NewReader(""))
		req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTMzOTU5ODYsIm5hbWUiOeW9rdSJ9.sUxRgnXqK1dgc-34IWjvHycGoTuU2IGF2vzdml2s8wg")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		e.ServeHTTP(rec, req)
		if assert.NoError(t, newRoomPost(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}

	})
}

// func TestEnterRoom(t *testing.T) {
// 	e := initEcho()

// 	t.Run("Test enter room that exists", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodGet, "/room/1", nil)
// 		req.Header.Set(echo.HeaderAuthorization, jwtkey)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		e.ServeHTTP(rec, req)
// 		if assert.NoError(t, enterRoom(c)) {
// 			assert.Equal(t, http.StatusOK, rec.Code)
// 		}
// 	})

// 	t.Run("Test enter room that exists without jwt", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodGet, "/room/1", nil)

// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		e.ServeHTTP(rec, req)
// 		if assert.NoError(t, enterRoom(c)) {
// 			assert.Equal(t, http.StatusBadRequest, rec.Code)
// 		}
// 	})

// 	t.Run("Test enter room that exists with invalid jwt", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodGet, "/room/1", nil)
// 		req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTMzOTU5ODYsIm5hbWUiOeW9rdSJ9.sUxRgnXqK1dgc-34IWjvHycGoTuU2IGF2vzdml2s8wg")
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		e.ServeHTTP(rec, req)
// 		if assert.NoError(t, enterRoom(c)) {
// 			assert.Equal(t, http.StatusUnauthorized, rec.Code)
// 		}
// 	})

// 	t.Run("Test enter room that not exists", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodGet, "/room/2", nil)
// 		req.Header.Set(echo.HeaderAuthorization, jwtkey)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		e.ServeHTTP(rec, req)
// 		if assert.NoError(t, enterRoom(c)) {
// 			assert.Equal(t, http.StatusNotFound, rec.Code)
// 		}
// 	})
// }

func TestMain(m *testing.M) {
	var err error

	engine = newEngine("sqlite3", "lotusdb_test.db")
	fixtures, err = testfixtures.New(
		testfixtures.Database(engine.DB().DB),
		testfixtures.Dialect("sqlite"),
		testfixtures.Directory("fixtures/"),
	)

	if err != nil {
		fmt.Println("Open db failed")
	}

	if err = fixtures.Load(); err != nil {
		fmt.Printf(err.Error())
	}

	os.Exit(m.Run())
}
