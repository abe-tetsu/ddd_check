package domain

import "time"

type ChannelID ID // want "ファイル名channel_data.goとID名ChannelIDが一致していません"

func NewChannelDataID(channelID string) ChannelID { // OK
	return ChannelID(channelID)
}

type ChannelData struct {
	ID        ChannelID
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewChannel(channelID string) ChannelData { // want "コンストラクタ名NewChannelにファイル名ChannelDataが含まれていません"
	return ChannelData{
		ID:        NewChannelDataID(channelID),
		ChannelID: channelID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func Newchanneldataaa(channelID string) ChannelData { // TODO: 含まれているかどうか・大文字小文字区別なしで確認しているから、これも通る
	return ChannelData{
		ID:        NewChannelDataID(channelID),
		ChannelID: channelID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewChannelData(channelID string) ChannelData { // OK
	return ChannelData{
		ID:        NewChannelDataID(channelID),
		ChannelID: channelID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewChanneltesttestData(channelID string) ChannelData { // want "コンストラクタ名NewChanneltesttestDataにファイル名ChannelDataが含まれていません"
	return ChannelData{
		ID:        NewChannelDataID(channelID),
		ChannelID: channelID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
