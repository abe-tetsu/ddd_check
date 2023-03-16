package analysis

import "strings"

// ファイル名をを / で分割して最後の要素を取得する
// その要素の .go を削除して、キャメルケースに変換する
func ConvertFileName(fileNameOld string) string {
	fileNameOnly := strings.Split(fileNameOld, "/")
	fileName := strings.Replace(fileNameOnly[len(fileNameOnly)-1], ".go", "", -1)

	// ローワーケースからアッパーキャメルケースに変換
	// 例: user_id -> UserID
	return strings.Replace(strings.Title(strings.Replace(fileName, "_", " ", -1)), " ", "", -1)
}
