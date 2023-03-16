package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

type Result struct {
	// IDが正しく定義されているかどうか
	//IsIDExist      bool
	IDError        token.Pos
	IDErrorMessage string

	// IDのコンストラクタが正しく定義されているかどうか
	//IsIDConstructorExist      bool
	IDConstructorError        token.Pos
	IDConstructorErrorMessage string

	// 構造体が正しく定義されているかどうか
	//IsStructExist      bool
	StructError        token.Pos
	StructErrorMessage string

	// 構造体のコンストラクタが正しく定義されているかどうか
	//IsStructConstructorExist      bool
	StructConstructorError        token.Pos
	StructConstructorErrorMessage string
}

// ケース
// 1. そもそも ID 型で定義されていない => コンストラクタをみる必用がない
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
		// 結果を格納する構造体を初期化
		result := Result{}

		f, err := parser.ParseFile(pass.Fset, fileName, nil, 0)
		if err != nil {
			return nil, err
		}

		idIdent := AnalyzerRun(&result, f, fileName, ConvertFileName(fileName))
		if idIdent != nil {
			// IDのコンストラクタが定義されていることを確認
			AnalyzerIDConstructorRun(&result, idIdent, f, ConvertFileName(fileName))
		}

		structIdent := StructAnalyzerRun(&result, f, fileName, ConvertFileName(fileName))
		if structIdent != nil {
			//fmt.Println("structIdent", structIdent)
			// 構造体のコンストラクタが定義されていることを確認
			AnalyzerStructConstructorRun(&result, structIdent, f, ConvertFileName(fileName))
		}

		// 結果をレポート
		reportResult(result, pass)
	}
	return nil, nil
}

func AnalyzerRun(result *Result, f *ast.File, fileNameFull, fileName string) *ast.Ident {
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
	if idIdent.Name != fileName+"ID" {
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

func AnalyzerIDConstructorRun(result *Result, idIdent *ast.Ident, f *ast.File, fileName string) {
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

	// コンストラクタの名前がNew+ファイル名+IDであることを確認
	if constructorIdent.Name != "New"+fileName+"ID" {
		result.IDConstructorError = constructorIdent.Pos()
		result.IDConstructorErrorMessage = fmt.Sprintf("コンストラクタ名%vがNew%vIDではありません", constructorIdent.Name, fileName)
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

func StructAnalyzerRun(result *Result, f *ast.File, fileNameFull, fileName string) *ast.Ident {
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
	//fmt.Println("比較します", "ファイル名", fileName, "構造体名", structIdent.Name)
	if structIdent.Name != fileName {
		result.StructError = structIdent.Pos()
		result.StructErrorMessage = fmt.Sprintf("ファイル名%vとID名%vが一致していません", strings.Split(fileNameFull, "/")[len(strings.Split(fileNameFull, "/"))-1], structIdent.Name)
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

func AnalyzerStructConstructorRun(result *Result, structIdent *ast.Ident, f *ast.File, fileName string) {
	isExistConstructorIdent := false
	var constructorIdent *ast.Ident

	// 構造体のコンストラクタが定義されていることを確認する
	ast.Inspect(f, func(n ast.Node) bool {
		if !isExistConstructorIdent {
			switch n := n.(type) {
			case *ast.FuncDecl:
				// コンストラクタが定義されていることを確認
				constructorIdent = ConstructorAnalyzer(n, structIdent)
				if constructorIdent != nil {
					isExistConstructorIdent = true
				}
			}
		}
		return true
	})

	if !isExistConstructorIdent || constructorIdent == nil {
		result.StructConstructorError = f.Pos()
		result.StructConstructorErrorMessage = "構造体のコンストラクタが定義されていません"
		return
	}

	// コンストラクタの名前がNew+ファイル名であることを確認
	if constructorIdent.Name != "New"+fileName {
		result.StructConstructorError = constructorIdent.Pos()
		result.StructConstructorErrorMessage = fmt.Sprintf("コンストラクタ名%vがNew%vではありません", constructorIdent.Name, fileName)
		return
	}
}

func reportResult(result Result, pass *analysis.Pass) {
	// IDと構造体が存在しない時にファイルの先頭にエラーを出す
	if result.IDErrorMessage == "ID型で定義されていません" && result.StructErrorMessage == "構造体が定義されていません" {
		pass.Reportf(result.IDError, "ID型と構造体が定義されていません")
		return
	}

	// IDのコンストラクタと構造体のコンストラクタが存在しない時にファイルの先頭にエラーを出す
	if result.IDConstructorErrorMessage == "IDのコンストラクタが定義されていません" && result.StructConstructorErrorMessage == "構造体のコンストラクタが定義されていません" {
		pass.Reportf(result.IDConstructorError, "IDと構造体のコンストラクタの定義がされていません")
		return
	}

	if result.IDErrorMessage != "" {
		pass.Reportf(result.IDError, result.IDErrorMessage)
	}

	if result.IDConstructorErrorMessage != "" {
		pass.Reportf(result.IDConstructorError, result.IDConstructorErrorMessage)
	}

	if result.StructErrorMessage != "" {
		pass.Reportf(result.StructError, result.StructErrorMessage)
	}

	if result.StructConstructorErrorMessage != "" {
		pass.Reportf(result.StructConstructorError, result.StructConstructorErrorMessage)
	}
}
