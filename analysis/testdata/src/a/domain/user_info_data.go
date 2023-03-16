package domain

import "time"

type UserInfoDataID ID // OK

func NewUserInfoDataTestID(username string, userNumber int) UserInfoDataID { // want "コンストラクタ名NewUserInfoDataTestIDがNewUserInfoDataIDではありません"
	return UserInfoDataID(username + "_" + string(userNumber))
}

type UserData struct { // want "ファイル名user_info_data.goとID名UserDataが一致していません"
	ID         UserInfoDataID
	Username   string
	UserNumber int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUserData(username string, userNumber int) UserData { // want "コンストラクタ名NewUserDataがNewUserInfoDataではありません"
	return UserData{
		ID:         NewUserInfoDataTestID(username, userNumber),
		Username:   username,
		UserNumber: userNumber,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
