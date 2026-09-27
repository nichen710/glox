# Glox Environment

## Responsabilidad

El paquete Environment es el encargado de almacenar y gestionar las asociaciones (bindings) entre identificadores y valores durante la ejecucion del programa en tiempo de ejecucion. Provee operaciones para definir (`Define`), consultar (`Get`) y actualizar (`Assign`) variables tanto en el entorno global como en scopes anidados.

## Diferencias de implementacion con PLOX

- Recorrido dinámico de scopes: en `Get` y `Assign` recorremos recursivamente la cadena de punteros `enclosing` hacia arriba. Esto permite resolver variables y shadowing en tiempo de ejecución de manera autónoma, sin depender de un resolvedor estático ni distancias prefijadas.
- Manejo de errores idiomático: retornamos un `error` estándar con formato uniforme en lugar de lanzar excepciones. Esto permite al intérprete asociar el fallo al token y línea correspondientes de forma controlada.
- Constructores explícitos: se definen `NewEnvironment()` para el scope global y `NewEnclosedEnvironment(enclosing)` para scopes anidados. Esto hace explícita la intención al entrar a un bloque y asegura la inicialización del mapa interno.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Definición y lectura de variables
- Reasignación de variables existentes
- Manejo de errores al acceder o reasignar variables no declaradas
- Resolución de variables en scopes anidados y shadowing
