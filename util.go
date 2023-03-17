package ddd_check

import (
	"github.com/abe-tetsu/ddd_check/result"
	"golang.org/x/tools/go/analysis"
	"strings"
)

// ファイル名をを / で分割して最後の要素を取得する
// その要素の .go を削除して、キャメルケースに変換する
func ConvertFileName(fileNameOld string) string {
	fileNameOnly := strings.Split(fileNameOld, "/")
	fileName := strings.Replace(fileNameOnly[len(fileNameOnly)-1], ".go", "", -1)

	// ローワーケースからアッパーキャメルケースに変換
	// 例: user_id -> UserID
	return strings.Replace(strings.Title(strings.Replace(fileName, "_", " ", -1)), " ", "", -1)
}

func ReportResult(result result.Result, pass *analysis.Pass) {
	// IDと構造体が存在しない時にファイルの先頭にエラーを出す
	if result.IDErrorMessage == "ID型で定義されていません" && result.StructErrorMessage == "構造体が定義されていません" {
		pass.Reportf(result.IDError, "ID型と構造体が定義されていません")
		return
	}

	// IDのコンストラクタと構造体のコンストラクタが存在しない時にファイルの先頭にエラーを出す
	if result.IDConstructorErrorMessage == "IDのコンストラクタが定義されていません" && len(result.StructConstructorErrorMessage) != 0 {
		if result.StructConstructorErrorMessage[0] == "構造体のコンストラクタが定義されていません" {
			pass.Reportf(result.IDConstructorError, "IDと構造体のコンストラクタの定義がされていません")
		}
		return
	}

	if result.IDErrorMessage != "" {
		pass.Reportf(result.IDError, result.IDErrorMessage)
	}

	if result.IDConstructorErrorMessage != "" {
		pass.Reportf(result.IDConstructorError, result.IDConstructorErrorMessage)
	}

	if result.StructErrorMessage != "" {
		pass.Reportf(result.StructError, result.StructErrorMessage)
	}

	if len(result.StructConstructorErrorMessage) != 0 {
		for i, v := range result.StructConstructorErrorMessage {
			pass.Reportf(result.StructConstructorError[i], v)
		}
	}
}
