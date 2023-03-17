package result

import "go/token"

type Result struct {
	// IDが正しく定義されているかどうか
	//IsIDExist      bool
	IDError        token.Pos
	IDErrorMessage string

	// IDのコンストラクタが正しく定義されているかどうか
	//IsIDConstructorExist      bool
	IDConstructorError        token.Pos
	IDConstructorErrorMessage string

	// 構造体が正しく定義されているかどうか
	//IsStructExist      bool
	StructError        token.Pos
	StructErrorMessage string

	// 構造体のコンストラクタが正しく定義されているかどうか
	//IsStructConstructorExist      bool
	StructConstructorError        []token.Pos
	StructConstructorErrorMessage []string
}
