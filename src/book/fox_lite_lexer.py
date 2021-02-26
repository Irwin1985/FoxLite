from src.book.fox_lite_token import Token, TokenType, lookup_ident


def is_space(ch):
    """
    Determina si el carácter es un espacio en blanco, tabulación o retorno de carro.
    El Line Feed (ENTER) no lo ignoramos porque es relevante en FoxLite.
    :param ch:
    :return: Boolean
    """
    return ch in (' ', '\t', '\r')


def is_letter(ch):
    """
    Determina si el carácter es alfanumérico o guion bajo.
    :param ch:
    :return:
    """
    return ch.isalpha() or ch.isdigit() or ch == '_'


class Lexer:
    """
      Lexer: también conocido como tokenizer, se encarga de leer y
      descomponer el código fuente en pequeñas unidades llamadas tokens.
    """
    def __init__(self, source_code):
        self.source = source_code
        self.current_pos = 0  # Mantiene la posición del carácter actual.
        self.current_char = self.source[self.current_pos]
        self.last_token_value = None  # Último valor del token generado.

    def new_token(self, token_type, token_value):
        """
        Crea un token y actualiza el valor del último token generado.
        :param token_type: el tipo del token a generar.
        :param token_value: el valor o lexema del token a generar.
        :return: Token()
        """
        token = Token(token_type=token_type, token_value=token_value)
        self.last_token_value = token_value
        return token

    def advance(self):
        """
        Mueve el puntero a la siguiente posición del código fuente.
        """
        self.current_pos += 1
        if self.current_pos >= len(self.source):
            self.current_char = None
        else:
            self.current_char = self.source[self.current_pos]

    def peek(self):
        """
        Mira 1 caracter hacia adelante en el código fuente.
        :return: char
        """
        peek_position = self.current_pos + 1
        if peek_position >= len(self.source):
            return None
        else:
            return self.source[peek_position]

    # Ignoramos los espacios en blanco
    def ignore_blanks(self):
        """
        Ignora los espacios en blanco, tabulaciones y retornos de carro.
        :return: None
        """
        # Avanzamos los caracteres mientras hayan espacios o EOF
        while self.current_char is not None and is_space(self.current_char):
            self.advance()

    def ignore_comments(self):
        """
        Ignora los comentarios de FoxLite '&&'
        Avanza los caracteres hasta el final de la línea.
        :return: None
        """
        while self.current_char is not None and self.current_char != '\n':
            self.advance()

    def number(self):
        """
        Extrae un número entero y genera un token de tipo INT.
        :return: Token()
        """
        lexeme = ''
        while self.current_char is not None and self.current_char.isdigit():
            lexeme += self.current_char
            self.advance()

        return self.new_token(TokenType.INT, lexeme)

    def identifier(self):
        """
        Obtiene un identificador.
        :return: Token()
        """
        lexeme = ''
        while self.current_char is not None and is_letter(self.current_char):
            lexeme += self.current_char
            self.advance()

        token_type = lookup_ident(lexeme)
        return self.new_token(token_type, lexeme)

    def dotted_indentifier(self):
        """
        Extrae un identificador encerrado con '.' ej: .t., .f., .null.
        :return: Token()
        """
        lexeme = '.'
        self.advance()  # Avanza el primer punto '.'

        while self.current_char is not None and self.current_char != '.':
            lexeme += self.current_char
            self.advance()

        lexeme += '.'
        self.advance()  # Avanza el segundo punto '.'

        token_type = lookup_ident(lexeme)
        return self.new_token(token_type, lexeme)

    def string(self, string_delim):
        """
        Obtiene una secuencia de caracteres (string)
        :param string_delim: Foxlite permite comilla simple y doble.
        :return: Token()
        """
        self.advance()  # Avanza el primer delimitador del string.
        lexeme = ''
        while self.current_char is not None and self.current_char != string_delim:
            lexeme += self.current_char
            self.advance()

        self.advance()  # Avanza el segundo delimitador del string.
        return self.new_token(TokenType.STRING, lexeme)

    def next_token(self):
        """
        Extrae el siguiente token desde el código fuente.
        :return: Token()
        """
        while self.current_char is not None:
            # Salto de Línea
            if self.current_char == '\n':
                self.advance()  # Avanza el '\n'
                # Si el token anterior no es LBREAK entonces generamos el token LBREAK.
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
            if self.current_char == '"' or self.current_char == "'":
                return self.string(self.current_char)

            # Caracteres de 1 dígito de longitud
            if self.current_char == '+':
                self.advance()
                return self.new_token(TokenType.PLUS, '+')

            if self.current_char == '-':
                self.advance()
                return self.new_token(TokenType.MINUS, '-')

            if self.current_char == '*':
                self.advance()
                return self.new_token(TokenType.MUL, '*')

            if self.current_char == '/':
                self.advance()
                return self.new_token(TokenType.DIV, '/')

            if self.current_char == '(':
                self.advance()
                return self.new_token(TokenType.LPAREN, '(')

            if self.current_char == ')':
                self.advance()
                return self.new_token(TokenType.RPAREN, ')')

            if self.current_char == ',':
                self.advance()
                return self.new_token(TokenType.COMMA, ',')

            # Caracteres de 1 o más dígitos de longitud
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

            # Carácter de impresión
            if self.current_char == '?':
                self.advance()
                return self.new_token(TokenType.PRINT, '?')

            raise Exception(f"Error: carácter desconocido '{self.current_char}'")

        return self.new_token(TokenType.EOF, None)
