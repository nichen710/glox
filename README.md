# glox

Interprete de Lox en Go.

## Grupo

* **Nicolas Chen** 105907 - [nichen710](https://github.com/nichen710)
* **Mateo Cabrera Rodríguez** 108118 - [m-cabrerar](https://github.com/m-cabrerar)

## Consigna

Como parte práctica de la materia Lenguajes y Compiladores se desarrolla el presente trabajo práctico. Consiste en implementar un lenguaje de programación. El lenguaje a implementar será Lox, un lenguaje de programación creado para el libro [Crafting Interpreters](https://craftinginterpreters.com). Lox es un lenguaje de alto nivel, con tipado dinámico. Su sintaxis está inspirada en C.

Esta implementación cuenta con las siguientes funcionalidades:
- Booleanos, números, strings, nils.
- Calculos ariméticos simples (+, -, *, /) y lógica booleana.
- Control de flujo (if, for, while).
- Declaración de funciones y variables.

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

Para probar las diferentes capas de procesamiento del lenguaje, se pueden utilizar los siguientes modos:
- `-scanning`: se transforma el código fuente en tokens.
- `-parsing`: los tokens se ordenan en sentencias y expresiones, respetando la gramática de Lox. Estas expresiones se ordenan en un Arbol de Sintaxis Abstracta, que respeta los órdenes de precedencia.
- `-resolve`: recorre el Arbol y calcula a qué scope corresponde cada variable.

Por ejemplo:
```sh
./glox -scanning
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

| Benchmark | glox | clox | plox  | other glox |
|:---------:|:----:|:----:|:-----:|:----------:|
| closures  |0.41s |0.04s |33.89s |    7.92s   |
| recursion |0.37s |0.02s |16.19s |    2.00s   |
| sums      |0.26s |0.02s |24.00s |    5.23s   |

Para evitar estar viendo una diferencia en overhead también se ejecturon las mismas pruebas pero con mayor volumen para clox y glox:

| Benchmark | glox  | clox |
|:---------:|:-----:|:----:|
| closures  |3.72s  |0.46s |
| recursion |1.22s  |0.08s |
| sums      |24.66s |2.78s |

De estos resultados podemos observar que plox es la implementación más lenta de las 3. Go al ser un lenguaje compilado es mucho más eficiente que python. Plox es un interprete siendo ejecutado por un interprete. Además, las llamadas a funciones son caras en Cpython. Por otro lado, Clox es un interprete que corre como código nativo, ya está compilado.
También podemos ver que la otra implementación de Lox en Go que encontramos es más lenta que la nuestra. Para entender por qué vimos el código en los puntos donde nosotros tomamos decisiones importantes (ya mencionadas antes) y notamos lo siguiente:
- Donde nosotros utilizamos IDs para identificar los nodos, la otra implementación convierte todo el nodo a texto y utiliza eso como clave. Cualquier lectura o asignación a una variable es más lenta.
- Donde nosotros devolvemos errores de la manera intencionada en Go, la otra solución utiliza panic y recover, lo cual es costoso. Esto hace a los returns en Lox lentos.

Estas diferencias en benchmark son una clara diferencia de Lenguaje y Implementación. Todas son el mismo lenguaje, Lox, implementado de maneras distintas con diseños distintos que tienen consecuencias en performance, pero no afectan los resultados de las ejecuciones.
