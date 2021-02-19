"""
    Evaluador del programa.
"""

import src.book.fox_lite_ast as ast
import src.book.fox_lite_object as obj


TRUE = obj.Boolean(value=True)
FALSE = obj.Boolean(value=False)
NULL = obj.Null()


class Evaluator:
    def eval(self, node, env):
        # Evaluar el programa principal
        if type(node) is ast.Program:
            return self.eval_program(node, env)
        # Evaluar bloques de sentencias
        elif type(node) is ast.Block:
            self.eval_block(node, env)
        # Evaluar enteros
        elif type(node) is ast.Integer:
            return obj.Integer(value=node.value)
        # Evaluar String
        elif type(node) is ast.String:
            return obj.String(value=node.value)
        # Evaluar booleanos
        elif type(node) is ast.Boolean:
            return TRUE if node.value else FALSE
        # Evaluar Null
        elif type(node) is ast.Null:
            return NULL
        # Evaluar expresiones unarias
        elif type(node) is ast.UnaryOp:
            operator = node.operator
            right = self.eval(node.right, env)

            if self.is_error(right):
                return right

            return self.eval_unary_expression(operator, right)
        # Evaluar expresiones binarias
        elif type(node) is ast.BinaryOp:
            operator = node.operator

            left = self.eval(node.left, env)
            if self.is_error(left):
                return left

            right = self.eval(node.right, env)
            if self.is_error(right):
                return right

            return self.eval_binary_expression(left, operator, right)
        # Evaluar declaraciónes de Variables
        elif type(node) is ast.VariableDecl:
            #  Las variables por defecto se declaran en .F.
            return env.set(node.token.value, FALSE, node.scope)
        # Evaluar asignaciones de Variables
        elif type(node) is ast.Assignment:
            # Resolvemos su valor antes de guardarlo en la Tabla de Símbolos.
            val = self.eval(node.value, env)
            if self.is_error(val):
                return val

            return env.set(node.token.value, val)
        # Evaluar identificadores
        elif type(node) is ast.Identifier:
            return self.eval_identifier(node, env)

    def eval_program(self, program, env):
        result = None

        for statement in program.statements:
            result = self.eval(statement, env)

            if result is not None:
                if result.type() == obj.Type.RETURN:
                    return result.value  # Se retorna la expresión del return.
                elif result.type() == obj.Type.ERROR:
                    return result  # Se retorna el objeto error.

        return result

    def eval_block(self, block, env):
        result = None

        for statement in block.statements:
            result = self.eval(statement, env)
            if self.is_error(result):
                return result

            if result is not None and result.type() in (obj.Type.ERROR, obj.Type.RETURN):
                return result

        return result

    def eval_identifier(self, node, env):
        val = env.get(node.value)
        if val is None:
            return self.new_error(f"variable '{node.value}' no definida.")

        return val

    def eval_unary_expression(self, operator, right):
        if operator == "!":
            if right.type() != obj.Type.BOOLEAN:
                self.new_error("tipo de dato incompatible. Se esperaba un BOOLEAN")
            elif right == TRUE:
                return FALSE
            elif right == FALSE:
                return TRUE
        elif operator == "-":
            if right.type() != obj.Type.INTEGER:
                self.new_error("tipo de dato incompatible. Se esperaba un INTEGER")
            else:
                return obj.Integer(value=-right.value)

    def eval_binary_expression(self, left, operator, right):
        if left.type() == obj.Type.INTEGER and right.type() == obj.Type.INTEGER:
            if operator in ('+', '-', '*', '/'):
                return self.eval_native_integer_to_integer_object(left, operator, right)
            elif operator in ('<', '>', '<=', '>=', '==', '!='):
                return self.eval_native_integer_to_boolean_object(left, operator, right)
            else:
                return self.new_error(f"operador no soportado para el tipo INTEGER: '{operator}'")
        elif left.type() == obj.Type.STRING and right.type() == obj.Type.STRING:
            if operator == '+':
                return obj.String(value=left.value + right.value)
            else:
                return self.new_error(f"operador no soportado para el tipo STRING: '{operator}'")
        elif left.type() == obj.Type.BOOLEAN and right.type() == obj.Type.BOOLEAN:
            if operator in ('<', '>', '<=', '>=', '==', '!=', 'and', 'or'):
                return self.eval_native_boolean_to_boolean_object(left, operator, right)
            else:
                return self.new_error(f"operador no soportado para el tipo BOOLEAN: '{operator}'")

        elif left.type() != right.type():
            return self.new_error(f'incompatibilidad de tipos: {left.type()}, {right.type()}')
    """
        Evalúa los operandos integer nativo a objeto boolean.
        Las operaciónes realizadas son aritméticas.
    """
    def eval_native_integer_to_integer_object(self, left, operator, right):
        left_val = left.value
        right_val = right.value

        if operator == '+':
            return obj.Integer(value=left_val + right_val)
        elif operator == '-':
            return obj.Integer(value=left_val - right_val)
        elif operator == '*':
            return obj.Integer(value=left_val * right_val)
        elif operator == '/':
            if right_val == 0:
                return self.new_error("división por cero.")
            return obj.Integer(value=left_val / right_val)
        else:
            return self.new_error(f"operador desconocido: '{operator}'")

    """
        Evalúa los operandos integer nativo a objeto boolean.
        Las operaciónes realizadas son relacionales.
    """
    def eval_native_integer_to_boolean_object(self, left, operator, right):
        left_val = left.value
        right_val = right.value

        if operator == '<':
            return TRUE if left_val < right_val else FALSE
        elif operator == '<=':
            return TRUE if left_val <= right_val else FALSE
        elif operator == '>':
            return TRUE if left_val > right_val else FALSE
        elif operator == '>=':
            return TRUE if left_val >= right_val else FALSE
        elif operator == '==':
            return TRUE if left_val == right_val else FALSE
        elif operator == '!=':
            return TRUE if left_val != right_val else FALSE
        else:
            return self.new_error(f"operador desconocido: '{operator}'")

    """
        Evalúa los operandos boolean nativo a objeto boolean.
        Las operaciónes realizadas son lógicas.
    """
    def eval_native_boolean_to_boolean_object(self, left, operator, right):
        left_val = left.value
        right_val = right.value

        if operator == 'and':
            return TRUE if left_val and right_val else FALSE
        elif operator == 'or':
            if left_val:
                return TRUE  # No hace falta evaluar el segundo operando.
            elif right_val:
                return TRUE
            else:
                return FALSE
        else:
            left_val = obj.Integer(value=1 if left_val else 0)
            right_val = obj.Integer(value=1 if right_val else 0)
            return self.eval_native_integer_to_boolean_object(left_val, operator, right_val)

    def new_error(self, message):
        return obj.Error(message=message)

    def is_error(self, node):
        return node is not None and node.type() == obj.Type.ERROR
