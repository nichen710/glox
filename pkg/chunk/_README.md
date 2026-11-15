# Glox Chunk

## Responsabilidad

El paquete Chunk tiene la responsabilidad de representar y almacenar las secuencias de bytecode compilado, su pool de constantes, los códigos de operación y el mapeo de líneas de código fuente para el reporte de errores en tiempo de ejecución y el desensamblado.

## Diferencias de implementacion con PLOXB

- Valores tipados con `value.Value`: en lugar de almacenar valores dinámicos genéricos, `Constants` utiliza `value.Value` evitando el boxing innecesario en tiempo de ejecución.
- Deduplicación de constantes en O(1): se utiliza un mapa interno para reutilizar el índice si el valor ya existe en el pool del chunk.
- Soporte para hasta 65.536 constantes: cuando el índice supera 255, se emite `OP_CONSTANT_LONG` con operando de 16 bits, manteniendo `OP_CONSTANT` para los primeros 256 valores.
- Tracking de líneas por byte: `Lines []int` se mantiene sincronizado con `Code` mediante `Write` y `WriteShort`, permitiendo consultar la línea exacta con `LineAt(offset)`.
- Nuevos opcodes para optimizaciones y closures: además de las instrucciones base, se incorporan `OP_CONSTANT_LONG`, `OP_NOT_EQUAL`, `OP_GREATER_EQUAL`, `OP_LESS_EQUAL`, `OP_CALL`, `OP_TAIL_CALL`, `OP_CLOSURE`, `OP_GET_UPVALUE`, `OP_SET_UPVALUE` y `OP_CLOSE_UPVALUE`.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Escritura de instrucciones y almacenamiento de constantes
- Desensamblado de instrucciones simples, constantes y saltos relativos
- Conversión a texto de `OpCode` conocidos y desconocidos
- Deduplicación de constantes y consulta de líneas (`LineAt`)
