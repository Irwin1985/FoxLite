FUNCTION DAME_TU_NOMBRE(NOMBRE)
    FUNCTION DAME_TU_APELLIDO(APELLIDO)
        RETURN "TU NOMBRE COMPLETO ES: " + NOMBRE + ", " + APELLIDO
    ENDFUNC
    RETURN DAME_TU_APELLIDO
ENDFUNC

FUNCTION CUADRADO(NUMERO)
    RETURN NUMERO * NUMERO
ENDFUNC

FUNCTION CUBO(NUMERO, FUN_CUADRADO)
    RETURN FUN_CUADRADO(NUMERO) * NUMERO
ENDFUNC