package id_analyzer

import "go/token"

type IDAnalyzerResult struct {
	// IDが正しく定義されているかどうか
	//IsIDExist      bool
	IDError        token.Pos
	IDErrorMessage string

	// IDのコンストラクタが正しく定義されているかどうか
	//IsIDConstructorExist      bool
	IDConstructorError        token.Pos
	IDConstructorErrorMessage string
}
