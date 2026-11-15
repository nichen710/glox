# Glox Optimizer

## Responsabilidad

El Optimizer tiene la responsabilidad de aplicar pasadas de optimización sobre el bytecode compilado y recursivamente sobre los chunks de las funciones almacenadas en su pool de constantes, mejorando el código generado sin alterar la semántica del programa.

## Diferencias de implementacion con PlOXB

- Fase dedicada post-compilación: en vez de ejecuta directamente el chunk salido del compilador, se incorporó una etapa intermedia de optimización antes de la ejecución o desensamblado.
- Arquitectura extensible con `Pass`: siguiendo el mismo estilo de pipeline que `pkg/resolver`, cada optimización implementa la interfaz `Pass`, permitiendo encadenar nuevas pasadas de forma modular desde `NewDefaultOptimizer()`.
- Jump Threading (`JumpThreadingPass`): detecta instrucciones `OP_JUMP` y `OP_JUMP_IF_FALSE` cuyo destino inmediato es otro `OP_JUMP` incondicional y redirige el salto directo al destino final (con protección contra ciclos), sin modificar el tamaño en bytes ni invalidar líneas.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Redirección de cadenas de saltos.
- Preservación de la semántica en estructuras de control encadenadas.
