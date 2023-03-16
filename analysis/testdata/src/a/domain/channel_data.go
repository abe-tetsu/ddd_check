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

func NewChannel(channelID string) ChannelData { // want "コンストラクタ名NewChannelがNewChannelDataではありません"
	return ChannelData{
		ID:        NewChannelDataID(channelID),
		ChannelID: channelID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
