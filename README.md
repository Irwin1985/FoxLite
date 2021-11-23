# FoxLite
Mis pensamientos acerca del diseño que le daría al lenguage de programación **Foxlite**.

# Introducción
Hola, mi nombre es Irwin y soy el creador de **FoxLite**, un lenguaje de programación que pretende ser una evolución de Fox. He programado en **Visual FoxPro** durante 12 años y le tengo un gran aprecio pues me ha dado de comer por durante todo este tiempo y lo menos que puedo hacer por él es evolucionarlo. Tal vez no te guste **FoxLite** pero creeme que el diseño que vas a ver aquí es lo más parecido a un **Visual FoxPro** evolucionado, así que comencemos a ver lo que deja, lo que hereda y lo que adquiere. 

Este diseño es meramente subjetivo así que la evolución de **Fox** está en mi mano pero no te preocupes porque conozco la gramática de **Fox** y la gramática de los lenguajes *"modernos"* por lo tanto me aseguraré que sus nuevos poderes lo ayuden a sobrevir en este nuevo ecosistema.

## Un repaso a Visual FoxPro
Comencemos viendo un trozo de sintaxis de **Visual FoxPro**

Lo siguiente son versiones del "Hola Mundo" en Fox.

```xBase
* Versión 1:
 ? "Hola Mundo"

* Versión 2:
 @ 1,1 SAY "Hola Mundo"

* Versión 3:
 WAIT WINDOW "Hola Mundo"

* Versión 4:
 MESSAGEBOX("Hola Mundo")

* Version 5: macro-sustitución
 cVariable = "Hola Mundo"
 ? &cVariable 

* Version 6:
 loForm = CREATEOBJECT("HolaForm")
 loForm.Show(1)
 DEFINE CLASS HolaForm AS Form
    AutoCenter= .T.
    Caption= "Hola Mundo"
    ADD OBJECT lblHola as Label ;
      WITH Caption = "Hola Mundo"
 ENDDEFINE
```

### Ahora un ejemplo extendido mostrando más de su sintaxis y gramática
```xBase
* COMENTARIOS
* Comentario de una línea con '*'
&& Comentario de una línea con '&&'

* TIPOS DE DATOS
? .T. && Verdadero (Boolean)
? .F. && Falso (Boolean)
? "FoxLite" && String
? 1985 && Number (no hay distinción)
? 3.14159265 && Number (no hay distinción)

* VARIABLES
PRIVATE lnID
lnID = 10

LOCAL lcNombre
lcNombre = "Juan"

PUBLIC gApellido
gApellido = "Gonzalez"

* CASTING
lcEdad = "36"
?VAL(lcEdad) && Convierte de String a Number
?STR(lnEdad) && Convierte de Number a String

* STRINGS
a = "Lorem ipsum dolor sit amet"
?a

* STRINGS MÚLTIPLES
TEXT TO a NOSHOW
Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt
ut labore et dolore magna aliqua.
ENDTEXT
?a

* OPERADORES ARITMÉTICOS
?10 + 5
?10 - 5
?10 * 5
?10 / 5

* OPERADORES RELACIONALES
?10 == 10
?10 != 10
?10 < 10
?10 > 10
?10 <= 10
?10 >= 10

* OPERADORES LÓGICOS
?.T. AND .F.
? .T. OR .F.
? NOT .T. && Versión 1
? !.T.    && Versión 2

* ARRAYS
DIMENSION aLenguajes[3]
aLenguajes[1] = "FoxLite"
aLenguajes[2] = "FoxPro"
aLenguajes[3] = "Python"

* DICCIONARIOS
data = CREATEOBJECT("Collection")
data.Add("Juan", "nombre")
data.Add(36, "edad")
data.Add(.F., "casado")

?data.Item(data.GetKey("nombre")) && Juan
?data.Item(data.GetKey("edad")) && 36
?data.Item(data.GetKey("casado")) && .F.

* IF / ELSE
IF .T.
    && Algo
ELSE
    && Otro
ENDIF

* DO CASE / OTHERWISE
DO CASE
CASE .T.
    && Algo
CASE .F.
    && Algo
OTHERWISE
    && Algo
ENDCASE

* LOOPS

DO WHILE .T.
    ? "Gracias infinitas Fox"
ENDDO

FOR I = 1 TO 10
    ? "Gracias " + STR(I) + " Fox"
ENDFOR

* FUNCIONES / PROCEDURES

FUNCTION Sumar(x, y)
    RETURN x + y
ENDFUNC

PROCEDURE Restar(x, y)
    RETURN x - y
ENDPROC

* CLASES
DEFINE CLASS Persona AS CUSTOM
    nombre = "Juan"
    edad = 36
    casado = .F.

    FUNCTION Presentarse
        RETURN "Hola me llamo " + THIS.nombre + " y tengo " + STR(THIS.edad) + " años"
    ENDFUNC
ENDDEFINE

* INSTANCIACIÓN DE CLASES
oJuan = CREATEOBJECT("Persona")
?oJuan.Presentarse()
```

