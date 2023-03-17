package struct_analyzer

import (
	"fmt"
	"github.com/abe-tetsu/ddd_check/result"
	"go/ast"
	"strings"
)

func AnalyzerStructConstructorRun(result *result.Result, structIdent *ast.Ident, f *ast.File, fileName string) {
	var constructorIdents []*ast.Ident

	// 構造体のコンストラクタが定義されていることを確認する
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			// コンストラクタが定義されていることを確認
			constructorIdents = append(constructorIdents, ConstructorAnalyzer(n, structIdent))
		}

		return true
	})

	if len(constructorIdents) == 0 {
		//fmt.Println("構造体のコンストラクタが定義されていません")
		result.StructConstructorError = append(result.StructConstructorError, f.Pos())
		result.StructConstructorErrorMessage = append(result.StructConstructorErrorMessage, "構造体のコンストラクタが定義されていません")
		return
	}

	//fmt.Println("コンストラクタの数: ", len(constructorIdents))
	for _, constructorIdent := range constructorIdents {
		if constructorIdent == nil {
			continue
		}

		constructorIdentNameLower := strings.ToLower(constructorIdent.Name)
		fileNameLower := strings.ToLower(fileName)
		// コンストラクタは複数個定義されている可能性がある。その場合は、コンストラクタ名にファイル名が含まれているかどうかで判定する。
		if !strings.Contains(constructorIdentNameLower, fileNameLower) {
			result.StructConstructorError = append(result.StructConstructorError, constructorIdent.Pos())
			result.StructConstructorErrorMessage = append(result.StructConstructorErrorMessage, fmt.Sprintf("コンストラクタ名%vにファイル名%vが含まれていません", constructorIdent.Name, fileName))
		}
	}
	return
}

func ConstructorAnalyzer(n *ast.FuncDecl, idIdent *ast.Ident) *ast.Ident {
	if n == nil {
		return nil
	}

	if n.Type == nil {
		return nil
	}

	if n.Type.Results == nil {
		return nil
	}

	// 返り値が1つでない場合は処理しない
	if len(n.Type.Results.List) != 1 {
		return nil
	}

	// 返り値がID型でない場合は処理しない
	returnType, ok := n.Type.Results.List[0].Type.(*ast.Ident)
	if !ok {
		return nil
	}

	if returnType.Name != idIdent.Name {
		return nil
	}

	// コンストラクタの名前を取得
	return n.Name
}
