"""
    Esta clase se encarga de guardar y recuperar todos los símbolos del programa.
"""

from src.book.fox_lite_object import Error


class Environment:
    def __init__(self, parent=None):
        self.symbol_table = {}
        self.parent = parent
    """
    Resuelve el nombre de un símbolo dado.
    1. Primero busca el símbolo en el ámbito actual
    2. Si no existe lo busca como 'public' o 'private' en el ámbito padre.
    3. Si no existe entonces crea el símbolo 'private' por defecto.
    4. Si existe entonces actualiza el valor del símbolo.
    """
    def set(self, name, val, scope='default', current_scope=False):
        # Buscamos el símbolo en el entorno actual
        result = self.symbol_table.get(name)

        if result is not None:
            if result[0] in ('local', 'private') and scope == 'public':  # Verificar la redefinición no válida.
                return Error(message=f"redefinición no válida de la variable '{name}'.")

            result[1] = val  # Actualizamos el valor del símbolo.
            self.symbol_table[name] = result
        else:
            if self.parent is not None and not current_scope:
                result = self.parent.get(name)  # Buscamos el símbolo en el entorno padre.
                if result is not None:
                    if result[0] in ('public', 'private'):  # Solo permitimos 'public' o 'private'
                        result[1] = val  # Actualizamos el valor del símbolo en el entorno padre.
                        self.parent[name] = result
                    else:
                        self.symbol_table[name] = [scope, val]  # Registramos el símbolo.
                else:
                    self.symbol_table[name] = [scope, val]  # Registramos el símbolo.
            else:
                self.symbol_table[name] = [scope, val]  # Registramos el símbolo.

        return val

    def get(self, name):
        result = self.symbol_table.get(name)

        if result is None and self.parent is not None:
            result = self.parent.get(name)

        return None if result is None else result[1]
