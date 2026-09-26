# Glox Interpreter

## Responsabilidad

El Interpreter tiene la responsabilidad de evaluar las expresiones del Arbol de Sintaxis Abstracta (AST) y calcular su valor en tiempo de ejecucion segun la semantica del lenguaje LOX, reportando errores de tipos y runtime.

## Diferencias de implementacion con PLOX

- Evaluacion directa por tipo: en esta implementacion usamos un switch de tipos para identificar que expresion estamos evaluando (literal, agrupamiento, unaria o binaria) y resolver su valor de forma directa y sencilla.
- Manejo de errores en runtime: en vez de cortar la ejecucion lanzando excepciones con raise, retornamos el resultado junto a un error (`RuntimeError`). Esto nos permite indicar claramente en que linea ocurrio el fallo (por ejemplo al dividir por cero o al mezclar tipos invalidos) y manejarlo de forma prolija.
- Formateo de resultados: agregamos la funcion `Stringify` para mostrar los valores en pantalla tal como los define Lox, evitando por ejemplo que los numeros se impriman con decimales de mas.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Evaluacion de literales y agrupamientos
- Operadores unarios (- y !)
- Operadores aritmeticos (+, -, \*, /, %)
- Concatenacion de cadenas
- Operadores de comparacion e igualdad
- Evaluacion de verdad (truthiness)
- Manejo de errores en runtime (tipos incompatibles, division y modulo por cero)
