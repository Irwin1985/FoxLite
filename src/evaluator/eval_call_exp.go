package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"fmt"
)

func evalCallExpression(node *ast.CallExp, env *object.Environment) object.Object {
	function := Eval(node.Caller, env)
	if isError(function) {
		return function
	}
	// Evaluamos los argumentos (si es que existen)
	args := evalExpressions(node.Args, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return applyFunction(function, args)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, exp := range exps {
		res := Eval(exp, env)
		if isError(res) {
			return []object.Object{res}
		}
		result = append(result, res)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		result := Eval(fn.Body, extendedEnv)
		if ret, ok := result.(*object.Return); ok {
			return ret.Value
		}
		return result
	default:
		return object.NewError(fmt.Sprintf("unknown function"))
	}
	return None
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	// primero creamos un nuevo environment
	env := object.NewEnclosedEnv(fn.Env)

	for idx, param := range fn.Parameters {
		env.Set(param.Value.(string), 'l', args[idx])
	}

	return env
}
