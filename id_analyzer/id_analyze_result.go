package id_analyzer

import "go/token"

type IDAnalyzerResult struct {
	// IDが正しく定義されているかどうか
	IDError        token.Pos
	IDErrorMessage string

	// IDのコンストラクタが正しく定義されているかどうか
	IDConstructorError        token.Pos
	IDConstructorErrorMessage string
}
