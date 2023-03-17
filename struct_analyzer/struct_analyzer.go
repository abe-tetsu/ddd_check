package struct_analyzer

import (
	"fmt"
	"github.com/abe-tetsu/ddd_check/result"
	"go/ast"
	"strings"
)

func StructAnalyzerRun(result *result.Result, f *ast.File, fileNameFull, fileName string) *ast.Ident {
	var structIdent *ast.Ident
	isExistStructIdent := false
	ast.Inspect(f, func(n ast.Node) bool {
		if !isExistStructIdent {
			switch n := n.(type) {
			case *ast.GenDecl:
				// 構造体が定義されていることを確認し、構造体の名前を取得する
				structIdent = StructAnalyzer(n)
				if structIdent != nil {
					isExistStructIdent = true
				}
			}
		}
		return true
	})

	if !isExistStructIdent || structIdent == nil {
		result.StructError = f.Pos()
		result.StructErrorMessage = "構造体が定義されていません"
		return nil
	}

	// ファイル名と一致するか確認
	structIdentNameLower := strings.ToLower(structIdent.Name)
	fileNameLower := strings.ToLower(fileName)
	if structIdentNameLower != fileNameLower {
		result.StructError = structIdent.Pos()
		result.StructErrorMessage = fmt.Sprintf("ファイル名%vと構造体名%vが一致していません", strings.Split(fileNameFull, "/")[len(strings.Split(fileNameFull, "/"))-1], structIdent.Name)
	}

	return structIdent
}

func StructAnalyzer(n *ast.GenDecl) *ast.Ident {
	if len(n.Specs) == 0 {
		return nil
	}

	for _, spec := range n.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			return nil
		}

		// structであることを確認
		_, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			return nil
		}

		return typeSpec.Name
	}

	// ここにきたら構造体が定義されていない
	return nil
}
