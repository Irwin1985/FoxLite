
"""
Un AST del inglés (Abstract Syntax Tree) es un Árbol de Sintaxis Abstracta cuyo objetivo es representar la
semántica del programa a evaluar.
"""


class AST:
    def string(self):
        pass


class Program(AST):
    def __init__(self):
        self.statements = []

    def string(self):
        out = ""
        if len(self.statements) > 0:
            for statement in self.statements:
                out += statement.string()

        return out


class Identifier(AST):
    def __init__(self, value):
        self.value = value


class Integer(AST):
    def __init__(self, value):
        self.value = value


class String(AST):
    def __init__(self, value):
        self.value = value


class Boolean(AST):
    def __init__(self, value):
        self.value = value


class Null(AST):
    def __init__(self):
        self.value = None


class Block(AST):
    def __init__(self):
        self.statements = []

    def string(self):
        out = ""
        if len(self.statements) > 0:
            for statement in self.statements:
                out += statement.string()

        return out


class FunctionDecl(AST):
    def __init__(self, name, params=None, body=None):
        self.name = name
        self.params = params if params is not None else []
        self.body = body

    def string(self):
        out = "function " + self.name
        out += "("
        param_list = ""
        if len(self.params) > 0:
            for param in self.params:
                param_list += param.string()
            out += ",".join(param_list)
        out += ")"

        return out


class DoWhile(AST):
    def __init__(self, condition=None, block=None):
        self.condition = condition
        self.block = block


class FunctionCall(AST):
    def __init__(self):
        self.name = None
        self.arguments = []

    def string(self):
        pass


class IfStatement(AST):
    def __init__(self):
        self.condition = None
        self.consequence = None
        self.alternative = None


class BinaryOp(AST):
    def __init__(self, left, operator, right):
        self.left = left
        self.operator = operator
        self.right = right

    def string(self):
        pass


class UnaryOp(AST):
    def __init__(self, operator, right):
        self.operator = operator
        self.right = right

    def string(self):
        pass


class VariableDecl(AST):
    def __init__(self, token, scope):
        self.token = token
        self.scope = scope


class Assignment(AST):
    def __init__(self, token, value=None):
        self.token = token
        self.value = value


class ReturnStmt(AST):
    def __init__(self, value=None):
        self.value = value


class PrintStmt(AST):
    def __init__(self):
        self.arguments = []
