# Glox Compiler

## Responsabilidad

El Compiler tiene la responsabilidad de compilar en una única pasada (*single-pass*) el flujo de tokens generado por el Scanner directamente a bytecode (`*chunk.Chunk`), sin construir un Árbol de Sintaxis Abstracta intermedio. Resuelve la precedencia de expresiones mediante un Pratt parser, administra scopes léxicos, funciones y upvalues para closures.

## Diferencias de implementacion con PLOXB

- Extensibilidad con `StatementCompiler`: en lugar de concentrar todos los statements en un único método o archivo, se aplica el mismo patrón declarativo que en `pkg/parser` (`StatementCompiler`), desacoplando cada sentencia (`var`, `fun`, `return`, `print`, `if`, `while`, `for`, `{ ... }`, expression statement) en su propio archivo.
- Manejo de errores idiomático en Go: los errores de compilación se retornan explícitamente mediante `CompileError` (con token y línea) en vez de lanzar excepciones.
- Soporte completo de funciones, closures y tail calls: extiende la base de `ploxb` (`main`) compilando funciones con compiladores hijos (`enclosing`), resolviendo upvalues locales y transitivos (`resolveUpvalue`), cerrando variables capturadas al salir de bloque (`OP_CLOSE_UPVALUE`) y convirtiendo llamadas en posición de retorno a `OP_TAIL_CALL`.
- Superinstructions de comparación: `!=`, `>=` y `<=` emiten una única instrucción (`OP_NOT_EQUAL`, `OP_GREATER_EQUAL`, `OP_LESS_EQUAL`) en lugar de combinar dos opcodes.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Compilación de sentencias, literales, operadores unarios, binarios y superinstructions
- Precedencia y asociatividad del Pratt parser
- Variables globales y locales (shadowing y reutilización de slots)
- Control de flujo (`if`/`else`, `while`, `for`, `and`, `or` con cortocircuito)
- Declaración de funciones, llamadas, recursión y closures con upvalues
- Deduplicación de constantes y límites (`MaxConstants`, `MaxLocals`, `MaxJump`)
- Errores de compilación (targets inválidos, auto-inicialización, `return` en top-level)

