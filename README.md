# glox

Interprete de Lox en Go.

## Grupo

- **Nicolas Chen** 105907 - [nichen710](https://github.com/nichen710)
- **Mateo Cabrera Rodríguez** 108118 - [m-cabrerar](https://github.com/m-cabrerar)

## Consigna

Como parte práctica de la materia Lenguajes y Compiladores se desarrolla el presente trabajo práctico. Consiste en implementar un lenguaje de programación. El lenguaje a implementar será Lox, un lenguaje de programación creado para el libro [Crafting Interpreters](https://craftinginterpreters.com). Lox es un lenguaje de alto nivel, con tipado dinámico. Su sintaxis está inspirada en C.

Esta implementación cuenta con las siguientes funcionalidades:

- Booleanos, números, strings, arrays y nils.
- Calculos ariméticos simples (+, -, \*, /, %) y lógica booleana.
- Control de flujo (if, for, while).
- Declaración de funciones, closures y variables.
- Programación orientada a objetos completa.

## Ejecución

Para utilizar el interprete, primero debe ser compilado:

```sh
go build -o ./glox ./cmd/glox
```

Para la ejecución:

```sh
./glox
```

También se puede ejecutar un archivo con los comandos ya escritos en lox:

```sh
./glox script.lox
```

Para probar las diferentes capas de procesamiento del lenguaje o elegir el motor de ejecución, se pueden utilizar los siguientes modos:

- `-scanning`: se transforma el código fuente en tokens.
- `-parsing`: los tokens se ordenan en sentencias y expresiones, respetando la gramática de Lox. Estas expresiones se ordenan en un Arbol de Sintaxis Abstracta, que respeta los órdenes de precedencia.
- `-resolve`: recorre el Arbol y calcula a qué scope corresponde cada variable.
- `-compile`: compila el código fuente directamente en una única pasada y muestra el bytecode generado (ya optimizado).
- `-vm`: ejecuta el programa compilando a bytecode y corriéndolo en la máquina virtual, en lugar de usar el intérprete Tree-Walk por defecto.
- `-no-opt`: omite las pasadas de optimización sobre el bytecode.

Por ejemplo:

```sh
./glox -scanning
./glox -vm script.lox
```

#### Ejemplo de uso real

Como prueba práctica de las capacidades interactivas del lenguaje, se incluye un juego de adivinar el número en `examples/guess_number.lox` que utiliza funciones, control de flujo y funciones nativas de consola. Se puede jugar en cualquiera de los dos motores:

```sh
./glox examples/guess_number.lox
./glox -vm examples/guess_number.lox
```

## Estructura del proyecto

#### `cmd/glox`

Es el entrypoint. Acá está el `main.go` que se encarga de punto de entrada y salida de información al sistema mediante la consola. Crea y ejecuta las distintas etapas del intérprete.

#### `pkg/`

Acá se encuentra la lógica del sistema, aislada en distintos paquetes con sus responsabilidades y testeos.

## Testing

Además de los tests de la cátedra, cada paquete con lógica tiene su archivo con tests. También se agregaron tests de volúmen para calcular benchmarks comparables con otras implemtenaciones, en `./benchmark-tests`

## Implentación y diferencias con Plox

La implentación se realizó con lo aprendido en clases y respetando el lenguaje definido en el libro. Al ser implementado en Go, encontramos ciertas particularidades a destacar.

- No hay excepciones para errores esperados. Los errores los manejamos como valores retornados. Esto nos obligó a pensar claramente el camino de los errores en cada función.
- Por la misma razón, el return de Lox está implementado retornando un error. Es parecido a la implentación del libro y de Plox, pero estos lanzan un error que se puede catchear, Go retorna una ReturnSignal con pinta de error.
- Como el error no corta la ejecución, sino que es devuelto, esta implementación del scanner puede terminar de leer el código inválido y retornar todos los errores encontrados en lugar de tan solo el pimero.
- IDs explícitos para los nodos del AST. En plox la clave de identificación era el objeto. En Go, los structs (que representan un nodo) se pasan por valor, no por puntero, por lo que surgió la necesidad de utilizar identificadores en la implementación.
- Utilizamos fábricas de statements. Para tener un codigo más prolijo, todos los statement tienen su propio archivo, en lugar de estar todos en el parser. Esto emprolija el código y nos facilita el agregar statements en el futuro.
- Go es un lenguaje estático, a diferencia de Lox y Python. Esto nos crea la necesidad de siempre checkear los tipos de las variables antes de utilizarlas.

## Benchmarks

Se realizaron tests de volumen para evaluar distintos aspectos. Recursividad, costo de sumas aritméticas y costo de resolución de variables en closures.

Para cada aspecto, se corrió una prueba 5 veces en cada interprete para comparar los resultados, se reporta el promedio de esos 5 valores. El volúmen se acotó para poder correr las pruebas en plox.

#### Interpretes evaluados

- Glox: nuestra implementación de Lox en Go.
- [Plox](https://github.com/FdelMazo/plox): la implementación de la cátredra en python.
- [Clox](https://github.com/munificent/craftinginterpreters): la implentación oficial del libro en C.
- [Otro Glox](https://github.com/chidiwilliams/glox): otra implementación pública de glox, por Chidi Williams.

| Benchmark | glox  | clox  |  plox  | other glox |
| :-------: | :---: | :---: | :----: | :--------: |
| closures  | 0.41s | 0.04s | 33.89s |   7.92s    |
| recursion | 0.37s | 0.02s | 16.19s |   2.00s    |
|   sums    | 0.26s | 0.02s | 24.00s |   5.23s    |

Para evitar estar viendo una diferencia en overhead también se ejecturon las mismas pruebas pero con mayor volumen para clox y glox:

| Benchmark |  glox  | clox  |
| :-------: | :----: | :---: |
| closures  | 3.72s  | 0.46s |
| recursion | 1.22s  | 0.08s |
|   sums    | 24.66s | 2.78s |

De estos resultados podemos observar que plox es la implementación más lenta de las 3. Go al ser un lenguaje compilado es mucho más eficiente que python. Plox es un interprete siendo ejecutado por un interprete. Además, las llamadas a funciones son caras en Cpython. Por otro lado, Clox es un interprete que corre como código nativo, ya está compilado.
También podemos ver que la otra implementación de Lox en Go que encontramos es más lenta que la nuestra. Para entender por qué vimos el código en los puntos donde nosotros tomamos decisiones importantes (ya mencionadas antes) y notamos lo siguiente:

- Donde nosotros utilizamos IDs para identificar los nodos, la otra implementación convierte todo el nodo a texto y utiliza eso como clave. Cualquier lectura o asignación a una variable es más lenta.
- Donde nosotros devolvemos errores de la manera intencionada en Go, la otra solución utiliza panic y recover, lo cual es costoso. Esto hace a los returns en Lox lentos.

Estas diferencias en benchmark son una clara diferencia de Lenguaje y Implementación. Todas son el mismo lenguaje, Lox, implementado de maneras distintas con diseños distintos que tienen consecuencias en performance, pero no afectan los resultados de las ejecuciones.

#### Diferencias con otras implementaciones en Go

Se realizaron pruebas y se logró comprobar empíricamente las pricipales causas de la diferencia en tiempos de ejecución con la implementación de Chidi Williams.

Se ejecutaron nuevamente los benchmarks utilizando el profiler de Go `pprof`. En la implementación más lenta, se observó que entre el 60% y 80% del tiempo de ejecución está siendo utilizado por `fmt.Sprintf`, que es lo que se utiliza para armar la clave de cada variable leida o asignada.

Otra comprobación realizada fue modificar las partes que nos parecen podrían ser el problema. En la implementación de Chidi corregimos el problema de la clave de texto, eso hizo que se ejecutara mucho más rápido. Aún así, en los tests que llaman a funciones seguía performando considerablemente peor que nuestra implementación. Luego cambiamos el uso de panics por returns para devolver como valor, lo que mejoró la llamada a funciones y acercó bastante los tiempos de ejecución entre ambas implementaciones.

| Benchmark | other glox | other glox modificado | glox  |
| :-------: | :--------: | :-------------------: | :---: |
| closures  |   6.25s    |         0.67s         | 0.38s |
| recursion |   1.91s    |         0.36s         | 0.31s |
|   sums    |   5.03s    |         0.47s         | 0.25s |

#### Tree-Walk vs Bytecode VM y efecto de la optimización

Con la incorporación del compilador a bytecode y la máquina virtual (`-vm`), comparamos ambos motores internos sobre los mismos tests de volumen. La VM resulta más rápida en llamadas a funciones, recursión y closures porque evita recorrer nodos del AST y crear entornos dinámicos en cada llamada, trabajando directamente sobre una pila contigua y frames planos.

| Benchmark | Tree-Walk | Bytecode VM (`-vm`) |
| :-------: | :-------: | :-----------------: |
| closures  |  0.737s   |       0.446s        |
| recursion |  0.360s   |       0.231s        |
|   sums    |  0.328s   |       0.332s        |

Además, medimos el impacto de la pasada de optimización (_Jump Threading_, que redirige saltos encadenados directo al destino final) comparando `-vm` contra `-vm -no-opt` en dos escenarios:

- **Código con pocos saltos (`sums.lox`)**: prácticamente no hay diferencia (`0.318s` sin optimizar vs `0.322s` optimizado), ya que casi no existen cadenas de saltos que colapsar.
- **Código con múltiples saltos encadenados (`chained_jumps.lox`, cadena de 60 `else if` en un loop)**: baja de `0.172s` sin optimizar a `0.159s` optimizado, porque la VM ejecuta un único salto directo en lugar de atravesar toda la cadena de `OP_JUMP` intermedios en cada iteración. Aun así, la ganancia total no es tan grande como podría esperarse porque el pass de optimización igual debe recorrer cada chunk completo antes de la ejecución.

También se evaluó un algoritmo completo combinando POO y arrays (**QuickSort iterativo** con stack explícito, en `quicksort_150.lox`, `quicksort_10000_worst.lox` y `quicksort_go.go`) comparando Tree-Walk, Bytecode VM optimizado y Go nativo:

|         Caso QuickSort          | Tree-Walk  | Bytecode VM (`-vm`) |  Go nativo  | Comparaciones |
| :-----------------------------: | :--------: | :-----------------: | :---------: | :-----------: |
|       **150 aleatorios**        | `0.00323s` |     `0.00113s`      | `0.000022s` |     1.119     |
| **10.000 inversos (peor caso)** |  `51.43s`  |      `35.62s`       |  `0.174s`   |  49.995.000   |

## Extras agregados

Más allá del Lox base, se sumaron extensiones tanto en el diseño interno de la VM como en el lenguaje, manteniendo paridad entre el intérprete Tree-Walk y el Bytecode VM:

- **Constantes largas (`OP_CONSTANT_LONG`)**: originalmente el pool de constantes aceptaba hasta 256 valores indexados con 1 byte. Para no limitar programas más grandes ni duplicar el tamaño de todas las instrucciones, se mantuvo `OP_CONSTANT` de 1 byte para los primeros 256 índices y se agregó `OP_CONSTANT_LONG` con operando de 2 bytes (hasta 65.536 constantes) cuando se supera ese umbral, junto con deduplicación automática de constantes repetidas.
- **Funciones nativas**: se incorporó un registro uniforme de funciones nativas en ambos motores para interactuar con el sistema y la consola (`clock`, `random`, `sleep`, `input`, `toNumber`, `clear`), además de las operaciones sobre estructuras de datos.
- **POO completo**: se implementaron clases, instanciación, lectura y escritura de propiedades, métodos enlazados con `this`, constructores, herencia simple y llamadas a la superclase, incluyendo las validaciones estáticas correspondientes en el resolver y en el compilador.
- **Arrays y operaciones ABM**: se agregaron literales de array (`[...]`), lectura por índice y asignación por índice con validación de índice entero y rango. Los arrays se representan por referencia para que puedan mutarse al pasarse a funciones o métodos, y se acompañan de cinco funciones nativas básicas: `len` (longitud de arrays o strings), `push` (agregar al final), `pop` (quitar del final), `insert` (insertar en una posición) y `remove` (eliminar por índice).

## Anexo: Dificultades y trade-offs

- **Optimización y cálculo de offsets (_Jump Threading_)**: al agregar una etapa de optimización sobre bytecode plano, cualquier opcode nuevo con operandos requiere mantener actualizado el cálculo de tamaño de instrucciones para no desalinear offsets al recorrer el chunk. Para no invalidar saltos existentes ni el mapeo de líneas, _Jump Threading_ modifica únicamente destinos de salto del mismo tamaño, aunque recorrer cada chunk antes de ejecutar agrega un costo fijo que reduce el impacto neto en programas cortos o con pocos saltos.
- **POO en la Bytecode VM (`this`, `init` y _tail calls_)**: las funciones comunes usan el slot `@0` para su primer parámetro, mientras que los métodos e inicializadores necesitan reservar `@0` para `this`. Cambiar todos los frames habría roto la convención existente, por lo que se mantuvo el layout original en funciones regulares y se reservó `@0` solo en métodos/inicializadores, guardando en cada frame la posición de retorno para que tanto `return` en `init` como las optimizaciones de _tail call_ dejen el stack consistente. Además, se incorporaron `OP_INVOKE` y `OP_SUPER_INVOKE` para evitar crear objetos _bound method_ intermedios en llamadas directas.
- **Arrays por referencia y límites de recursión**: en Go los slices se pasan por valor de cabecera, por lo que un `append` dentro de `push` o `insert` no sería visible para quien llama si se pasara el slice directo. Se mitigó envolviendo los elementos en un struct puntero (`LoxArray` en Tree-Walk y `value.Array` en la VM) y centralizando las validaciones de índice entero y rango. Por otro lado, al probar QuickSort sobre 10.000 elementos en el peor caso, la versión recursiva superaba el límite de 256 frames de la VM (`MaxFrames`); en lugar de inflar el límite artificialmente, se implementó la variante iterativa apoyada en un array como stack explícito.
