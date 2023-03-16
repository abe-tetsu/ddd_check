package domain

import "time"

type UserInfoDataID ID // OK

func NewUserInfoDataID(username string, userNumber int) UserInfoDataID {
	return UserInfoDataID(username + "_" + string(userNumber))
}

type UserData struct {
	ID         UserInfoDataID
	Username   string
	UserNumber int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUserData(username string, userNumber int) UserData {
	return UserData{
		ID:         NewUserInfoDataID(username, userNumber),
		Username:   username,
		UserNumber: userNumber,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
