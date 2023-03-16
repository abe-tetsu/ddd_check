package analysis

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const docStruct = "analysis is ..."

// AnalyzerIDConstructor is ...
var AnalyzerStruct = &analysis.Analyzer{
	Name: "analysis",
	Doc:  docStruct,
	Run:  runStruct,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runStruct(pass *analysis.Pass) (any, error) {
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
	////fs := token.NewFileSet()
	//for _, fileName := range fileNameList {
	//	f, err := parser.ParseFile(pass.Fset, fileName, nil, 0)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	idIdent := AnalyzerRun(pass, f, fileName, ConvertFileName(fileName))
	//
	//	if idIdent == nil {
	//		continue
	//	}
	//
	//	// IDのコンストラクタが定義されていることを確認
	//	AnalyzerIDConstructorRun(pass, idIdent, f, fileName, ConvertFileName(fileName))
	//}
	return nil, nil
}
