package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"bufio"
	"fmt"
	"os"
)

func evalInputStmt(node *ast.Input, env *object.Environment) object.Object {
	prompt := Eval(node.Message, env)
	if isError(prompt) {
		return prompt
	}
	if prompt.Type() != object.StringObj {
		return object.NewError(fmt.Sprintf("cannot print out `%s`: expected `string`", prompt.Type()))
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt.Inspect()) // mostrar el mensaje
	if scanner.Scan() {
		val := &object.String{Value: scanner.Text()}
		// guardar el string
		env.Set(node.Output.Value.(string), 'l', val)
	}
	if scanner.Err() != nil {
		return object.NewError(fmt.Sprintf("could not read from console: %v", scanner.Err()))
	}
	return None
}
