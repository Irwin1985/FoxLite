"""
Un AST del inglés (Abstract Syntax Tree) es un Árbol de Sintaxis Abstracta cuyo objetivo es representar la
semántica del programa a evaluar.
"""


class AST:
    pass


class Program(AST):
    def __init__(self):
        self.statements = []

    def __str__(self):
        output = ''
        for statement in self.statements:
            output += repr(statement)

        return output

    __repr__ = __str__


class Identifier(AST):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return self.value

    __repr__ = __str__


class Integer(AST):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return str(self.value)

    __repr__ = __str__


class String(AST):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return self.value

    __repr__ = __str__


class Boolean(AST):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return '.t.' if self.value else '.f.'

    __repr__ = __str__


class Null(AST):
    def __init__(self):
        self.value = None

    def __str__(self):
        return '.null.'

    __repr__ = __str__


class Block(AST):
    def __init__(self):
        self.statements = []

    def __str__(self):
        output = ''
        for statement in self.statements:
            output += repr(statement)

        return output

    __repr__ = __str__


class FunctionDecl(AST):
    def __init__(self, name, params=None, body=None):
        self.name = name
        self.params = params if params is not None else []
        self.body = body

    def __str__(self):
        output = 'function'
        output += ' ' + repr(self.name)
        output += '('

        params = []
        for param in self.params:
            params.append(repr(param))

        output += ','.join(params)
        output += ')'

        output += '\n'

        output += '\t' + repr(self.body)

        output += '\n'
        output += 'endfunc'

        return output

    __repr__ = __str__


class FunctionCall(AST):
    def __init__(self):
        self.name = None
        self.arguments = []

    def __str__(self):
        output = repr(self.name)
        output += '('
        args = []

        for arg in self.arguments:
            args.append(repr(arg))

        output += ','.join(args)
        output += ')'

        return output

    __repr__ = __str__


class DoWhile(AST):
    def __init__(self, condition=None, block=None):
        self.condition = condition
        self.block = block

    def __str__(self):
        output = 'do while '
        output += repr(self.condition)
        output += '\n'
        output += '\t' + repr(self.block)
        output += '\n'
        output += 'enddo'
        return output

    __repr__ = __str__


class IfStatement(AST):
    def __init__(self):
        self.condition = None
        self.consequence = None
        self.alternative = None

    def __str__(self):
        output = 'if '
        output += repr(self.condition)
        output += '\n'
        output += '\t' + repr(self.consequence)

        if len(self.alternative.statements) > 0:
            output += '\n'
            output += 'else'
            output += '\n'
            output += '\t' + repr(self.alternative)

        output += '\n'
        output += 'endif'
        return output

    __repr__ = __str__


class BinaryOp(AST):
    def __init__(self, left, operator, right):
        self.left = left
        self.operator = operator
        self.right = right

    def __str__(self):
        output = repr(self.left)
        output += ' '
        output += self.operator
        output += ' '
        output += repr(self.right)
        return output

    __repr__ = __str__


class UnaryOp(AST):
    def __init__(self, operator, right):
        self.operator = operator
        self.right = right

    def __str__(self):
        output = self.operator
        output += repr(self.right)
        return output

    __repr__ = __str__


class VariableDecl(AST):
    def __init__(self, name, scope):
        self.name = name
        self.scope = scope

    def __str__(self):
        output = self.scope
        output += ' ' + repr(self.name)
        return output

    __repr__ = __str__


class Assignment(AST):
    def __init__(self, name, value=None):
        self.name = name
        self.value = value

    def __str__(self):
        output = repr(self.name)
        output += ' = '
        output += repr(self.value)
        return output

    __repr__ = __str__


class ReturnStmt(AST):
    def __init__(self, value=None):
        self.value = value

    def __str__(self):
        output = 'return '
        output += repr(self.value)
        return output

    __repr__ = __str__


class PrintStmt(AST):
    def __init__(self):
        self.arguments = []

    def __str__(self):
        output = '?'
        args = []

        for arg in self.arguments:
            args.append(repr(arg))

        output += ','.join(args)
        return output

    __repr__ = __str__