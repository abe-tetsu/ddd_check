package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

const docID = "analysis is ..."

// AnalyzerStruct is ...
var AnalyzerID = &analysis.Analyzer{
	Name: "analysis",
	Doc:  docID,
	Run:  runID,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runID(pass *analysis.Pass) (any, error) {
	// パッケージ名がdomainでない場合は処理しない
	if pass.Pkg.Name() != "domain" {
		return nil, nil
	}

	// ファイル名を取得する
	fileNameList := make([]string, 0, len(pass.Files))
	for _, f := range pass.Files {
		file := pass.Fset.File(f.Pos())
		fileNameList = append(fileNameList, file.Name())
	}

	// ファイル単位で解析する
	//fs := token.NewFileSet()
	for _, fileName := range fileNameList {
		f, err := parser.ParseFile(pass.Fset, fileName, nil, 0)
		if err != nil {
			return nil, err
		}

		AnalyzerRun(pass, f, fileName, ConvertFileName(fileName))
	}
	return nil, nil
}

func AnalyzerRun(pass *analysis.Pass, f *ast.File, fileNameFull, fileName string) {
	fmt.Println("AnalyzerRun: ", fileName, "=====================")

	var ident *ast.Ident
	isExist := false
	ast.Inspect(f, func(n ast.Node) bool {
		if !isExist {
			switch n := n.(type) {
			case *ast.GenDecl:
				// IDが定義されていることを確認し、IDの名前を取得して、ファイル名と一致するか確認
				ident = IDAnalyzer(n)
				if ident != nil {
					isExist = true
					return true
				}
			}
		}
		return true
	})

	if !isExist || ident == nil {
		pass.Reportf(f.Pos(), "ID型で定義されていません")
		return
	}

	// ファイル名と一致するか確認
	if ident.Name != fileName+"ID" {
		pass.Reportf(ident.Pos(), "ファイル名%vとID名%vが一致していません", strings.Split(fileNameFull, "/")[len(strings.Split(fileNameFull, "/"))-1], ident.Name)
		return
	}

	// IDのコンストラクタが定義されていることを確認する
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
