package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	mockDB = map[string]string{}
)

func TestRoomPermission(t *testing.T) {

	tests := []struct {
		name       string
		user       User
		permission Permission
		room       string
		want       bool
	}{
		{"Test enter room with permission", User{"kyoku", map[Permission]map[string]bool{{"enter_room"}: {"114514": true}}, player}, Permission{"enter_room"}, "114514", true},
		{"Test enter room without permission", User{"heishin", map[Permission]map[string]bool{}, player}, Permission{"enter_room"}, "114514", false},
		{"Test enter room without corresponding permision", User{"kyoku", map[Permission]map[string]bool{{"enter_room"}: {"1919810": true}}, player}, Permission{"enter_room"}, "114514", false},
		{"Test delete room with identity of admin", User{"sis", map[Permission]map[string]bool{}, admin}, Permission{"delete_room"}, "114514", true},
		{"Test delete room with identity of super admin", User{"con", map[Permission]map[string]bool{}, superAdmin}, Permission{"delete_room"}, "114514", true},
		{"Test with identity of banned", User{"troll", map[Permission]map[string]bool{}, banned}, Permission{"delete_room"}, "114514", false},
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
	e := echo.New()
	t.Run("Test login with correct username and password", func(t *testing.T) {
		loginJSON := `{"username":"kyoku","password":"heishin"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

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
