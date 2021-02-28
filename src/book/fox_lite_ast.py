"""
Un AST del inglés (Abstract Syntax Tree) es un Árbol de Sintaxis Abstracta cuyo objetivo es representar la
semántica del programa a evaluar.
"""


class AST:
    pass


class Program(AST):
    def __init__(self):
        self.statements = []


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


class FunctionDecl(AST):
    def __init__(self, name, params=None, body=None):
        self.name = name
        self.params = params if params is not None else []
        self.body = body


class FunctionCall(AST):
    def __init__(self):
        self.name = None
        self.arguments = []


class DoWhile(AST):
    def __init__(self, condition=None, block=None):
        self.condition = condition
        self.block = block


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


class UnaryOp(AST):
    def __init__(self, operator, right):
        self.operator = operator
        self.right = right


class VariableDecl(AST):
    def __init__(self, name, scope):
        self.name = name
        self.scope = scope


class Assignment(AST):
    def __init__(self, name, value=None):
        self.name = name
        self.value = value


class ReturnStmt(AST):
    def __init__(self, value=None):
        self.value = value


class PrintStmt(AST):
    def __init__(self):
        self.arguments = []
