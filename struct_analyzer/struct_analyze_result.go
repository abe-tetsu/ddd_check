package struct_analyzer

import "go/token"

type StructAnalyzerResult struct {
	// 構造体が正しく定義されているかどうか
	//IsStructExist      bool
	StructError        token.Pos
	StructErrorMessage string

	// 構造体のコンストラクタが正しく定義されているかどうか
	//IsStructConstructorExist      bool
	StructConstructorError        []token.Pos
	StructConstructorErrorMessage []string
}
