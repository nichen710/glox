# Glox Resolver

## Responsabilidad

El Resolver tiene la responsabilidad de realizar el análisis semántico estático sobre el Árbol de Sintaxis Abstracta (AST), resolviendo el alcance léxico de las variables y detectando errores semánticos antes de la ejecución.

## Diferencias de implementacion con PLOX

- Pipeline desacoplado del intérprete: en esta implementacion el Resolver es independiente y produce una `BindingTable` que se le inyecta luego al intérprete.
- Arquitectura extensible con `PipelineStep` y `ResolverBuilder`: en esta implementacion se diseñó la interfaz genérica `PipelineStep` y un `ResolverBuilder` declarativo que permite encadenar futuras pasadas de análisis o linter dentro del pipeline de forma modular.
- Identidad de nodos mediante identificadores únicos (`ID`): en esta implementacion se asigna un identificador a cada nodo relevante para indexar la tabla de bindings de forma directa y eficiente.
- Validación de sentencias `return` y contexto de POO: se rastrea el contexto de funciones, inicializadores, clases y subclases para validar `return`, `this`, `super` y evitar auto-herencia.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Resolución de profundidad de variables locales y encadenamiento de scopes
- Resolución de parámetros de funciones y clausuras (closures)
- Detección de auto-inicialización en variables locales (`var a = a;`)
- Detección de variables duplicadas en el mismo ámbito local
- Prohibición de retornos a nivel raíz (top-level return)
- Errores semánticos de POO (`this` fuera de clase, `super` sin superclase, auto-herencia, retorno de valor en `init`)
- Ejecución y extensibilidad del pipeline con pasos personalizados
