package domain // want "IDと構造体のコンストラクタの定義がされていません"

import "time"

type UserFriendListID ID // OK

type UserFriendList struct {
	ID         UserFriendListID
	Username   string
	UserNumber int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
