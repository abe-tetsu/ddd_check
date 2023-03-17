package id_analyzer

import (
	"fmt"
	"go/ast"
	"strings"
)

func AnalyzerIDConstructorRun(result *IDAnalyzerResult, idIdent *ast.Ident, f *ast.File, fileName string) {
	isExistConstructorIdent := false
	var constructorIdent *ast.Ident

	// IDのコンストラクタが定義されていることを確認する
	ast.Inspect(f, func(n ast.Node) bool {
		if !isExistConstructorIdent {
			switch n := n.(type) {
			case *ast.FuncDecl:
				// コンストラクタが定義されていることを確認
				constructorIdent = ConstructorAnalyzer(n, idIdent)
				if constructorIdent != nil {
					isExistConstructorIdent = true
				}
			}
		}
		return true
	})

	if !isExistConstructorIdent || constructorIdent == nil {
		result.IDConstructorError = f.Pos()
		result.IDConstructorErrorMessage = "IDのコンストラクタが定義されていません"
		return
	}

	// ファイル名と一致するか確認
	// TODO: 大文字と小文字を区別しない実装にしている
	constructorNameLower := strings.ToLower(constructorIdent.Name)
	fileNameLower := strings.ToLower("New" + fileName + "ID")
	if constructorNameLower != fileNameLower {
		result.IDConstructorError = constructorIdent.Pos()
		result.IDConstructorErrorMessage = fmt.Sprintf("ファイル名%vとコンストラクタ名%vが一致していません", fileName, constructorIdent.Name)
	}
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
