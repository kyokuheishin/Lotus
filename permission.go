package main

type Permission struct {
	name string
}

func checkCanManageRoomMember(u User, rm RoomMember, roomID int, db interface{}) (message string, ok bool) {
	if !(rm.Role == "creator" || rm.Role == "kp") {
		return "Permission denied", false
	}

	return "OK", true
}

func checkCanSendMessage(u User, rm RoomMember, roomID int, db interface{}) (message string, ok bool) {
	if rm.Role == "muted" {
		return "The user is muted in this room", false
	}
	return "OK", true
}

// func checkEnterRoomPermission(u User, rDb RoomInformation, rForm RoomInformation, db interface{}) (message string, ok bool) {
// 	if u.Role == "banned" {
// 		return "User is banned", false
// 	}

// 	if rDb.IsPlaying == true && rDb.AllowOnlookerWhenPlaying == false {
// 		return "A onlooker is not allowed in this room", false
// 	}

// 	if rDb.IsPasswordRequired == true {
// 		if err := bcrypt.CompareHashAndPassword([]byte(rDb.Password), []byte(rForm.Password)); err != nil {
// 			return "Wrong password", false
// 		}
// 	}

// 	return "Ok", true

// }

func checkRoomPermission(u UserLegacy, permission Permission, room string) bool {

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
