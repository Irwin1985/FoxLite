from src.lexer import Lexer
from src.tokens import Token, TokenType

import unittest


class TestLexer(unittest.TestCase):
    def test_tokens(self):
        source_code = """
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
        expected_tokens = [
            Token(TokenType.IDENT, "x"),
            Token(TokenType.ASSIGN, "="),
            Token(TokenType.INT, "10"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.IDENT, "y"),
            Token(TokenType.ASSIGN, "="),
            Token(TokenType.INT, "20"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.IF, "if"),
            Token(TokenType.IDENT, "x"),
            Token(TokenType.GREATER_EQ, ">="),
            Token(TokenType.IDENT, "y"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.MESSAGEBOX, "messagebox"),
            Token(TokenType.LPAREN, "("),
            Token(TokenType.STRING, "x es mayor"),
            Token(TokenType.RPAREN, ")"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.ELSE, "else"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.PRINT, "?"),
            Token(TokenType.STRING, "x es menor"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.ENDIF, "endif"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.DO, "do"),
            Token(TokenType.WHILE, "while"),
            Token(TokenType.IDENT, "x"),
            Token(TokenType.LESS, "<"),
            Token(TokenType.INT, "99"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.IDENT, "x"),
            Token(TokenType.ASSIGN, "="),
            Token(TokenType.IDENT, "x"),
            Token(TokenType.PLUS, "+"),
            Token(TokenType.INT, "1"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.PRINT, "?"),
            Token(TokenType.STRING, "Contando por"),
            Token(TokenType.COMMA, ","),
            Token(TokenType.IDENT, "x"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.IF, "if"),
            Token(TokenType.IDENT, "x"),
            Token(TokenType.EQUAL, "=="),
            Token(TokenType.INT, "55"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.RETURN, "return"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.ENDIF, "endif"),
            Token(TokenType.LBREAK, "LBREAK"),
            Token(TokenType.ENDDO, "enddo"),
            Token(TokenType.LBREAK, "LBREAK"),
        ]
        lexer = Lexer(source_code=source_code)

        for expected_token in expected_tokens:
            actual = lexer.next_token()
            self.assertEqual(expected_token.type, actual.type)
            self.assertEqual(expected_token.value, actual.value)


if __name__ == '__main__':
    unittest.main()
