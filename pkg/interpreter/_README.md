# Glox Interpreter

## Responsabilidad

El Interpreter tiene la responsabilidad de ejecutar las sentencias y evaluar las expresiones del Arbol de Sintaxis Abstracta (AST) calculando su valor en tiempo de ejecucion segun la semantica del lenguaje LOX, reportando errores de tipos y runtime.

## Diferencias de implementacion con PLOX

- Manejo de errores en runtime: en esta implementacion en vez de interrumpir la ejecucion levantando excepciones, retornamos errores tipados como valores para reportar con precision el token y la linea del fallo siguiendo el estandar de Go.
- Control de retorno con señales: en esta implementacion en vez de lanzar excepciones para retornar de una funcion, propagamos una señal de retorno tipada que es capturada al invocar la llamada.
- Salida configurable: en esta implementacion permitimos configurar el destino de salida del interprete, facilitando la verificacion y captura de salidas en las pruebas unitarias.

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
