package ast

import (
	"FoxLite/lang/token"
	"bytes"
	"fmt"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) statementNode() {}
func (p *Program) String() string {
	var out bytes.Buffer
	if len(p.Statements) > 0 {
		for _, stmt := range p.Statements {
			out.WriteString(fmt.Sprintf("%s\n", stmt.String()))
		}
	}
	return out.String()
}

type BlockStmt struct {
	Statements []Statement
	Level      int
}

func (b *BlockStmt) statementNode() {}
func (b *BlockStmt) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	var tab string
	for i := 1; i <= b.Level; i++ {
		tab += "    "
	}
	if len(b.Statements) > 0 {
		for _, stmt := range b.Statements {
			out.WriteString(fmt.Sprintf("%s%s\n", tab, stmt.String()))
		}
	}
	out.WriteString("}")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}

type VarStmt struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (v *VarStmt) statementNode() {}
func (v *VarStmt) String() string {
	return fmt.Sprintf("var %s = %s", v.Name.String(), v.Value.String())
}

type ReturnStmt struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStmt) statementNode() {}
func (r *ReturnStmt) String() string {
	return fmt.Sprintf("return %s", r.Value.String())
}

type PrefixExpression struct {
	Token token.Token
	Op    string
	Right Expression
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s %s)", p.Op, p.Right.String())
}

type InfixExpression struct {
	Token token.Token
	Left  Expression
	Op    string
	Right Expression
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Op, i.Right.String())
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) String() string {
	return s.Value
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) String() string {
	if b.Value {
		return ".T."
	} else {
		return ".F."
	}
}

type NilLiteral struct {
	Token token.Token
}

func (n *NilLiteral) expressionNode() {}
func (n *NilLiteral) String() string {
	return "nil"
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStmt
	Alternative *BlockStmt
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if(")
	out.WriteString(i.Condition.String())
	out.WriteString(")")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}
