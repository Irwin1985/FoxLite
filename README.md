# FoxLite
Mis pensamientos acerca del diseño que le daría al lenguage de programación **Foxlite**.

# Introducción
Hola, mi nombre es Irwin y soy el creador de **FoxLite**, un lenguaje de programación que pretende ser una evolución de Fox. He programado en **Visual FoxPro** durante 12 años y le tengo un gran aprecio pues me ha dado de comer por durante todo este tiempo y lo menos que puedo hacer por él es evolucionarlo. Tal vez no te guste **FoxLite**, pero créeme que el diseño que vas a ver aquí es lo más parecido a un **Visual FoxPro** evolucionado, así que comencemos a ver lo que deja, lo que hereda y lo que mejora. 

Este diseño es meramente subjetivo así que la evolución de **Fox** está en mi mano, pero no te preocupes porque conozco la gramática de **Fox** y la gramática de los lenguajes *"modernos"* por lo tanto me aseguraré que sus nuevos poderes lo ayuden a sobrevivir en este nuevo ecosistema.

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

Lo anterior es apenas una porción de la sintaxis de **Visual FoxPro**, puede hacer más cosas, pero ustedes son **Foxeros** así que ya tienen una idea de lo que es capaz de hacer.

**FoxLite** heredará gran parte de la semántica de **FoxPro** pero también es verdad que perderá parte de la sintaxis para favorecer su modernidad.

Veamos algunas de los aspectos que **FoxLite** no obtendrá de **FoxPro**:

## Limando asperezas en el proceso evolutivo

- **Literales Booleanos:** honestamente **.T.** y **.F.** no me molestan pero he decidido tender una rama hacía el árbol **ALGOL** por lo que ahora serán **true** y **false**. La vieja versión es incluso mejor, ya que escribes menos *(3 letras en lugar de 4)*, pero esto lo hago para que **FoxLite** sea bien visto por la comunidad cuya raíz desde luego parte de **ALGOL**.

```Javascript
   lbVerdad = true
   lbMentira = false
```

- **Literales Arrays:** debo admitir que la sintaxis **DIMENSION arrayName()** nunca me gustó así que esta exclusión es personal *(al igual que todas las demás)* 😋. **FoxLite** adoptará la forma simplista que muchos lenguajes tienen, pero que fue popularizada por **Javascript**.

```Javascript
   laNumeros = [1, 2, 3]
   laFrutas = ["Manzana", "Mango", "Mora"]
```

- **Sensibilidad a las Mayúsculas:** está claro que los lenguajes modernos son sensitivos a las mayúsculas así que FoxLite tiene que seguir el mismo estándar, esto no es un capricho sino que más bien es para favorecer la escritura y la legibilidad del código. **Visual FoxPro** es insensible a las mayúsculas y eso tiene sus ventajas, pero también es verdad que le resta legibilidad al tener todo el código en mayúsculas *(que suele ser lo más habitual)* lo cual genera pequeños problemas, por ejemplo, siempre escribo en minúsculas, pero cuando trabajo con compañeros siempre me terminan formateando el código a mayúsculas con la utilidad **Beautify** y tengo que volver a llevar mi código a minúsculas para que luego me lo vuelvan a reformatear.

Los lenguajes modernos incluso vienen con una utilidad integrada para formatear el código fuente, cosa que me parece estupenda y la tendré en mente para incluírsela a **FoxLite** en su versión **Consola**.

- **Procedimientos:** los procedures serán remplazados por las funciones, ya que nosotros los usamos indistintamente en **FoxPro**, ahora en **FoxLite** quiero dejar solo las funciones y otra cosa muy importante es que toda función retorna un valor implícita o explícitamente, la forma implícita es la última expresión de su bloque que será retornada, la forma explícita desde luego es con la palabra reservada **return**.