Lo anterior es apenas una porción de la sintaxis de **Visual FoxPro**, puede hacer más cosas pero ustedes son **Foxeros** así que ya tienen una idea de lo que es capaz de hacer.

**FoxLite** heredará gran parte de la semántica de **FoxPro** pero también es verdad que perderá parte de la sintaxis para favorecer su modernidad.

Veamos algunas de los aspectos que **FoxLite** no obtendrá de su anscestro:

## Limando asperezas en el proceso evolutivo

- **Literales Booleanos:** honestamente **.T.** y **.F.** no me molestan pero he decidido tender una rama hacía el árbol **ALGOL** por lo que ahora serán **true** y **false**. La vieja versión es incluso mejor ya que escribes menos *(3 letras en lugar de 4)* pero esto lo hago para que **FoxLite** sea bien visto por la comunidad cuya raíz desde luego parte de **ALGOL**.
```Javascript
   verdad = true
   mentira = false
```
- **Literales Arrays:** debo admitir que la sintaxis **DIMENSION arrayName()** nunca me gustó así que esta exclusión es personal *(al igual que todas las demás)* 😋. **FoxLite** adoptará la forma simplística que muchos lenguajes tienen pero que fue popularizada por **Javascript**.

```Javascript
   numeros = [1, 2, 3]
   frutas = ["Manzana", "Mango", "Mora"]
```
- **Funciones Builtins:** estas son las funciones que un lenguaje tiene integradas en su núcleo y que por lo tanto podemos utilizar en cualquier script. **Visual FoxPro** tiene montones de ellas clasificadas por tipos de datos. Aunque particularmente me gusta trabajar con ellas pienso que es mejor adherirlas a su tipo correspondiente y así *limitar* la cantidad de funciones integradas que **FoxLite** debe cargar en sus hombros. Esto tiene un coste que aún sigo sopesando pero creo que al final me decantaré por este diseño. 

Veamos algunos ejemplos:
```xBase
   && Version VFP
   nombre = "FoxLite   "
   ?LEN(nombre) && 10
   ?LEN(ALLTRIM(nombre)) && 7
```
El ejemplo anterior está escrito en **Visual FoxPro** y muestra el uso de 2 *funciones integradas*: **ALLTRIM()** y **LEN()** donde la primera trabaja con *Characters* o *String* y la segunda con *números*.

Ahora veamos la versión en **FoxLite**:
```Javascript
   // version foxlite
   nombre = "FoxLite   "
   ?nombre.len() // 10
   ?nombre.trim().len() // 7
```
Como habrás notado **ALLTRIM()** ha perdido parte de su pelaje y ahora es solo **TRIM()** que es una versión resumida y significa lo mismo, esta nueva versión quizá no te agrade mucho pero es una forma de mantener las funciones integradas adheridas a sus tipos. Es verdad que no previenen los errores porque si invocas la función **trim()** sobre un tipo numérico obtendrás un error de incompatibilidad de tipos pero esto tampoco tiene que ser una desventaja porque para eso existen los [linters](https://es.wikipedia.org/wiki/Lint) que ayudan a detectar errores en tiempo de diseño. De esto no tienes que preocuparte porque un linter se puede desarrollar e incrustar dentro de un editor bien sea propio de **FoxLite** o un tercero como **VsCode, Atom, etc.**

**ALLTRIM()** y **LEN()** son solo el abrebocas de todo el cambio que sufrirán las funciones integradas. Algunas las agradecerás y otras las lamentarás pero es por el bien de FoxLite y su nuevo ecosistema.

**SUBSTR()** está sentenciada a muerte porque lo mismo se puede lograr de la siguiente manera:
```Javascript
   // Versión Fox
   nombre = "FoxLite"
   ?SUBSTR(nombre, 1, 3) // Fox
   // Versión FoxLite
   ?nombre[0:3]
```

- **Sensibilidad a las Mayúsculas:** está claro que los lenguajes modernos son sensitivos a las mayúsculas así que FoxLite tiene que seguir el mismo estándar, esto no es un capricho sino que más bien es para favorecer la escritura y la legibilidad del código.


