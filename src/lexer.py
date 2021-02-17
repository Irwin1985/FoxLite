"""
  Lexer: también conocido como tokenizer, se encarga de leer y descomponer el código fuente en pequeñas unidades
  llamadas tokens.
"""

from src.tokens import Token, TokenType, lookup_ident


class Lexer:
    def __init__(self, source_code):
        self.source = source_code
        self.current_pos = 0
        self.current_char = self.source[self.current_pos]
        self.last_token_value = None

    # Mostramos un error cuando nos encontremos un caracter extraño.
    def error(self):
        raise Exception(f'Error: caractér desconocido {self.current_char}')

    # Crea un token y actualiza el último token generado.
    def new_token(self, token_type, token_value):
        token = Token(type=token_value, value=token_value)
        self.last_token_value = token_value
        return token

    # Avanza un caracter en el código fuente
    def advance(self):
        self.current_pos += 1
        if self.current_pos > len(self.source) - 1:
            self.current_char = None
        else:
            self.current_char = self.source[self.current_pos]

    # Mira 1 caracter hacia adelante en el código fuente
    def peek(self):
        peek_position = self.current_pos + 1
        if peek_position > len(self.source) - 1:
            return None
        else:
            return self.source[peek_position]

    # Ignoramos los comentarios
    def ignore_comments(self):
        # Avanzamos los caracteres hasta el final de la línea
        while self.current_char is not None and self.current_char != '\n':
            self.advance()

    # Ignoramos los espacios en blanco
    def ignore_blanks(self):
        # Avanzamos los caracteres mientras hayan espacios o EOF
        while self.current_char is not None and is_space(self.current_char):
            self.advance()

    # Obtiene un entero
    def number(self):
        result = ''
        while self.current_char is not None and self.current_char.isdigit():
            result += self.current_char
            self.advance()

        return self.new_token(TokenType.INT, result)

    # Obtiene un identificador
    def identifier(self):
        result = ''
        while self.current_char is not None and is_letter(self.current_char):
            result += self.current_char
            self.advance()

        token_type = lookup_ident(result)
        return self.new_token(token_type, result)

    # Obtiene una palabra reservada de entre 2 puntos ej: .t., .f., .null.
    def dotted_indentifier(self):
        self.advance()  # Avanza el primer punto '.'
        token = self.identifier()
        self.advance()  # Avanza el segundo punto '.'

        return token

    # Obtiene una secuencia de caracteres (string)
    def string(self):
        self.advance()  # Avanza el primer '"'
        result = ''
        while self.current_char is not None and self.current_char != '"':
            result += self.current_char
            self.advance()

        self.advance()  # Avanzamos el segundo '"'
        return self.new_token(TokenType.STRING, result)

    # Extrae el siguiente token desde el código fuente
    def next_token(self):
        while self.current_char is not None:

            # Salto de Línea
            if self.current_char == '\n':
                self.advance()  # Avanza el '\n'
                # Si el último token es un salto de línea entonces ignoramos este.
                if self.last_token_value is not None and self.last_token_value != 'LBREAK':
                    return self.new_token(TokenType.LBREAK, 'LBREAK')

            # Ignoramos los espacios en blanco
            if self.current_char.isspace():
                self.ignore_blanks()
                continue

            # Ignoramos los comentarios
            if self.current_char == '&' and self.peek() == '&':
                self.ignore_comments()
                continue

            # Números (enteros)
            if self.current_char.isdigit():
                return self.number()

            # Identificadores
            if is_letter(self.current_char):
                return self.identifier()

            # Palabras reservadas encerradas en puntos
            if self.current_char == '.':
                return self.dotted_indentifier()

            # Secuencia de caracteres (string)
            if self.current_char == '"':
                return self.string()

            # Caracteres de 1 digito de longitud
            if self.current_char == '+':
                self.advance()
                return self.new_token(TokenType.PLUS, '+')

            if self.current_char == '-':
                self.advance()
                return self.new_token(TokenType.PLUS, '-')

            if self.current_char == '*':
                self.advance()
                return self.new_token(TokenType.PLUS, '*')

            if self.current_char == '/':
                self.advance()
                return self.new_token(TokenType.PLUS, '/')

            if self.current_char == '(':
                self.advance()
                return self.new_token(TokenType.LPAREN, '(')

            if self.current_char == ')':
                self.advance()
                return self.new_token(TokenType.RPAREN, ')')

            if self.current_char == ',':
                self.advance()
                return self.new_token(TokenType.COMMA, ',')

            # Caracteres de 1 o más digitos de longitud
            if self.current_char == '=':
                self.advance()  # Avanza el '='
                if self.current_char == '=':
                    self.advance()  # Avanza el segundo '='
                    return self.new_token(TokenType.EQUAL, '==')

                return self.new_token(TokenType.ASSIGN, '=')

            if self.current_char == '!':
                self.advance()  # Avanza el '!'
                if self.current_char == '=':
                    self.advance()  # Avanza el '='
                    return self.new_token(TokenType.NOT_EQUAL, '!=')

                return self.new_token(TokenType.NOT, '!')

            if self.current_char == '<':
                self.advance()  # Avanza el '<'
                if self.current_char == '=':
                    self.advance()  # Avanza el '='
                    return self.new_token(TokenType.LESS_EQ, '<=')

                return self.new_token(TokenType.LESS, '<')

            if self.current_char == '>':
                self.advance()  # Avanza el '>'
                if self.current_char == '=':
                    self.advance()  # Avanza el '='
                    return self.new_token(TokenType.GREATER_EQ, '>=')

                return self.new_token(TokenType.GREATER, '>')

            # Caracter de impresión
            if self.current_char == '?':
                self.advance()
                return self.new_token(TokenType.PRINT, '?')

            self.error()  # Caracter desconocido

        return self.new_token(TokenType.EOF, None)


"""
 Métodos Helper del Lexer
"""


def is_space(ch):
    return ch in (' ', '\t', '\r')


def is_letter(ch):
    return ch.isalpha() or ch.isdigit() or ch == '_'


if __name__ == '__main__':
    source = """
    && Ejemplo FoxLite
    x = 10 && Declaración de variable 'x'
    y = 20 && Declaración de variable 'y'
    
    && Prueba de IF
    
    if x >= y
      messagebox("x es mayor")
    else
      ? "x es menor"
    endif
    
    
    do while x < 99
       x = x + 1
       ? "Contando por", x
       
       if x == 55
          return
       endif
    enddo    
    """
    lexer = Lexer(source_code=source)
    tok = lexer.next_token()

    while tok.value is not None:
        print(tok)
        tok = lexer.next_token()
