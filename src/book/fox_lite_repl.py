from src.book.fox_lite_lexer import Lexer
from src.book.fox_lite_parser import Parser
from src.book.fox_lite_evaluator import Evaluator
from src.book.fox_lite_environment import Environment

FOX_LITE = """
   /\   /\   ¡FoxLite!
  //\\_//\\     ____
  \_     _/    /   /
   / * * \    /^^^]
   \_\O/_/    [   ]
    /   \_    [   /
    \     \_  /  /
     [ [ /  \/ _/
    _[ [ \  /_/
"""

ERROR_FACE = """
    |\__/|
   /     \
  /_.@ @,_\
     \@/
"""


def repl():
    print(FOX_LITE)
    print('¡Hola Mundo! Este es lenguage de programación FoxLite')
    print('¡Adelante! Ingresa algunos comandos...\n')
    env = Environment()
    while True:
        try:
            source_code = input(">> ")
        except EOFError:
            break
        if not source_code:
            continue

        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()

        if len(parser.errors) != 0:
            print_parser_errors(parser.errors)

        evaluator = Evaluator()
        evaluated = evaluator.eval(node=program, env=env)
        if evaluated is not None:
            print(evaluated.resolve())


def print_parser_errors(errors):
    print(ERROR_FACE)
    print("¡Ay caramba! parece que tenemos un problema.")
    print(" parser errors:")
    for error in errors:
        print(error)


if __name__ == '__main__':
    repl()
