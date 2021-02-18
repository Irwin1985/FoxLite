"""
   Clase Token: representa la unidad más pequeña en la construcción del lenguaje. Será utilizado por el lexer
   para descomponer el código fuente y por el parser para validar la gramática y crear el AST.
"""

from enum import Enum


class TokenType(Enum):

    # Fin del código fuente
    EOF = 'EOF'

    # Nueva línea
    LBREAK = 'LBREAK'

    # Identificadores y Literales
    IDENT = 'IDENT'  # add, foobar, x, y, ...
    INT = 'INT'  # 1343456

    # Operadores aritméticos
    ASSIGN = '='
    PLUS = '+'
    MINUS = '-'
    MUL = '*'
    DIV = '/'

    # Operadores relacionales
    NOT = '!'
    LESS = '<'
    GREATER = '>'
    LESS_EQ = '<='
    GREATER_EQ = '>='
    EQUAL = '=='
    NOT_EQUAL = '!='

    # Operadores Lógicos
    AND = 'and'
    OR = 'or'

    # Caracteres especiales
    COMMA = ','
    DOT = '.'
    LPAREN = '('
    RPAREN = ')'
    AMPER = '&'
    PRINT = '?'
    STRING = 'STRING'

    # Funciones
    FUNCTION = 'function'
    ENDFUNC = 'endfunc'
    MESSAGEBOX = 'messagebox'

    # Do While
    DO = 'do'
    WHILE = 'while'
    ENDDO = 'enddo'

    # If / EndIf
    IF = 'if'
    ELSE = 'else'
    ENDIF = 'endif'
    RETURN = 'return'

    # Declaración de Variables
    PUBLIC = 'public'
    LOCAL = 'local'
    PRIVATE = 'private'

    # Boolean
    TRUE = 'true'
    FALSE = 'false'
    NULL = 'null'


class Token:
    def __init__(self, token_type, token_value):
        self.type = token_type
        self.value = token_value

    def __str__(self):
        return f"type: {self.type}, value: '{self.value}'"

    __repr__ = __str__


# Diccionario de Palabras Reservadas
keywords = {
    "function": TokenType.FUNCTION,
    "endfunc": TokenType.ENDFUNC,
    "messagebox": TokenType.MESSAGEBOX,
    "do": TokenType.DO,
    "while": TokenType.WHILE,
    "enddo": TokenType.ENDDO,
    "if": TokenType.IF,
    "else": TokenType.ELSE,
    "endif": TokenType.ENDIF,
    "return": TokenType.RETURN,
    "public": TokenType.PUBLIC,
    "local": TokenType.LOCAL,
    "private": TokenType.PRIVATE,
    ".t.": TokenType.TRUE,
    ".f.": TokenType.FALSE,
    ".null.": TokenType.NULL,
    "and": TokenType.AND,
    "or": TokenType.OR,
}


def lookup_ident(ident):
    tok = keywords.get(ident)
    return tok if tok is not None else TokenType.IDENT
