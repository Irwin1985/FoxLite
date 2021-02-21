import src.book.fox_lite_ast as ast
import src.book.fox_lite_object as obj
import src.book.fox_lite_environment as environment
from src.book.fox_lite_builtins import builtins


TRUE = obj.Boolean(value=True)
FALSE = obj.Boolean(value=False)
NULL = obj.Null()


def is_error(value_obj):
    """
    Verifica si el objeto dato es de tipo obj.Error()
    :param value_obj:
    :return: obj
    """
    return value_obj is not None and value_obj.type() == obj.Type.ERROR


def is_relational_operator(operator):
    """
    Determina si el operador dado es de origen relacional.
    :param operator:
    :return: Boolean
    """
    return operator in ('<', '>', '<=', '>=', '==', '!=')


def new_error(message):
    """
    Retorna un objeto de tipo obj.Error() con la información del error.
    :param message:
    :return: obj.Error()
    """
    return obj.Error(message=message)


def eval_native_integer_to_boolean_object(left_obj, operator, right_obj):
    """
    Evalúa los operandos integer nativo a objeto boolean.
    Las operaciónes realizadas son relacionales.
    :param left_obj:
    :param operator:
    :param right_obj:
    :return: obj.Boolean()
    """
    left_val = left_obj.value
    right_val = right_obj.value

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


def eval_native_boolean_to_boolean_object(left_obj, operator, right_obj):
    """
    Evalúa los operandos boolean nativo a objeto boolean.
    Las operaciones realizadas son relacionales.
    Como los operados son de tipo Boolean, primero se tienen que
    convertir a su equivalente numérico. (.f. => 0, .t. => 1)
    :param left_obj:
    :param operator:
    :param right_obj:
    :return: obj.Boolean()
    """
    left_val = left_obj.value
    right_val = right_obj.value

    left_obj = obj.Integer(value=1 if left_val else 0)
    right_obj = obj.Integer(value=1 if right_val else 0)

    return eval_native_integer_to_boolean_object(left_obj, operator, right_obj)


def eval_native_integer_to_integer_object(left_obj, operator, right_obj):
    """
    Evalúa los operandos integer nativo a objeto integer.
    Las operaciones realizadas son aritméticas.
    :param left_obj:
    :param operator:
    :param right_obj:
    :return: obj.Integer()
    """
    left_val = left_obj.value
    right_val = right_obj.value

    if operator == '+':
        return obj.Integer(value=left_val + right_val)
    elif operator == '-':
        return obj.Integer(value=left_val - right_val)
    elif operator == '*':
        return obj.Integer(value=left_val * right_val)
    elif operator == '/':
        if right_val == 0:
            return new_error("División por cero.")
        return obj.Integer(value=left_val / right_val)


def eval_native_string_to_object(left_obj, operator, right_obj):
    """
    Evalúa los operandos string nativo a objeto string.
    Las operaciones realizadas son aritméticas.
    :param left_obj:
    :param operator:
    :param right_obj:
    :return: object
    """
    left_val = left_obj.value
    right_val = right_obj.value

    if operator == '+':
        return obj.String(value=left_val + right_val)
    elif operator == '==':
        return obj.Boolean(value=left_val == right_val)
    elif operator == '!=':
        return obj.Boolean(value=left_val != right_val)


def eval_binary_expression(left_obj, operator, right_obj):
    """
    Realiza operaciones binarias y devuelve el objeto correspondiente.
    :param left_obj:
    :param operator:
    :param right_obj:
    :return: object
    """
    if left_obj.type() == obj.Type.INTEGER and right_obj.type() == obj.Type.INTEGER:
        if operator in ('+', '-', '*', '/'):
            return eval_native_integer_to_integer_object(left_obj, operator, right_obj)
        elif is_relational_operator(operator):
            return eval_native_integer_to_boolean_object(left_obj, operator, right_obj)
        else:
            return new_error(f"Operador no soportado para el tipo INTEGER: '{operator}'")

    elif left_obj.type() == obj.Type.STRING and right_obj.type() == obj.Type.STRING:
        if operator in ('+', '==', '!='):
            return eval_native_string_to_object(left_obj, operator, right_obj)
        else:
            return new_error(f"Operador no soportado para el tipo STRING: '{operator}'")

    elif left_obj.type() == obj.Type.BOOLEAN and right_obj.type() == obj.Type.BOOLEAN:
        if is_relational_operator(operator):
            return eval_native_boolean_to_boolean_object(left_obj, operator, right_obj)
        else:
            return new_error(f"Operador no soportado para el tipo BOOLEAN: '{operator}'")

    elif left_obj.type() != right_obj.type():
        return new_error(f'Incompatibilidad de tipos: {left_obj.type()}, {right_obj.type()}')


