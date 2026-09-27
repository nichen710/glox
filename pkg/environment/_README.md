# Glox Environment

## Responsabilidad

El paquete Environment es el encargado de almacenar y gestionar las asociaciones (bindings) entre identificadores y valores durante la ejecucion del programa en tiempo de ejecucion. Provee operaciones para definir (`Define`), consultar (`Get`) y actualizar (`Assign`) variables.

## Diferencias de implementacion con PLOX

- Almacenamiento directo mediante mapa de Go (`map[string]any`) para un acceso O(1) rápido y eficiente a las variables.
- Manejo de errores idiomático retornando tuplas con `error` en lugar de lanzar excepciones, permitiendo al intérprete encapsular el fallo en un `RuntimeError` estructurado con la información del token y su línea.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Definición y lectura de variables
- Reasignación de variables existentes
- Manejo de errores al acceder o reasignar variables no declaradas
