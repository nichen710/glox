# Glox Parser

## Responsabilidad

El Parser es el siguiente paso lógico despues de el scaneo. Los tokens son ordenados en frases con sentido (o si no tienen sentido se devuelve un error). El parser arma un Árbol de Sintaxis Abstracta (AST) que ordena el input en grupos lógicos y con precedencia. Este árbol se armar utilizando el recorrido de top-down parser visto en clases, donde los nodos son las expresiones y la lógica de armado garantiza que se sigan las reglas de nuestra gramática.

## Diferencias de implementacion con PLOX

- Manejo de errores adaptado al standar de Go. Como en el parser, el error se almacena en un campo en lugar de lanzarlo con raise. Esto introdujo la necesidad de checkeos en las llamadas recursivas de Expression para que, en caso de error, se deje de construir el árbol inválido.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Expresiones literales
- Expresiones unarias
- Expresiones binarias
- Expresiones de agrupamiento
- Tests de precedencia
- Tests de errores
