import unittest

from src.book.fox_lite_environment import Environment
from src.book.fox_lite_evaluator import Evaluator
from src.book.fox_lite_lexer import Lexer
from src.book.fox_lite_parser import Parser
import src.book.fox_lite_object as obj


class TestEvaluator(unittest.TestCase):
    def test_eval_assignment(self):
        source = """
        local foo
        foo = 10
        public foo
        """
        evaluated = self.assert_test_eval(source)

    def test_eval_boolean_expression(self):
        tests = [
            [".t.", True],
            [".f.", False],
            [".t. == .t.", True],
            [".f. == .f.", True],
            [".t. == .f.", False],
            [".t. == .f.", False],
            [".f. != .t.", True],
            [".t. < .f.", False],
            [".f. < .f.", False],
            [".t. > .f.", True],
            [".f. > .t.", False],
            [".t. <= .f.", False],
            [".f. <= .t.", True],
            [".t. >= .f.", True],
            [".f. >= .t.", False],
            [".t. and .f.", False],
            [".f. and .t.", False],
            [".t. or .f.", True],
            [".f. or .t.", True],
            ["(1 < 2) == .t.", True],
            ["(1 < 2) == .f.", False],
            ["(1 > 2) == .t.", False],
            ["1 < 2", True],
            ["1 > 2", False],
            ["1 < 1", False],
            ["1 > 1", False],
            ["1 == 1", True],
            ["1 != 1", False],
            ["1 == 2", False],
            ["1 != 2", True],
        ]
        for tt in tests:
            source = tt[0]
            expected = tt[1]
            evaluated = self.assert_test_eval(source)
            self.assert_test_boolean_object(evaluated, expected)

    def test_eval_integer_expression(self):
        tests = [
            ["5", 5],
            ["10", 10],
            ["-5", -5],
            ["-10", -10],
            ["5 + 5 + 5 + 5 - 10", 10],
            ["2 * 2 * 2 * 2 * 2", 32],
            ["-50 + 100 + -50", 0],
            ["5 + 2 * 10", 25],
            ["20 + 2 * -10", 0],
            ["50 / 2 * 2 + 10", 60],
            ["2 * (5 + 10)", 30],
            ["3 * 3 * 3 + 10", 37],
            ["(5 + 10 * 2 + 15 / 3) * 2 + -10", 50],
        ]

        for tt in tests:
            source = tt[0]
            expected = tt[1]

            evaluated = self.assert_test_eval(source)
            self.assert_test_integer_object(evaluated, expected)

    def assert_test_integer_object(self, obj_type, expected):
        result = obj_type  # Object
        if type(result) is not obj.Integer:
            print(f'object is not Integer. got={type(obj)}')
            return False

        if result.value != expected:
            print(f'object has wrong value. got={result.value}, want={expected}')
            return False
        return True

    def assert_test_boolean_object(self, obj_type, expected):
        result = obj_type
        if type(result) is not obj.Boolean:
            print(f'object is not Boolean. got={type(result)}')
        if result.value != expected:
            print(f'object has wrong value. got={result.value}, want={expected}')
            return False
        return True

    def assert_test_eval(self, source):
        lexer = Lexer(source)
        parser = Parser(lexer)
        program = parser.program()
        eva = Evaluator()
        env = Environment()

        return eva.eval(ast_node=program, env=env)


if __name__ == '__main__':
    unittest.main()