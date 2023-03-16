package analysis

import (
	"go/ast"
	"go/parser"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

const docID = "analysis is ..."

// AnalyzerIDConstructor is ...
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

		idIdent := AnalyzerRun(pass, f, fileName, ConvertFileName(fileName))

		if idIdent == nil {
			continue
		}

		// IDのコンストラクタが定義されていることを確認
		AnalyzerIDConstructorRun(pass, idIdent, f, fileName, ConvertFileName(fileName))

	}
	return nil, nil
}

// 1. そもそも ID 型で定義されていない => コンストラクタをみる必用がない
// 2. ID 型で定義されているが、名前がファイル名と一致していない => コンストラクタも見て、名前が一致しているか確認
// // 2-1. コンストラクタが定義されていない
// // 2-2. コンストラクタが定義されているが、名前がファイル名と一致していない
// // 2-3. コンストラクタが定義されていて、名前がファイル名と一致している
// 3. ID 型で定義されていて、名前がファイル名と一致している => コンストラクタも見て、名前が一致しているか確認
// // 3-1. コンストラクタが定義されていない
// // 3-2. コンストラクタが定義されているが、名前がファイル名と一致していない
// // 3-3. コンストラクタが定義されていて、名前がファイル名と一致している
func AnalyzerRun(pass *analysis.Pass, f *ast.File, fileNameFull, fileName string) *ast.Ident {
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
		pass.Reportf(f.Pos(), "ID型で定義されていません")
		return nil
	}

	// ファイル名と一致するか確認
	if idIdent.Name != fileName+"ID" {
		pass.Reportf(idIdent.Pos(), "ファイル名%vとID名%vが一致していません", strings.Split(fileNameFull, "/")[len(strings.Split(fileNameFull, "/"))-1], idIdent.Name)
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

func AnalyzerIDConstructorRun(pass *analysis.Pass, idIdent *ast.Ident, f *ast.File, fileNameFull, fileName string) {
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
		pass.Reportf(f.Pos(), "IDのコンストラクタが定義されていません")
		return
	}

	// コンストラクタの名前がNew+ファイル名+IDであることを確認
	if constructorIdent.Name != "New"+fileName+"ID" {
		pass.Reportf(constructorIdent.Pos(), "コンストラクタ名%vがNew%vIDではありません", constructorIdent.Name, fileName)
		return
	}
}

func ConstructorAnalyzer(n *ast.FuncDecl, idIdent *ast.Ident) *ast.Ident {
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