def eval_unary_expression(operator, right_obj):
    """
    Realiza las operaciones unarias y devuelve el objeto resultante.
    :param operator:
    :param right_obj:
    :return: object
    """
    if operator == "!":
        if right_obj.type() != obj.Type.BOOLEAN:
            return new_error("Tipo de dato incompatible. Se esperaba un BOOLEAN")
        elif right_obj == TRUE:
            return FALSE
        elif right_obj == FALSE:
            return TRUE
    elif operator == "-":
        if right_obj.type() != obj.Type.INTEGER:
            return new_error("Tipo de dato incompatible. Se esperaba un INTEGER")
        else:
            return obj.Integer(value=-right_obj.value)


def eval_identifier(ast_node, env):
    """
    Evalúa el valor de un identificador.
    :param ast_node:
    :param env:
    :return: object
    """
    val_obj = env.get(ast_node.value)
    if val_obj is None:
        return new_error(f"Variable '{ast_node.value}' no definida.")

    return val_obj


def extend_function_env(func_obj, arguments):
    """
    Crea un nuevo environment para la función dada.
    :param func_obj:
    :param arguments:
    :return: Environment
    """
    extended_env = environment.Environment(parent=func_obj.env)
    # Registramos los símbolos de la función
    for param_idx, param in enumerate(func_obj.params):
        extended_env.set(
            name=param.value,
            val=arguments[param_idx],
            scope='default',
            current_scope=True)

    return extended_env


