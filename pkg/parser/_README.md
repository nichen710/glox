# Glox Parser

## Responsabilidad

El Parser es el siguiente paso lógico despues de el scaneo. Los tokens son ordenados en frases con sentido (o si no tienen sentido se devuelve un error). El parser arma un Árbol de Sintaxis Abstracta (AST) que ordena el input en grupos lógicos y con precedencia. Este árbol se arma utilizando el recorrido de top-down parser visto en clases, donde los nodos son las expresiones y sentencias (statements) y la lógica de armado garantiza que se sigan las reglas de nuestra gramática.

## Diferencias de implementacion con PLOX

- Manejo de errores adaptado al standard de Go. Como en el parser, el error se almacena en un campo en lugar de lanzarlo con raise. Esto introdujo la necesidad de checkeos en las llamadas recursivas de Expression para que, en caso de error, se deje de construir el árbol inválido.
- Abstracción de operadores lógicos y binarios: uso de `parseBinary` y `parseLogical` para desacoplar y evitar duplicación de código en la jerarquía de precedencia (`or`, `and`, `equality`, `comparison`, `term`, `factor`).
- Extensibilidad con `StatementFactory`: para el parseo de sentencias se utiliza una lista de fábricas que desacoplan cada tipo de sentencia en su propio archivo, cayendo por defecto en sentencias de expresión.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Expresiones literales
- Expresiones unarias
- Expresiones binarias
- Expresiones lógicas
- Expresiones de agrupamiento
- Tests de precedencia
- Tests de errores
- Sentencias
