"""
Funciones integradas en FoxLite
"""
import src.book.fox_lite_object as obj


def len_func(args):
    """
    Devuelve la longitud de una cadena.
    :param args:
    :return: obj.Integer()
    """
    if len(args) != 1:
        return obj.Error(message=f"número de argumentos inválidos. Se esperaba 1 y llegaron {len(args)})")

    if args[0].type() == obj.Type.STRING:
        return obj.Integer(value=len(args[0].value))
    else:
        return obj.Error(message=f"el tipo {args[0].type()} no es soportado por la función len().")


def type_func(args):
    """
    Devuelve el tipo de dato del argumento pasado.
    :param args:
    :return: obj.String()
    """
    if len(args) != 1:
        return obj.Error(message=f"número de argumentos inválidos. Se esperaba 1 y llegaron {len(args)})")

    obj_type = args[0].type()
    type_value = 'U'
    if obj_type == obj.Type.STRING:
        type_value = 'C'
    elif obj_type == obj.Type.BOOLEAN:
        type_value = 'L'
    elif obj_type == obj.Type.INTEGER:
        type_value = 'N'
    elif obj_type == obj.Type.NULL:
        type_value = 'X'

    return obj.String(value=type_value)


def messagebox_func(args):
    """
    Muestra un cuadro de diálogo en pantalla.
    :param args:
    :return:
    """
    if len(args) != 2:
        return obj.Error(message=f"número de argumentos inválidos. Se esperaban 2 y llegaron {len(args)}")

    import ctypes
    ctypes.windll.user32.MessageBoxW(None, args[0].resolve(), args[1].resolve(), 0)


def alltrim_len(args):
    """
    Elimina los espacios en blanco a la izquierda y a la derecha de una cadena.
    :param args:
    :return:
    """
    if len(args) != 1:
        return obj.Error(message=f"número de argumentos inválidos. Se esperaba 1 y llegaron {len(args)}")

    return obj.String(value=args[0].resolve().strip())


builtins = {
    "len": obj.Builtin(function=len_func),
    "type": obj.Builtin(function=type_func),
    "messagebox": obj.Builtin(function=messagebox_func),
    "alltrim": obj.Builtin(function=alltrim_len),
}