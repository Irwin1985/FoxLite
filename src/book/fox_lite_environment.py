"""
    Esta clase se encarga de guardar y recuperar todos los símbolos del programa.
"""


class Environment:
    def __init__(self):
        self.symbol_table = {}
        self.parent = None

    def get(self, name):
        result = self.symbol_table.get(name)

        if result is None and self.parent is not None:
            result = self.parent.get(name)

        return result

    def set(self, name, val):
        self.symbol_table[name] = val


def new_environment(parent):
    env = Environment()
    env.parent = parent
    return env
