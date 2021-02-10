package main

type Permission struct {
	name string
}

func checkRoomPermission(u User, permission Permission, room string) bool {

	if u.identity == superAdmin || u.identity == admin {
		return true
	}

	if u.identity == banned {
		return false
	}

	if val, ok := u.permissions[permission]; ok {
		if val, ok := val[room]; ok {
			return val
		}
	}
	return false
}
