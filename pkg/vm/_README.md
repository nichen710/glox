# Glox VM

## Responsabilidad

La VM tiene la responsabilidad de ejecutar el bytecode contenido en un chunk utilizando una arquitectura basada en pila y marcos de llamada. Administra variables globales y locales, invocación de funciones de usuario y nativas, upvalues abiertos/cerrados para closures y el reporte de errores en tiempo de ejecución.

## Diferencias de implementacion con PlOXB

- Organización modular por responsabilidad: la lógica de la VM está dividida en archivos cohesivos, `vm.go` para el loop principal, `frames.go` para `CallFrame` y llamadas, `upvalues.go` para closures, `ops.go` para operaciones, `stack.go` para la pila y `natives.go` para funciones nativas.
- Valores tipados con `value.Value`: reemplaza tipos dinámicos genéricos por un struct etiquetado por `Kind`, evitando el boxing y haciendo seguras las comparaciones de igualdad.
- Upvalues seguros ante realocaciones en Go: los upvalues abiertos referencian el índice del stack en lugar de un puntero directo al slice, evitando punteros inválidos si el stack crece con `append`.
- Funciones, closures, nativas y optimización de llamadas tail: soporta `CallFrame`s anidados, funciones nativas, `OP_CONSTANT_LONG` (16 bits) y reutilización del frame actual en `OP_TAIL_CALL` cerrando previamente los upvalues abiertos.
- Despacho directo sin closures en Go: las operaciones aritméticas y de comparación evitan crear closures anónimos por instrucción, y los errores de runtime informan tanto el offset como la línea fuente.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Impresión y evaluación de literales, operadores aritméticos, comparaciones y concatenación
- Variables globales y locales
- Control de flujo
- Funciones, recursión, retorno implícito y funciones nativas
- Closures con contadores independientes, upvalues multinivel, mutación compartida y cierre al salir de bloque
- Tail Call Optimization en recursión profunda y con upvalues
- Errores de runtime (división/módulo por cero, tipos inválidos, aridad, stack underflow, stack overflow, número de línea)
