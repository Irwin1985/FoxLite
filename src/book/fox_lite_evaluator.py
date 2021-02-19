"""
    Evaluador del programa.
"""

import src.book.fox_lite_ast as ast
import src.book.fox_lite_object as obj
import src.book.fox_lite_environment as environment
from src.book.fox_lite_builtins import builtins


TRUE = obj.Boolean(value=True)
FALSE = obj.Boolean(value=False)
NULL = obj.Null()


def is_error(node):
    """
    Verifica si el objeto resultante es de tipo obj.Error()
    :param node:
    :return:
    """
    return node is not None and node.type() == obj.Type.ERROR


def new_error(message):
    """
    Retorna un objeto de tipo obj.Error() con la información del error.
    :param message:
    :return:
    """
    return obj.Error(message=message)


def eval_native_integer_to_boolean_object(left, operator, right):
    """
    Evalúa los operandos integer nativo a objeto boolean.
    Las operaciónes realizadas son relacionales.
    :param left:
    :param operator:
    :param right:
    :return:
    """
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
        return new_error(f"operador desconocido: '{operator}'")


def eval_native_boolean_to_boolean_object(left, operator, right):
    """
    Evalúa los operandos boolean nativo a objeto boolean.
    Las operaciones realizadas son lógicas.
    :param left:
    :param operator:
    :param right:
    :return:
    """
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
        return eval_native_integer_to_boolean_object(left_val, operator, right_val)


def eval_native_integer_to_integer_object(left, operator, right):
    """
    Evalúa los operandos integer nativo a objeto integer.
    Las operaciones realizadas son aritméticas.
    :param left:
    :param operator:
    :param right:
    :return:
    """
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
            return new_error("división por cero.")
        return obj.Integer(value=left_val / right_val)
    else:
        return new_error(f"operador desconocido: '{operator}'")


def eval_binary_expression(left, operator, right):
    """
    Realiza operaciones binarias y devuelve el objeto correspondiente.
    :param left:
    :param operator:
    :param right:
    :return:
    """
    if left.type() == obj.Type.INTEGER and right.type() == obj.Type.INTEGER:
        if operator in ('+', '-', '*', '/'):
            return eval_native_integer_to_integer_object(left, operator, right)
        elif operator in ('<', '>', '<=', '>=', '==', '!='):
            return eval_native_integer_to_boolean_object(left, operator, right)
        else:
            return new_error(f"operador no soportado para el tipo INTEGER: '{operator}'")
    elif left.type() == obj.Type.STRING and right.type() == obj.Type.STRING:
        if operator == '+':
            return obj.String(value=left.value + right.value)
        else:
            return new_error(f"operador no soportado para el tipo STRING: '{operator}'")
    elif left.type() == obj.Type.BOOLEAN and right.type() == obj.Type.BOOLEAN:
        if operator in ('<', '>', '<=', '>=', '==', '!=', 'and', 'or'):
            return eval_native_boolean_to_boolean_object(left, operator, right)
        else:
            return new_error(f"operador no soportado para el tipo BOOLEAN: '{operator}'")

    elif left.type() != right.type():
        return new_error(f'incompatibilidad de tipos: {left.type()}, {right.type()}')


def eval_unary_expression(operator, right):
    """
    Realiza las operaciones unarias y devuelve el objeto resultante.
    :param operator:
    :param right:
    :return:
    """
    if operator == "!":
        if right.type() != obj.Type.BOOLEAN:
            new_error("tipo de dato incompatible. Se esperaba un BOOLEAN")
        elif right == TRUE:
            return FALSE
        elif right == FALSE:
            return TRUE
    elif operator == "-":
        if right.type() != obj.Type.INTEGER:
            new_error("tipo de dato incompatible. Se esperaba un INTEGER")
        else:
            return obj.Integer(value=-right.value)


def eval_identifier(node, env):
    """
    Resuelve el valor de un identificador.
    :param node:
    :param env:
    :return:
    """
    val = env.get(node.value)
    if val is None:
        return new_error(f"variable '{node.value}' no definida.")

    return val


