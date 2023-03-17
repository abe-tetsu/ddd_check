package domain

import "time"

type UserInfoDataID ID // OK

func NewUserInfoDataTestID(username string, userNumber int) UserInfoDataID { // want "ファイル名UserInfoDataとコンストラクタ名NewUserInfoDataTestIDが一致していません"
	return UserInfoDataID(username + "_" + string(userNumber))
}

type UserData struct { // want "ファイル名user_info_data.goと構造体名UserDataが一致していません"
	ID         UserInfoDataID
	Username   string
	UserNumber int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUserData(username string, userNumber int) UserData { // want "コンストラクタ名NewUserDataにファイル名UserInfoDataが含まれていません"
	return UserData{
		ID:         NewUserInfoDataTestID(username, userNumber),
		Username:   username,
		UserNumber: userNumber,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func NewUserDataFromTest(username string, userNumber int) UserData { // want "コンストラクタ名NewUserDataFromTestにファイル名UserInfoDataが含まれていません"
	return UserData{
		ID:         NewUserInfoDataTestID(username, userNumber),
		Username:   username,
		UserNumber: userNumber,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func NewUserInfoData(username string, userNumber int) UserData { // OK
	return UserData{
		ID:         NewUserInfoDataTestID(username, userNumber),
		Username:   username,
		UserNumber: userNumber,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func test() {
	NewUserInfoData("test", 1)
}
