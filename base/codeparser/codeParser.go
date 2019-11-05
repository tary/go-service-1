package codeparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"net/url"
	"runtime/debug"

	"github.com/giant-tech/go-service/base/utility"

	"github.com/cihub/seelog"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

// newCodeNode new code
func newCodeNode(nodeName string) *CodeNode {
	isFunc := (len(nodeName) == 0)
	return &CodeNode{IsFunc: isFunc, NodeName: nodeName}
}

// CodeVisitor 代码访问器
type CodeVisitor struct {
	Nodes      []*CodeNode
	CurNodeIdx int
	HaveIdent  bool
}

// ParseCode 解析代码
func ParseCode(content string) (nodes []*CodeNode, err error) {
	defer func() {
		if e := recover(); e != nil {
			seelog.Error("ParseCode panic:", e, string(debug.Stack()))
			if viper.GetString("Config.Recover") == "0" {
				panic(e)
			}

			retStr := fmt.Sprintf("%+v", e)
			err = fmt.Errorf(retStr)
		}
	}()

	newStr, err := url.QueryUnescape(content)
	if err != nil {
		seelog.Error("QueryUnescape err: ", err)
		return nodes, err
	}

	v := newCodeVisitor()
	expr, parseErr := parser.ParseExpr(newStr)
	if parseErr != nil {
		err = parseErr
		seelog.Error("ParseCode err: ", err)
		return nodes, parseErr
	}

	ast.Walk(v, expr)

	v.CurNodeIdx = 0

	if len(v.Nodes) > 0 {
		//翻转nodes
		for i, j := 0, len(v.Nodes)-1; i < j; i, j = i+1, j-1 {
			v.Nodes[i], v.Nodes[j] = v.Nodes[j], v.Nodes[i]
		}
	}

	return v.Nodes, nil
}

// newCodeVisitor new code vis
func newCodeVisitor() *CodeVisitor {
	return &CodeVisitor{CurNodeIdx: -1}
}

// Visit 访问节点
func (v *CodeVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.Ident: //所有名字
		//fmt.Println("*ast.Ident, name: ", d.Name)

		//第一个属性或者函数需要特殊处理
		if !v.HaveIdent {
			v.HaveIdent = true

			if v.CurNodeIdx < 0 {
				v.Nodes = append(v.Nodes, newCodeNode(d.Name))
				v.CurNodeIdx++
			}

			if v.Nodes[v.CurNodeIdx].IsFunc && len(v.Nodes[v.CurNodeIdx].NodeName) == 0 {
				v.Nodes[v.CurNodeIdx].NodeName = d.Name
			}

			if v.Nodes[v.CurNodeIdx].NodeName != d.Name {
				v.Nodes = append(v.Nodes, newCodeNode(d.Name))
				v.CurNodeIdx++
			}
		}

		for v.CurNodeIdx >= 0 && v.Nodes[v.CurNodeIdx].NodeName != d.Name {
			v.CurNodeIdx--
		}

	case *ast.BasicLit: //函数参数
		//fmt.Println("*ast.BasicLit, Value: ", d.Value)
		if v.Nodes[v.CurNodeIdx].IsFunc {
			v.Nodes[v.CurNodeIdx].Params = append(v.Nodes[v.CurNodeIdx].Params, utility.Unquote(d.Value))
		} else {
			seelog.Error("*ast.BasicLit, not function")
		}

	case *ast.CallExpr: //函数()
		//fmt.Println("*ast.CallExpr, Fun: ", d.Fun, ", Args: ", d.Args)
		v.Nodes = append(v.Nodes, newCodeNode(""))
		v.CurNodeIdx++

	case *ast.SelectorExpr: //.后面的
		//fmt.Println("*ast.SelectorExpr, Name: ", d.Sel.Name, ", Obj: ", d.Sel.Obj)
		if v.CurNodeIdx >= 0 && v.Nodes[v.CurNodeIdx].IsFunc && len(v.Nodes[v.CurNodeIdx].NodeName) == 0 {
			v.Nodes[v.CurNodeIdx].NodeName = d.Sel.Name
		} else {
			v.Nodes = append(v.Nodes, newCodeNode(d.Sel.Name))
			v.CurNodeIdx++
		}

	default:
		log.Errorf("default, %T\n", n)
	}

	return v
}
