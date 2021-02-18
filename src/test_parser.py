import unittest
from src.fox_lite_parser import Parser
from src.fox_lite_lexer import Lexer
from src.fox_lite_ast import (
    Boolean,
    Integer,
    Identifier,
    Null,
    VariableDecl,
)


class TestParser(unittest.TestCase):
    def test_assignment_statement(self):
        lexer = Lexer("foo = bar")
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements))

        assignment = program.statements[0]
        ident_var = assignment.name
        ident_value = assignment.value

        self.assertEqual("foo", ident_var.value, f'ident.value is not foo. got={ident_var.value}')
        self.assertEqual("bar", ident_value.value, f'ident_value.value is not bar. got={ident_value.value}')

    def test_do_while_statement(self):
        source_code = """
        do while .t.
            return 1
        enddo
        """
        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        do_while = program.statements[0]

        condition = do_while.condition
        self.assertEqual(True, condition.value, f'condition.value is not True, got={condition.value}')

        body = do_while.block.statements[0]
        integer = body.value

        self.assertEqual(1, integer.value, f'integer.value is not 1, got={integer.value}')

    def test_function_call_with_args(self):
        lexer = Lexer("foo(1, 2)")
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        func = program.statements[0]
        self.assertEqual("foo", func.name.value, f'func.name no es foo, got={func.name.value}')

        # Validar los argumentos
        arg1 = func.arguments[0]
        arg2 = func.arguments[1]
        self.assertEqual(1, arg1.value, f'arg1.value no es 1, got={arg1.value}')
        self.assertEqual(2, arg2.value, f'arg2.value no es 2, got={arg2.value}')

    def test_function_call_no_params(self):
        lexer = Lexer("foo()")
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        func = program.statements[0]
        self.assertEqual("foo", func.name.value, f'func.name no es foo, got={func.name.value}')

    def test_function_with_params(self):
        source_code = """
        function retorna_1(foo, bar)
            return 1
        endfunc
        """
        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        func = program.statements[0]
        self.assertEqual("retorna_1", func.name.value)

        identi1 = func.params[0]
        self.assertEqual("foo", identi1.value, f'identi1.value no es foo, got={identi1.value}')

        identi2 = func.params[1]
        self.assertEqual("bar", identi2.value, f'identi1.value no es bar, got={identi2.value}')

        return_ast = func.body.statements[0]
        integer = return_ast.value

        self.assertEqual(1, integer.value, f'integer.value no es 1, got={integer.value}')

    def test_function_with_parenthesis_and_no_params(self):
        source_code = """
        function retorna_1()
            return 1
        endfunc
        """
        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        func = program.statements[0]
        self.assertEqual("retorna_1", func.name.value)

        return_ast = func.body.statements[0]
        integer = return_ast.value

        self.assertEqual(1, integer.value, f'integer.value no es 1, got={integer.value}')

    def test_function_without_parenthesis(self):
        source_code = """
        function retorna_1
            return 1
        endfunc
        """
        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        func = program.statements[0]
        self.assertEqual("retorna_1", func.name.value)

        return_ast = func.body.statements[0]
        integer = return_ast.value

        self.assertEqual(1, integer.value, f'integer.value no es 1, got={integer.value}')

    def test_if_else_statement(self):
        source_code = """
        if .t.
            return 1
        else
            return 2
        endif
        """
        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        if_stmt = program.statements[0]

        self.assertEqual(if_stmt.condition.value, True)  # Condición del if (.t.)
        # Consecuencia
        return_if_ast = if_stmt.consequence.statements[0]
        integer_ast = return_if_ast.value
        self.assertEqual(1, integer_ast.value, f'integer_ast.value is not 1, got={integer_ast.value}.')
        # Alternativa
        return_else_ast = if_stmt.alternative.statements[0]
        integer_ast = return_else_ast.value
        self.assertEqual(2, integer_ast.value, f'integer_ast.value is not 2, got={integer_ast.value}.')

    def test_if_statement(self):
        source_code = """
        if .t.
            return 1
        endif
        """
        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements),
                         f'program.statements does not contain 1 statements. got={len(program.statements)}')

        if_stmt = program.statements[0]

        self.assertEqual(if_stmt.condition.value, True)  # Condición del if (.t.)
        return_ast = if_stmt.consequence.statements[0]
        integer_ast = return_ast.value
        self.assertEqual(1, integer_ast.value, f'integer_ast.value is not 1, got={integer_ast.value}.')

    def test_binary_expressions(self):
        tests = [
            ["5 + 5", 5, "+", 5],
            ["5 - 5", 5, "-", 5],
            ["5 * 5", 5, "*", 5],
            ["5 / 5", 5, "/", 5],
            ["5 < 5", 5, "<", 5],
            ["5 > 5", 5, ">", 5],
            ["5 == 5", 5, "==", 5],
            ["5 != 5", 5, "!=", 5],
            [".t. == .t.", True, "==", True],
            [".t. != .f.", True, "!=", False],
            [".f. == .f.", False, "==", False],
        ]

        for tt in tests:
            source = tt[0]
            left = tt[1]
            operator = tt[2]
            right = tt[3]

            lexer = Lexer(source)
            parser = Parser(lexer)
            program = parser.program()
            self.assert_check_parser_errors(parser)

            self.assertEqual(1, len(program.statements),
                             f'program.statements does not contain 1 statements. got={len(program.statements)}')

            exp = program.statements[0]

            self.assertEqual(exp.left.value, left, f'exp.left is not {left}. got={exp.left.value}')
            self.assertEqual(exp.operator, operator, f"exp.operator is not '{operator}'. got={exp.operator}")
            self.assertEqual(exp.right.value, right, f'exp.right is not {right}. got={exp.right.value}')

    def test_unary_expressions(self):
        tests = [
            ["!5", "!", 5],
            ["-15", "-", 15],
            ["!.t.", "!", True],
            ["!.f.", "!", False],
        ]

        for tt in tests:
            source = tt[0]
            operator = tt[1]
            value = tt[2]

            lexer = Lexer(source)
            parser = Parser(lexer)
            program = parser.program()
            self.assert_check_parser_errors(parser)

            self.assertEqual(1, len(program.statements),
                             f'program.statements does not contain 1 statements. got={len(program.statements)}')

            exp = program.statements[0]

            self.assertEqual(exp.operator, operator, f"exp.operator is not '{operator}'. got={exp.operator}")
            self.assertEqual(exp.right.value, value, f'exp.right is not {value}. got={exp.right.value}')

    def test_variable_declaration(self):
        tests = [
            ["public a", VariableDecl(token=Identifier(value="a"), scope='public')],
            ["local b", VariableDecl(token=Identifier(value="b"), scope='local')],
            ["private c", VariableDecl(token=Identifier(value="c"), scope='private')],
        ]
        for tt in tests:
            source = tt[0]
            expected = tt[1]

            lexer = Lexer(source)
            parser = Parser(lexer)
            program = parser.program()
            self.assert_check_parser_errors(parser)
            self.assertEqual(1, len(program.statements))

            variable = program.statements[0]

            self.assertEqual(expected.token.value, variable.token.value,
                             f'expected.token.value is not {expected.token.value}. got={variable.token.value}')

            self.assertEqual(expected.scope, variable.scope,
                             f'expected.scope is not {expected.scope}. got={variable.scope}')

    def test_return_statement(self):
        tests = [
            ["return \n", Boolean(value=True)],
            ["return .f.", Boolean(value=False)],
            ["return 5", Integer(value=5)],
            ["return foobar", Identifier(value="foobar")],
            ["return .null.", Null()]
        ]
        for tt in tests:
            source = tt[0]
            expected = tt[1]

            lexer = Lexer(source)
            parser = Parser(lexer)
            program = parser.program()
            self.assert_check_parser_errors(parser)
            self.assertEqual(1, len(program.statements))

            return_stmt = program.statements[0]
            return_value = return_stmt.value
            self.assertEqual(expected.value, return_value.value,
                             f'return_stmt.value is not {expected.value}. got={return_value.value}')

    def test_null(self):
        lexer = Lexer(".null.")
        parser = Parser(lexer)
        program = parser.program()
        self.assert_check_parser_errors(parser)
        self.assertEqual(1, len(program.statements))

        null = program.statements[0]
        self.assertEqual(None, null.value, f'boolean.value is not None. got={null.value}')

    def test_boolean(self):
        tests = [
            [".t.", True],
            [".f.", False],
        ]
        for tt in tests:
            source = tt[0]
            expected = tt[1]

            lexer = Lexer(source)
            parser = Parser(lexer)
            program = parser.program()
            self.assert_check_parser_errors(parser)
            self.assertEqual(1, len(program.statements))

            boolean = program.statements[0]
            self.assertEqual(expected, boolean.value, f'boolean.value is not {expected}. got={boolean.value}')

    def test_integer(self):
        lexer = Lexer("5")
        parser = Parser(lexer)

        program = parser.program()
        self.assert_check_parser_errors(parser)

        self.assertEqual(1, len(program.statements),
                         f'program has not enough statements. got={len(program.statements)}')
        integer = program.statements[0]
        self.assertEqual(5, integer.value, f'integer.value is not 5. got={integer.value}')

    def test_identifier(self):
        lexer = Lexer("foobar")
        parser = Parser(lexer)

        program = parser.parse()
        self.assert_check_parser_errors(parser)

        self.assertEqual(1, len(program.statements),
                         f'program has not enough statements. got={len(program.statements)}')

        ident = program.statements[0]

        self.assertEqual("foobar", ident.value, f'ident.value not "foobar" got={ident.value}')

    def assert_check_parser_errors(self, parser):
        if len(parser.errors) == 0:
            return

        print(f'parser has {len(parser.errors)} errors.')
        for error in parser.errors:
            print(f'parser error: {error}')

        self.fail()


if __name__ == '__main__':
    unittest.main()