class Evaluator:
    def eval(self, ast_node, env):
        # Evaluar el programa principal
        if type(ast_node) is ast.Program:
            return self.eval_program(ast_node, env)
        # Evaluar bloques de sentencias
        elif type(ast_node) is ast.Block:
            return self.eval_block(ast_node, env)
        # Evaluar Integer
        elif type(ast_node) is ast.Integer:
            return obj.Integer(value=ast_node.value)
        # Evaluar String
        elif type(ast_node) is ast.String:
            return obj.String(value=ast_node.value)
        # Evaluar Boolean
        elif type(ast_node) is ast.Boolean:
            return TRUE if ast_node.value else FALSE
        # Evaluar Null
        elif type(ast_node) is ast.Null:
            return NULL
        # Evaluar expresiones unarias
        elif type(ast_node) is ast.UnaryOp:
            operator = ast_node.operator
            right_obj = self.eval(ast_node.right, env)

            if is_error(right_obj):
                return right_obj

            return eval_unary_expression(operator, right_obj)
        # Evaluar expresiones binarias
        elif type(ast_node) is ast.BinaryOp:
            operator = ast_node.operator

            if operator in ('or', 'and'):
                # Invocamos directamente sin evaluar los operandos.
                return self.eval_logical_expression(ast_node, env)

            left_obj = self.eval(ast_node.left, env)
            if is_error(left_obj):
                return left_obj

            right_obj = self.eval(ast_node.right, env)
            if is_error(right_obj):
                return right_obj

            return eval_binary_expression(left_obj, operator, right_obj)
        # Evaluar declaraciónes de Variables
        elif type(ast_node) is ast.VariableDecl:
            #  Las variables por defecto se declaran en .F.
            return env.set(ast_node.token.value, FALSE, ast_node.scope)
        # Evaluar asignaciones de Variables
        elif type(ast_node) is ast.Assignment:
            # Resolvemos su valor antes de guardarlo en la Tabla de Símbolos.
            val_obj = self.eval(ast_node.value, env)
            if is_error(val_obj):
                return val_obj

            return env.set(ast_node.token.value, val_obj)
        # Evaluar identificadores
        elif type(ast_node) is ast.Identifier:
            return eval_identifier(ast_node, env)
        # Sentencia If
        elif type(ast_node) is ast.IfStatement:
            return self.eval_if_statement(ast_node, env)
        # Declaración de Función
        elif type(ast_node) is ast.FunctionDecl:
            # Creamos el objeto Function
            function = obj.Function(
                name=ast_node.name.value,
                params=ast_node.params,
                body=ast_node.body,
                env=env,
            )
            env.set(function.name, function)
            return function
        # Llamada a Función
        elif type(ast_node) is ast.FunctionCall:
            function = env.get(ast_node.name.value)

            if function is None:
                # Intentamos buscar la función como builtin
                function = builtins.get(ast_node.name.value)
                if function is None:
                    return new_error(f"Función no definida: '{ast_node.name.value}'.")

            arguments = self.eval_arguments(ast_node.arguments, env)

            if len(arguments) == 1 and is_error(arguments[0]):
                return arguments[0]

            return self.execute_function(function, arguments)
        # Do While
        elif type(ast_node) is ast.DoWhile:
            while True:
                condition_obj = self.eval(ast_node.condition, env)

                if is_error(condition_obj):
                    return condition_obj

                if condition_obj == FALSE:
                    break

                result_obj = self.eval(ast_node.block, env)
                if result_obj is not None:
                    if result_obj.type() == obj.Type.RETURN:
                        return result_obj.value  # Se retorna la expresión del return.
                    elif result_obj.type() == obj.Type.ERROR:
                        return result_obj  # Se retorna el objeto obj.Error().

        # Sentencia Return
        elif type(ast_node) is ast.ReturnStmt:
            return_obj = self.eval(ast_node.value, env)

            if is_error(return_obj):
                return return_obj

            return obj.Return(value=return_obj)
        # Sentencia Print
        elif type(ast_node) is ast.PrintStmt:
            for argument in ast_node.arguments:
                result_obj = self.eval(argument, env)
                if is_error(result_obj):
                    return result_obj
                print(result_obj.to_string(), end=" ")
            print()  # Empty line

    def eval_program(self, program, env):
        result_obj = None

        for statement in program.statements:
            result_obj = self.eval(statement, env)

            if result_obj is not None:
                if result_obj.type() == obj.Type.RETURN:
                    return result_obj.value  # Se retorna la expresión del return.
                elif result_obj.type() == obj.Type.ERROR:
                    return result_obj  # Se retorna el objeto error.

        return result_obj

    def eval_block(self, block, env):
        result_obj = None

        for statement in block.statements:
            result_obj = self.eval(statement, env)
            if is_error(result_obj):
                return result_obj

            if result_obj is not None and result_obj.type() in (obj.Type.ERROR, obj.Type.RETURN):
                return result_obj

        return result_obj

    def eval_if_statement(self, ast_node, env):
        condition_obj = self.eval(ast_node.condition, env)

        if is_error(condition_obj):
            return condition_obj

        if condition_obj == TRUE:
            return self.eval(ast_node.consequence, env)
        elif ast_node.alternative is not None:
            return self.eval(ast_node.alternative, env)

    def execute_function(self, func, arguments):
        if type(func) is obj.Function:
            extended_env = extend_function_env(func, arguments)
            # Ejecutamos el cuerpo de la función
            evaluated_obj = self.eval(func.body, extended_env)
            if type(evaluated_obj) is obj.Return:
                return evaluated_obj.value

            return evaluated_obj
        elif type(func) is obj.Builtin:
            return func.function(arguments)
        else:
            return new_error(f"{func.type()} no es una función.")

    def eval_arguments(self, arguments, env):
        evaluated_arguments = []
        for argument in arguments:
            evaluated_obj = self.eval(argument, env)
            if is_error(evaluated_obj):
                return [evaluated_obj]

            evaluated_arguments.append(evaluated_obj)

        return evaluated_arguments

    def eval_logical_expression(self, ast_node, env):
        """
        Realiza operaciónes lógicas entre 2 booleanos.
        :param ast_node: ast.BinaryOp()
        :param env: Environment()
        :return: obj.Boolean()
        """
        # Evaluamos primero el operando de la izquierda.
        left_obj = self.eval(ast_node.left, env)
        if is_error(left_obj):
            return left_obj

        if left_obj.type() != obj.Type.BOOLEAN:
            return new_error(f"Los operadores AND y OR solo soportan operandos de tipo BOOLEAN.")

        if ast_node.operator == 'and':
            if left_obj == FALSE:
                return FALSE
        elif ast_node.operator == 'or':
            if left_obj == TRUE:
                return TRUE

        # La operación depende del operando de la derecha así que lo evaluamos.
        right_obj = self.eval(ast_node.right, env)
        if is_error(right_obj):
            return right_obj

        if right_obj.type() != obj.Type.BOOLEAN:
            return new_error(f"Los operadores AND y OR solo soportan operandos de tipo BOOLEAN.")

        return TRUE if right_obj == TRUE else FALSE
