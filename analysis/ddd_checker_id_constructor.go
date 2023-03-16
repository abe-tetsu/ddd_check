package analysis

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const docIDConstructor = "analysis is ..."

// AnalyzerIDConstructor is ...
var AnalyzerIDConstructor = &analysis.Analyzer{
	Name: "analysis",
	Doc:  docIDConstructor,
	Run:  runIDConstructor,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runIDConstructor(pass *analysis.Pass) (any, error) {
	//// パッケージ名がdomainでない場合は処理しない
	//if pass.Pkg.Name() != "domain" {
	//	return nil, nil
	//}
	//
	//// ファイル名を取得する
	//fileNameList := make([]string, 0, len(pass.Files))
	//for _, f := range pass.Files {
	//	file := pass.Fset.File(f.Pos())
	//	fileNameList = append(fileNameList, file.Name())
	//}
	//
	//// ファイル単位で解析する
	//for _, fileName := range fileNameList {
	//	f, err := parser.ParseFile(pass.Fset, fileName, nil, 0)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	AnalyzerIDConstructorRun(pass, f, fileName, ConvertFileName(fileName))
	//}
	return nil, nil
}

func AnalyzerIDConstructorRun(pass *analysis.Pass, idIdent *ast.Ident, f *ast.File, fileNameFull, fileName string) {
	fmt.Println("AnalyzerRun Constructor: ", fileName, " Start =====================")

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

	fmt.Println("AnalyzerRun: ", fileName, "===================== end")
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

	fmt.Println()
	fmt.Println("比較します returnType.Name: ", returnType.Name, " idIdent.Name: ", idIdent.Name)
	if returnType.Name != idIdent.Name {
		return nil
	}

	// コンストラクタの名前を取得
	return n.Name
}
