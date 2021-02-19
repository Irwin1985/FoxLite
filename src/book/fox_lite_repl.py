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
    print_header()
    env = Environment()
    source_code = ''
    while True:
        try:
            user_input = input(">> ")
        except EOFError:
            break
        if not user_input:
            continue

        if user_input[len(user_input)-1] == ';':
            # Unimos la línea de código
            source_code += '\n' + user_input[0:len(user_input)-1]
            continue
        else:
            source_code += '\n' + user_input

        lexer = Lexer(source_code)
        parser = Parser(lexer)
        program = parser.program()

        if len(parser.errors) != 0:
            print_parser_errors(parser.errors)

        evaluator = Evaluator()
        evaluated = evaluator.eval(node=program, env=env)
        if evaluated is not None:
            print(evaluated.resolve())

        source_code = ''


def print_parser_errors(errors):
    print(ERROR_FACE)
    print("¡Ay caramba! parece que tenemos un problema.")
    print(" parser errors:")
    for error in errors:
        print(error)


def print_header():
    print(FOX_LITE)
    print('¡Hola Mundo! Este es lenguage de programación FoxLite')
    print('¡Adelante! Ingresa algunos comandos...\n')


if __name__ == '__main__':
    repl()
