package evaluator

import (
	"FoxLite/src/ast"
	"FoxLite/src/object"
	"fmt"
)

func evalWhileStmt(node *ast.While, env *object.Environment) object.Object {
	for { // Un While se ejecuta dentro de un bucle infinito
		cond := Eval(node.Condition, env) // Primero evaluamos la condición del While
		if isError(cond) {
			return cond
		}
		if cond.Type() != object.BooleanObj { // Validamos que el resultado sea 'Bool'
			return object.NewError(fmt.Sprintf("non-bool type `%v` used as case condition", cond.Type()))
		}
		if cond.(*object.Boolean).Value { // Siempre y cuando sea Verdadero, ejecutamos el bloque.
			var action byte
			var res object.Object
			// Evaluamos el bloque del While
			for _, stmt := range node.Body.Statements {
				res = Eval(stmt, env)
				if isError(res) {
					return res
				}

				rType := res.Type()
				if rType == object.ExitObj {
					action = 'b' // break
					break        // rompemos el bucle
				} else if rType == object.LoopObj {
					action = 'l' // loop
					break        // rompemos y pasamos al siguiente loop
				} else if rType == object.ReturnObj {
					action = 'r' // return
					break        // rompemos y retornamos res
				}
			}

			// evaluamos el estado de la ejecución del loop actual
			if action == 'b' {
				break // rompemos el bucle
			} else if action == 'l' {
				continue // pasamos al siguiente bucle
			} else if action == 'r' {
				return res // retornamos res
			}
		} else {
			// Rompemos el Loop porque la condición es False
			break
		}
	}
	return None
}
