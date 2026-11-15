# Glox Compiler

## Responsabilidad

El Compiler tiene la responsabilidad de compilar en una única pasada (_single-pass_) el flujo de tokens generado por el Scanner directamente a bytecode (`*chunk.Chunk`), sin construir un Árbol de Sintaxis Abstracta intermedio. Resuelve la precedencia de expresiones mediante un Pratt parser, administra scopes léxicos, funciones y upvalues para closures.

## Diferencias de implementacion con PLOXB

- Extensibilidad con `StatementCompiler`: en lugar de concentrar todos los statements en un único método o archivo, se aplica el mismo patrón declarativo que en parser (`StatementCompiler`), desacoplando cada sentencia en su propio archivo.
- Manejo de errores idiomático en Go: los errores de compilación se retornan explícitamente mediante `CompileError` (con token y línea) en vez de lanzar excepciones.
- Soporte completo de funciones, closures y tail calls: se implementó soporte completo para funciones,closures y tail calls compilando funciones con compiladores hijos, resolviendo upvalues locales y transitivos, cerrando variables capturadas al salir de bloque y convirtiendo llamadas en posición de retorno a `OP_TAIL_CALL`.
- Superinstructions de comparación y constantes largas: `!=`, `>=` y `<=` emiten una única instrucción, y los literales con índice `>= 256` emiten automáticamente `OP_CONSTANT_LONG` (hasta 65.536 constantes por chunk).

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Compilación de sentencias, literales, operadores unarios, binarios y superinstructions
- Precedencia y asociatividad del Pratt parser
- Variables globales y locales
- Control de flujo
- Declaración de funciones, llamadas, recursión y closures con upvalues
- Deduplicación de constantes, `OP_CONSTANT_LONG` y límites
- Errores de compilación
