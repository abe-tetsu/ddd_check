package id_analyzer

import (
	"fmt"
	"github.com/abe-tetsu/ddd_check/result"
	"go/ast"
	"strings"
)

func AnalyzerRun(result *result.Result, f *ast.File, fileNameFull, fileName string) *ast.Ident {
	var idIdent *ast.Ident
	isExistIDIdent := false
	ast.Inspect(f, func(n ast.Node) bool {
		if !isExistIDIdent {
			switch n := n.(type) {
			case *ast.GenDecl:
				// IDが定義されていることを確認し、IDの名前を取得して、ファイル名と一致するか確認
				idIdent = IDAnalyzer(n)
				if idIdent != nil {
					isExistIDIdent = true
				}
			}
		}
		return true
	})

	if !isExistIDIdent || idIdent == nil {
		result.IDError = f.Pos()
		result.IDErrorMessage = "ID型で定義されていません"
		return nil
	}

	// ファイル名と一致するか確認
	// TODO: 大文字と小文字を区別しない実装にしている
	identNameLower := strings.ToLower(idIdent.Name)
	fileNameLower := strings.ToLower(fileName)
	if identNameLower != fileNameLower+"id" {
		result.IDError = idIdent.Pos()
		result.IDErrorMessage = fmt.Sprintf("ファイル名%vとID名%vが一致していません", strings.Split(fileNameFull, "/")[len(strings.Split(fileNameFull, "/"))-1], idIdent.Name)
	}

	return idIdent
}

func IDAnalyzer(n *ast.GenDecl) *ast.Ident {
	if len(n.Specs) == 0 {
		return nil
	}

	for _, spec := range n.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			return nil
		}

		// structではないことを確認
		_, ok = typeSpec.Type.(*ast.StructType)
		if ok {
			return nil
		}

		// 取得したIdentがIDか確認
		typeIdent, ok := typeSpec.Type.(*ast.Ident)
		if !ok {
			return nil
		}

		if typeIdent.Name == "ID" {
			return typeSpec.Name
		}
	}

	// ここにきたらID型で定義されていない
	return nil
}
