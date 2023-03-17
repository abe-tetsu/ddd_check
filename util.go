package ddd_check

import (
	"github.com/abe-tetsu/ddd_check/id_analyzer"
	"github.com/abe-tetsu/ddd_check/struct_analyzer"
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

func ReportResult(idAnalyzeResult id_analyzer.IDAnalyzerResult, structAnalyzeResult struct_analyzer.StructAnalyzerResult, pass *analysis.Pass) {
	// IDと構造体が存在しない時にファイルの先頭にエラーを出す
	if idAnalyzeResult.IDErrorMessage == "ID型で定義されていません" && structAnalyzeResult.StructErrorMessage == "構造体が定義されていません" {
		pass.Reportf(idAnalyzeResult.IDError, "ID型と構造体が定義されていません")
		return
	}

	// IDのコンストラクタと構造体のコンストラクタが存在しない時にファイルの先頭にエラーを出す
	if idAnalyzeResult.IDConstructorErrorMessage == "IDのコンストラクタが定義されていません" && len(structAnalyzeResult.StructConstructorErrorMessage) != 0 {
		if structAnalyzeResult.StructConstructorErrorMessage[0] == "構造体のコンストラクタが定義されていません" {
			pass.Reportf(idAnalyzeResult.IDConstructorError, "IDと構造体のコンストラクタの定義がされていません")
		}
		return
	}

	if idAnalyzeResult.IDErrorMessage != "" {
		pass.Reportf(idAnalyzeResult.IDError, idAnalyzeResult.IDErrorMessage)
	}

	if idAnalyzeResult.IDConstructorErrorMessage != "" {
		pass.Reportf(idAnalyzeResult.IDConstructorError, idAnalyzeResult.IDConstructorErrorMessage)
	}

	if structAnalyzeResult.StructErrorMessage != "" {
		pass.Reportf(structAnalyzeResult.StructError, structAnalyzeResult.StructErrorMessage)
	}

	if len(structAnalyzeResult.StructConstructorErrorMessage) != 0 {
		for i, v := range structAnalyzeResult.StructConstructorErrorMessage {
			pass.Reportf(structAnalyzeResult.StructConstructorError[i], v)
		}
	}
}