- **String multilínea con TEXT/TO:** aunque particularmente me gusta usarlo, el **TEXT TO** será remplazado por una versión más simplista inspirada por [vLang](https://vlang.io/).

- **Comentarios:** los comentarios en **FoxPro** tampoco me gustaron mucho así que **FoxLite** tendrá otros símbolos para los comentarios.

## Lo nuevo de FoxLite

No se le puede llamar *lenguaje moderno* sin que tenga características modernas verdad?, entonces vamos a ver algunas de las cosas que nos ofrecerá este lenguaje.
 
- **Notación Húngara estricta:** esta es quizá la idea más loca que se me ha ocurrido para **FoxLite**, pero si la estudiamos un poco de seguro diremos *"ah, pues claro, tiene sentido"*. Si vienes de **FoxPro** de seguro has visto o usado la [Notación Húngara](https://es.wikipedia.org/wiki/Notaci%C3%B3n_h%C3%BAngara), es aquella forma de declarar los identificadores de un programa según su ámbito y tipo, ejemplo: **lcNombre** donde *l* es de **LOCAL**, *c* es de **Character** y *Nombre* es la descripción del identificador. Entonces *¿Cómo encaja esa notación en **FoxLite**?*

Para ahorrarnos el trabajo de tener que escribir **LOCAL, PRIVATE o PUBLIC** vamos a usar la **Notación Húngara** cómo forma estricta de declaración de variables. Esto quiere decir que habrá un estilo único de escritura en el lenguaje lo cual es fantástico porque favorecerá la legibilidad y uniformidad del código. 

## Excepciones en la Notación Húngara
- El bucle **For** puede saltarse la notación húngara para favorecer la corta declaración de los iteradores.
```Javascript
    // Válido pero no recomendado
    For lcContador = 1 To 10
        ?lcContador
    
    // Bueno
    For i = 1 To 10
        ?i
```
- Toda variable declarada sin notación húngara será tratada como local.
  
- **Declaración de Variables:** gracias a la *Notación Húngara*, las variables no tiene por qué llevar delante las palabras reservadas **LOCAL, PRIVATE o PUBLIC**. El *enlace* se realizará con las primeras 2 letras seguidas de la descripción del identificador.

Veamos unos ejemplos:

```Javascript
    lcNombre = "Juan" // lo mismo que LOCAL 
    pnSaldo = 3.000 // lo mismo que PRIVATE
    glConfigFile = "c:\Congif.json" // lo mismo que PUBLIC
    
    // Para el caso de parámetros
    Func Sumar(tnNumero1, tnNumero2)
        Return tnNumero1 + tnNumero2
```

Lo anterior deja en evidencia que en **FoxLite** no habrá declaraciones de variables sin su respectiva asignación, por lo tanto toda variable que declares deberá llevar su respectivo valor para que el enlace interno sepa su ámbito, tipo y valor inicial.

- **Constantes:** **FoxLite** no tendrá *constantes simbólicas* como las tiene Fox y que no estoy pensando en un pre-procesado del código fuente. Lo que si va a tener son *constantes declaradas* y tendrán la siguiente sintaxis.

```Javascript
    const PI = 3.14159265
    lnRadio = 4
    ?"La circunferencia es: ", PI * Sqrt(lnRadio)
```

- **Los espacios en blanco importan:** mi meta con **FoxLite** es hacer que se escriba la menor cantidad de código posible, esto lo digo especialmente por aquellas palabras reservadas que todo lenguaje necesita para poder parsear el código. **FoxPro** a mi parecer tiene mucho de esto, sobre todo en las palabras de cierre como **ENDIF, ENDDO, ENDFUNC, ENDCASE, ENDFOR, etc**.

Mi idea es seguir el mismo camino que [**Python**](https://en.wikipedia.org/wiki/Python_syntax_and_semantics) y utilizar la indentación con espacios para eliminar esas palabras que al final nos estorban porque no son código ejecutable sino más bien una guía para el **Parser**.

- **Estilo de escritura en CamelCase:** para acompañar el estilo húngaro, el código de **FoxLite** tanto para *identificadores* como *palabras reservadas*, seguirá el estilo [Camel Case](https://es.wikipedia.org/wiki/Camel_case) el cual consiste en elevar a mayúsculas la primera letra del identificador, sobre todo para las descripciones compuestas en cuyo caso cada primera letra se debe elevar a mayúsculas también.

```Javascript
    liEmpleado = CreateObject("Collection")
    liEmpleado["nombre"] = "Juan"
    liEmpleado["salario"] = 3.500
    liEmpleado["horario"] = ["Lunes", "Miercoles", "Viernes"]
    
    // Imprimir los datos de un diccionario
    Func ImprimeEmpleado(tiEmpleado)
        For k, v in tiEmpleado
            ?k, v
```

- **Funciones:** aunque técnicamente las funciones no son nuevas porque **FoxPro** también las tiene, en **FoxLite** serán tratadas como [ciudadanas de primera clase](https://en.wikipedia.org/wiki/First-class_function).

Veamos algunos ejemplos:

- **Closures:** esta es quizá una de las características más ambiciosas, se trata de crear funciones dentro de otras funciones. Esto al principio puede ser confuso, pero una vez que las conozcan verán el poder que pueden ofrecer.

```Javascript
    // Función externa
    Func HolaMundo()
        pcSaludo = "Hola"
        // Función interna (closure)
        Func DecirMundo()
            Return pcSaludo + " Mundo!"
        Return DecirMundo
    // Invocar la función
    ?HolaMundo() // Imprime => "Hola Mundo!"
```

- **Funciones de alto orden:** esto es básicamente tratar a las funciones como al resto de los tipos de datos, es decir, que se puedan declarar como variables, pasar como argumentos y retornar desde otras funciones.

```Javascript
    // Eleva al cuadrado
    Func Cuadrado(tnNum)
        Return tnNum * tnNum

    // Eleva al cubo (recibe una función)
    Func Cubo(tnNum, tfCuadrado)
        Return tnNum * tfCuadrado(tnNum)

    // Invocar la función
    ?Cubo(3) // 27
    
    // Asignar una función a una variable
    lfCuadrado = Cuadrado()
    ?lfCuadrado(2) // 4
```

- **Diccionarios:** como vimos en la sintaxis de **FoxPro**, se pueden crear diccionarios o *Collection*, pero son un poco verbosas. **FoxLite** tratará los diccionarios de una forma más fácil y entendible.

```Javascript
// declarar el diccionario
datos = createobject("Collection")
datos["nombre"] = "Juan"
datos["edad"] = 36
datos["casado"] = false

// imprimir los datos
?datos["nombre"]
?datos["edad"]
?datos["casado"]
```

El ejemplo anterior nos revela que la creación de objetos a través del builtin **CreateObject** se mantendrá vigente.

- **Arrays:** los arrays también cambiarán su estructura a la forma moderna y creo que la mejor parte con respecto a los arrays de **FoxPro** es que también serán tratados como ciudadanos de primera clase.

```Javascript
// declarar el array
laFrutas = ["Manzana", "Mango", "Melocotón"]

// imprimir el array
?laFrutas[0]
?laFrutas[1]
?laFrutas[2]
```

- **String Multilínea:** un string se delimita por sus comillas simples o dobles, pero también existe otro delimitador llamado **backtick**, veamos un ejemplo:

```Javascript
lcString1 = 'string 
con
comillas 
simples'

lcString2 = "string
con
comillas
dobles"

lcString3 = `string
con
el
delimitador
backtick
`
```

- **Interpolación de string:** también se solía hacer con TEXT TO usando los dobles ángulos ```<<variable>>``` pero ahora se hará de una manera más sencilla.

```Javascript
lcNombre = "juan"
lcApellido = "perez"
? "Hola, mi nombre es $lcNombre y mi apellido $lcApellido."
```

- **JSON Nativo:** con dos funciones nativas ya podremos serializar y deserializar objetos.

```Javascript
// declarar el diccionario
liDatos = CreateObject("Collection")
liDatos["nombre"] = "Juan"
liDatos["edad"] = 36
liDatos["casado"] = false

// convertir a string JSON 
?JsonToStr(liDatos)

lcData = `
{
    "nombre": "Juan",
    "apellido": "Gonzalez",
    "edad": 36
}
`
loData = StrToJson(lcData) // bastante simple verdad?

```

- **HTTP Nativo:** las peticiones web serán tan sencillas como esto:
```Javascript
lcURL = "https://github.com/Irwin1985/FoxLite/blob/master/README.md"
?Http(lcURL)
```

- **Código Diferido:** es un código que se ejecuta al final de cada bloque de instrucciones de una función.

```Javascript
Func CargarFichero(tcFileName)
    lnHandle = fOpen(tcFileName)
    defer 
        fClose(lnHandle)

    while !fEof(lnHandle)
        ?fGets(lnHandle)
```

- **Funciones variádicas:** son las funciones que reciben 1 o más argumentos.

```Javascript
Func Saludar(taPersonas...)
    for p in taPersonas
        ?p
```

