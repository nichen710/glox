# Glox Scanner

## Responsabilidad

El Scanner tiene la responsabilidad de transformar el codigo fuente de LOX en tokens internos y reportar error de sintaxis basico.

## Diferencias de implementacion con PLOX

- Manejo de multiples errores: en esta implementacion en vez de interrumpir la ejecucion levantando una excepcion con raise ante el primer fallo, acumulamos los errores encontrados durante el escaneo para reportar todos los problemas lexicos de una sola vez.
- Metodo Scan(): en esta implementacion el metodo Scan() no solo devuelve los tokens, sino tambien los errores encontrados, siguiendo la forma idiomatica de Go para manejar errores y permitiendo detener el procesamiento antes de avanzar al parser si hubo fallos.
- Encapsulamientos de escaneo: en esta implementacion encapsulamos los casos especiales de escaneo en metodos privados (ej. stringLiteral, numberLiteral, identifier).

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Agrupacion y puntuacion
- Operadores aritmeticos
- Comparaciones e igualdad
- Espacios, tabulaciones y saltos de linea
- Comentarios de una sola linea
- Tokens con comentarios y saltos de linea
- Cadenas
- Numeros
- Palabras clave vs identificadores
- Declaraciones completas
- Errores de escaneo
