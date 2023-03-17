package ddd_check

import (
	"github.com/abe-tetsu/ddd_check/id_analyzer"
	"github.com/abe-tetsu/ddd_check/struct_analyzer"
	"go/parser"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

const docID = "analysis is checker for DDD"

// AnalyzerIDConstructor is ...
var AnalyzerID = &analysis.Analyzer{
	Name: "ddd_checker",
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
		// ファイル名に _test が含まれている場合は処理しない
		if strings.Contains(fileName, "_test") || strings.Contains(fileName, "_payload") {
			continue
		}

		// 結果を格納する構造体を初期化
		idAnalyzeResult := id_analyzer.IDAnalyzerResult{}
		structAnalyzeResult := struct_analyzer.StructAnalyzerResult{}

		f, err := parser.ParseFile(pass.Fset, fileName, nil, 0)
		if err != nil {
			return nil, err
		}

		idIdent := id_analyzer.IDAnalyzerRun(&idAnalyzeResult, f, fileName, ConvertFileName(fileName))
		if idIdent != nil {
			// IDのコンストラクタが定義されていることを確認
			id_analyzer.IDConstructorAnalyzerRun(&idAnalyzeResult, idIdent, f, ConvertFileName(fileName))
		}

		structIdent := struct_analyzer.StructAnalyzerRun(&structAnalyzeResult, f, fileName, ConvertFileName(fileName))
		if structIdent != nil {
			// 構造体のコンストラクタが定義されていることを確認
			struct_analyzer.StructConstructorAnalyzerRun(&structAnalyzeResult, structIdent, f, ConvertFileName(fileName))
		}

		// 結果をレポート
		ReportResult(idAnalyzeResult, structAnalyzeResult, pass)
	}
	return nil, nil
}
