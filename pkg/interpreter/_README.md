# Glox Interpreter

## Responsabilidad

El Interpreter tiene la responsabilidad de ejecutar las sentencias y evaluar las expresiones del Arbol de Sintaxis Abstracta (AST) calculando su valor en tiempo de ejecucion segun la semantica del lenguaje LOX, reportando errores de tipos y runtime.

## Diferencias de implementacion con PLOX

- Evaluacion directa por tipo: en esta implementacion usamos un switch de tipos tanto para sentencias (`Statement`) como para expresiones (`Expression`), ejecutando las acciones correspondientes de forma directa y sencilla.
- Manejo de errores en runtime: en vez de cortar la ejecucion lanzando excepciones con raise, retornamos el resultado junto a un error (`RuntimeError`). Esto nos permite indicar claramente en que linea ocurrio el fallo (por ejemplo al dividir por cero o al mezclar tipos invalidos) y manejarlo de forma prolija.
- Formateo de resultados: agregamos la funcion `Stringify` para mostrar los valores en pantalla tal como los define Lox, evitando por ejemplo que los numeros se impriman con decimales de mas.
- Salida configurable con `SetWriter`: permite redirigir la salida del interprete hacia un buffer para verificar la emision de sentencias `print` en las pruebas unitarias.
- Control de retornos con señales: en vez de lanzar excepciones con raise para el `return`, propagamos una señal centinela que se captura al invocar la función.
- Interfaz `Callable`: unifica bajo un mismo contrato estático tanto las funciones de usuario como las funciones nativas (ej. `clock`).

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Evaluacion de literales y agrupamientos
- Operadores unarios (- y !)
- Operadores aritmeticos (+, -, \*, /, %)
- Concatenacion de cadenas
- Operadores de comparacion e igualdad
- Operadores logicos y cortocircuito
- Evaluacion de verdad (truthiness)
- Manejo de errores en runtime (tipos incompatibles, division y modulo por cero)
- Ejecucion de sentencias
