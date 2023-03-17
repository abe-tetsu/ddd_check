package ddd_check

import (
	"github.com/abe-tetsu/ddd_check/id_analyzer"
	"github.com/abe-tetsu/ddd_check/result"
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

// ケース
// 1. そもそも ID 型で定義されていない => コンストラクタをみる必用がない
// // // 1-1-1. 構造体も定義されていない　=> コンストラクタをみる必要がない
// // // 1-1-2. 構造体が定義されているが、名前がファイル名と一致していない
// // // // 1-1-2-1. コンストラクタが定義されていない
// // // // 1-1-2-2. コンストラクタが定義されているが、名前がファイル名と一致していない
// // // // 1-1-2-3. コンストラクタが定義されていて、名前がファイル名と一致している
// // // 1-1-3. 構造体が定義されていて、名前がファイル名と一致している
// // // // 1-1-3-1. コンストラクタが定義されていない
// // // // 1-1-3-2. コンストラクタが定義されているが、名前がファイル名と一致していない
// // // // 1-1-3-3. コンストラクタが定義されていて、名前がファイル名と一致している

// 2. ID 型で定義されているが、名前がファイル名と一致していない => コンストラクタも見て、名前が一致しているか確認
// // 2-1. コンストラクタが定義されていない
// // 2-2. コンストラクタが定義されているが、名前がファイル名と一致していない
// // 2-3. コンストラクタが定義されていて、名前がファイル名と一致している
// 3. ID 型で定義されていて、名前がファイル名と一致している => コンストラクタも見て、名前が一致しているか確認
// // 3-1. コンストラクタが定義されていない
// // 3-2. コンストラクタが定義されているが、名前がファイル名と一致していない
// // 3-3. コンストラクタが定義されていて、名前がファイル名と一致している

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
		result := result.Result{}

		f, err := parser.ParseFile(pass.Fset, fileName, nil, 0)
		if err != nil {
			return nil, err
		}

		idIdent := id_analyzer.AnalyzerRun(&result, f, fileName, ConvertFileName(fileName))
		if idIdent != nil {
			// IDのコンストラクタが定義されていることを確認
			id_analyzer.AnalyzerIDConstructorRun(&result, idIdent, f, ConvertFileName(fileName))
		}

		structIdent := struct_analyzer.StructAnalyzerRun(&result, f, fileName, ConvertFileName(fileName))
		if structIdent != nil {
			// 構造体のコンストラクタが定義されていることを確認
			struct_analyzer.AnalyzerStructConstructorRun(&result, structIdent, f, ConvertFileName(fileName))
		}

		// 結果をレポート
		ReportResult(result, pass)
	}
	return nil, nil
}
