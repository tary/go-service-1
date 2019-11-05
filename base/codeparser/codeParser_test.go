package codeparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"net/url"
	"strings"
	"testing"
)

type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.Ident:
		fmt.Println("*ast.Ident, name: ", d.Name)

	case *ast.BasicLit:
		fmt.Println("*ast.BasicLit, Value: ", d.Value)

	case *ast.CallExpr:
		fmt.Println("*ast.CallExpr, Fun: ", d.Fun, ", Args: ", d.Args)

	case *ast.SelectorExpr:
		fmt.Println("*ast.SelectorExpr, Name: ", d.Sel.Name, ", Obj: ", d.Sel.Obj)

	default:
		fmt.Printf("default, %s%T\n", strings.Repeat("\t", int(v)), n)
	}

	//fmt.Printf("%s%T\n", strings.Repeat("\t", int(v)), n)
	return v + 1
}

func TestParser(t *testing.T) {
	expr, err := parser.ParseExpr("GetPlayerByName(1).ModulePlayerInfo..Test(\"haha\", 1, 2, 33, 444)")

	if err != nil {
		fmt.Println("expr: ", expr, ", err: ", err)
	} else {
		fmt.Println("expr: ", expr, ", Pos:", expr.Pos(), ", End: ", expr.End())
	}

	expr2, err2 := parser.ParseExpr("ModulePlayerInfo.Test(\"haha\", 1, 2, 33, 444)")

	if err2 != nil {
		fmt.Println("expr: ", expr2, ", err2: ", err2)
	} else {
		fmt.Println("expr: ", expr2, ", Pos:", expr2.Pos(), ", End: ", expr2.End())
	}

	expr3, err3 := parser.ParseExpr("GetPlayerByName(1).ModulePlayerInfo.Test(\"h,a.,,h.,a,\", 1, 2, 33.0, 444.1).Level")

	if err3 != nil {
		fmt.Println("expr: ", expr3, ", err3: ", err3)
	} else {
		fmt.Println("expr: ", expr3, ", Pos:", expr3.Pos(), ", End: ", expr3.End())
	}

	var v visitor
	ast.Walk(v, expr3)

	ast.Inspect(expr3, func(node ast.Node) bool {
		//fmt.Printf("%s%T\n", strings.Repeat("\t", int(v)), node)
		return true
	})
}
func TestScanner(t *testing.T) {
	src := []byte("GetPlayerByName(1).ModulePlayerInfo.Test(\"h,a.,,h.,a,\", 1, 2, 33.0, 444.1).Level")

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}
}

func TestNewParser(t *testing.T) {
	src := "GetPlayerByName(1).ModulePlayerInfo.Test(\"haha\", 1, 2, 33, 444)"

	dest := url.QueryEscape(src)
	newSrc, err1 := url.QueryUnescape(dest)

	fmt.Println("dest: ", dest, ", err: ", err1)
	fmt.Println("newSrc: ", newSrc)

	nodes, err := ParseCode(src)
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("len: ", len(nodes))
}

func TestPlayer1Parser(t *testing.T) {
	src := "GetPlayerByName().ModulePlayerInfo.Test(\"haha\", 1, 2, 33, 444).Level.Name"

	dest := url.QueryEscape(src)
	newSrc, err1 := url.QueryUnescape(dest)

	fmt.Println("dest: ", dest, ", err: ", err1)
	fmt.Println("newSrc: ", newSrc)

	nodes, err := ParseCode(src)
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("len: ", len(nodes))
}

func TestPlayerGM(t *testing.T) {
	src := "ModulePlayerInfo.GetLevel()"

	dest := url.QueryEscape(src)
	newSrc, err1 := url.QueryUnescape(dest)

	fmt.Println("dest: ", dest, ", err: ", err1)
	fmt.Println("newSrc: ", newSrc)

	nodes, err := ParseCode(src)
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("len: ", len(nodes))
}

func TestPlayer2GM(t *testing.T) {
	src := "ModulePlayerInfo.Level"

	dest := url.QueryEscape(src)
	newSrc, err1 := url.QueryUnescape(dest)

	fmt.Println("dest: ", dest, ", err: ", err1)
	fmt.Println("newSrc: ", newSrc)

	nodes, err := ParseCode(src)
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("len: ", len(nodes))
}

func TestPlayer3GM(t *testing.T) {
	src := "Level"

	dest := url.QueryEscape(src)
	newSrc, err1 := url.QueryUnescape(dest)

	fmt.Println("dest: ", dest, ", err: ", err1)
	fmt.Println("newSrc: ", newSrc)

	nodes, err := ParseCode(src)
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Println("len: ", len(nodes))
}
