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


builtins = {
    "len": obj.Builtin(function=len_func),
    "type": obj.Builtin(function=type_func),
}