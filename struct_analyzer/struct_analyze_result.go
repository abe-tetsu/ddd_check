package struct_analyzer

import "go/token"

type StructAnalyzerResult struct {
	// 構造体が正しく定義されているかどうか
	StructError        token.Pos
	StructErrorMessage string

	// 構造体のコンストラクタが正しく定義されているかどうか
	StructConstructorError        []token.Pos
	StructConstructorErrorMessage []string
}
