"""
Sistema de objetos de FoxLite
"""

from enum import Enum


class Type(Enum):
    INTEGER = "INTEGER"
    STRING = "STRING"
    BOOLEAN = "BOOLEAN"
    NULL = "NULL"
    RETURN = "RETURN"
    ERROR = "ERROR"


class ObjectType:
    pass


class Object:
    def type(self):
        pass

    def resolve(self):
        pass


class Integer(Object):
    def __init__(self, value):
        self.value = value

    def resolve(self):
        return str(self.value)

    def type(self):
        return Type.INTEGER


class Boolean(Object):
    def __init__(self, value):
        self.value = value

    def resolve(self):
        return str(self.value)

    def type(self):
        return Type.BOOLEAN


class Null(Object):
    def __init__(self):
        self.value = None

    def resolve(self):
        return "Null"

    def type(self):
        return Type.NULL


class String(Object):
    def __init__(self, value):
        self.value = value

    def resolve(self):
        return self.value

    def type(self):
        return Type.STRING


class Error(Object):
    def __init__(self, message):
        self.message = message

    def resolve(self):
        return self.message

    def type(self):
        return Type.ERROR