class Evaluator:
    def eval(self, node, env):
        # Evaluar el programa principal
        if type(node) is ast.Program:
            return self.eval_program(node, env)
        # Evaluar bloques de sentencias
        elif type(node) is ast.Block:
            return self.eval_block(node, env)
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

            if is_error(right):
                return right

            return eval_unary_expression(operator, right)
        # Evaluar expresiones binarias
        elif type(node) is ast.BinaryOp:
            operator = node.operator

            left = self.eval(node.left, env)
            if is_error(left):
                return left

            right = self.eval(node.right, env)
            if is_error(right):
                return right

            return eval_binary_expression(left, operator, right)
        # Evaluar declaraciónes de Variables
        elif type(node) is ast.VariableDecl:
            #  Las variables por defecto se declaran en .F.
            return env.set(node.token.value, FALSE, node.scope)
        # Evaluar asignaciones de Variables
        elif type(node) is ast.Assignment:
            # Resolvemos su valor antes de guardarlo en la Tabla de Símbolos.
            val = self.eval(node.value, env)
            if is_error(val):
                return val

            return env.set(node.token.value, val)
        # Evaluar identificadores
        elif type(node) is ast.Identifier:
            return eval_identifier(node, env)
        # Sentencia If
        elif type(node) is ast.IfStatement:
            return self.eval_if_statement(node, env)
        # Declaración de Función
        elif type(node) is ast.FunctionDecl:
            # Creamos el objeto Function
            function = obj.Function(
                name=node.name.value,
                params=node.params,
                body=node.body,
                env=env,
            )
            env.set(function.name, function)
            return function
        # Llamada a Función
        elif type(node) is ast.FunctionCall:
            function = env.get(node.name.value)

            if function is None:
                # Intentamos buscar la función como builtin
                function = builtins.get(node.name.value)
                if function is None:
                    return new_error(f"función no definida: '{node.name.value}'.")

            args = self.eval_arguments(node.arguments, env)

            if len(args) == 1 and is_error(args[0]):
                return args[0]

            return self.execute_function(function, args)
        # Do While
        elif type(node) is ast.DoWhile:
            result = None
            while True:
                condition = self.eval(node.condition, env)

                if is_error(condition):
                    return condition

                if condition == FALSE:
                    break

                result = self.eval(node.block, env)

            return result
        # Sentencia Return
        elif type(node) is ast.ReturnStmt:
            return_value = self.eval(node.value, env)

            if is_error(return_value):
                return return_value

            return obj.Return(value=return_value)

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
            if is_error(result):
                return result

            if result is not None and result.type() in (obj.Type.ERROR, obj.Type.RETURN):
                return result

        return result

    def eval_if_statement(self, node, env):
        condition = self.eval(node.condition, env)

        if is_error(condition):
            return condition

        if condition == TRUE:
            return self.eval(node.consequence, env)
        elif node.alternative is not None:
            return self.eval(node.alternative, env)

    def execute_function(self, func, args):
        if type(func) is obj.Function:
            extended_env = extend_function_env(func, args)
            # Ejecutamos el cuerpo de la función
            evaluated = self.eval(func.body, extended_env)
            if type(evaluated) is obj.Return:
                return evaluated.value

            return evaluated
        elif type(func) is obj.Builtin:
            return func.function(args)
        else:
            return new_error(f"{func.type()} no es una función.")

    def eval_arguments(self, args, env):
        result = []
        for argument in args:
            evaluated = self.eval(argument, env)
            if is_error(evaluated):
                return [evaluated]

            result.append(evaluated)

        return result


def extend_function_env(func, args):
    extended_env = environment.Environment(parent=func.env)
    # Registramos los símbolos de la función
    for param_idx, param in enumerate(func.params):
        extended_env.set(
            name=param.value,
            val=args[param_idx],
            scope='default',
            current_scope=True)

    return extended_env
