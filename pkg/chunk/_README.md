# Glox Chunk

## Responsabilidad

El paquete Chunk tiene la responsabilidad de representar y almacenar las secuencias de bytecode compilado, su pool de constantes, los códigos de operación y el mapeo de líneas de código fuente para el reporte de errores en tiempo de ejecución y el desensamblado.

## Diferencias de implementacion con PLOXB

- Valores tipados con `value.Value`: en lugar de almacenar valores dinámicos genéricos, `Constants` utiliza `value.Value` evitando el boxing innecesario en tiempo de ejecución.
- Deduplicación de constantes: `AddConstant` reutiliza el índice si el valor ya existe en el pool de constantes del chunk, reduciendo presión sobre el límite de 256 constantes por chunk.
- Tracking de líneas por byte: `Lines []int` se mantiene sincronizado con `Code` mediante `Write` y `WriteShort`, permitiendo consultar la línea exacta con `LineAt(offset)`.
- Nuevos opcodes para optimizaciones y closures: además de las instrucciones base, se incorporan `OP_NOT_EQUAL`, `OP_GREATER_EQUAL`, `OP_LESS_EQUAL`, `OP_CALL`, `OP_TAIL_CALL`, `OP_CLOSURE`, `OP_GET_UPVALUE`, `OP_SET_UPVALUE` y `OP_CLOSE_UPVALUE`.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Escritura de instrucciones y almacenamiento de constantes
- Desensamblado de instrucciones simples, constantes y saltos relativos
- Conversión a texto de `OpCode` conocidos y desconocidos
- Deduplicación de constantes y consulta de líneas (`LineAt`)
