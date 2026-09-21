# Glox Scanner

## Responsabilidad

El Scanner tiene la responsabilidad de transformar el codigo fuente de LOX en tokens internos y reportar error de sintaxis basico.

## Diferencias de implementacion con PLOX

- Campo '\_error': en esta implementacion usamos el campo '\_error' para almacenar el error encontrado durante el escaneo en vez de lanzar una excepcion mediante raise.
- Metodo Scan(): en esta implementacion el metodo Scan() no solo devuelve los tokens, sino tambien retorna el error encontrado. Esto permite detener el scaneo al momento de detectar el error y no procesar el restante codigo fuente. Ademas es la forma idiomatica en Go para manejar errores.
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
