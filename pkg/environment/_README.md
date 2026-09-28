# Glox Environment

## Responsabilidad

El paquete Environment es el encargado de almacenar y gestionar las asociaciones (bindings) entre identificadores y valores durante la ejecucion del programa en tiempo de ejecucion. Provee operaciones para definir (`Define`), consultar (`Get`) y actualizar (`Assign`) variables tanto en el entorno global como en scopes anidados.

## Diferencias de implementacion con PLOX

- Manejo de errores: en esta implementacion en vez de interrumpir la ejecucion levantando excepciones, retornamos errores tipados como valores para que el interprete pueda manejar los fallos de manera controlada siguiendo el estandar de Go.
- Separacion en las operaciones de acceso: en esta implementacion en vez de recibir distancias como parametros opcionales, separamos las consultas en metodos especificos segun si se busca en los entornos anidados o si se accede directamente a una distancia ya resuelta.
- Constructores explícitos: se definen `NewEnvironment()` para el scope global y `NewEnclosedEnvironment(enclosing)` para scopes anidados. Esto hace explícita la intención al entrar a un bloque y asegura la inicialización del mapa interno.
- Resolucion en entornos anidados: en esta implementacion recorremos los entornos envolventes para resolver o reasignar variables en los scopes padres cuando se busca de manera dinamica.

## Test

Para las pruebas se trataron de probar los casos de interes mediante test suites. Los casos que se probaron son:

- Definición y lectura de variables
- Reasignación de variables existentes
- Manejo de errores al acceder o reasignar variables no declaradas
- Resolución de variables en scopes anidados y shadowing
