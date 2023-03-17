// 完璧なテストケース
package domain

import "time"

type UserChatDataID ID // OK

type UserChatDataIDs []UserChatDataID

func NewUserChatDataID(username string, userNumber int) UserChatDataID { // OK
	return UserChatDataID(username + "_" + string(userNumber))
}

type UserChatData struct { // OK
	ID         UserChatDataID
	Username   string
	UserNumber int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUserChatData(username string, userNumber int) UserChatData { // OK
	return UserChatData{
		ID:         NewUserChatDataID(username, userNumber),
		Username:   username,
		UserNumber: userNumber,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
