# Asset transfer basic sample

Ejemplo básico de transferencia de activos
El ejemplo básico de transferencia de activos demuestra:

- Conexión de una aplicación cliente (dapp) a una red blockchain de Fabric.
- Envío de transacciones del chaincode para actualizar el world-state.
- Transacciones para consultar el world-state.
- Manejo de errores en la invocación de transacciones.


### Applicacion

Siga el flujo de ejecución en el código de la aplicación cliente y el resultado correspondiente al ejecutar la aplicación. Preste atención a la secuencia de:

- Invocaciones de transacciones (salida de la consola  "**--> Submit Transaction**" y "**--> Evaluate Transaction**").
- Resultados devueltos por transacciones (salida de la consola  "**\*\*\* Result**").

### Chaincode

El código de cadena (en la carpeta `cc-assettransfer-go`) implementa las siguientes funciones para admitir la aplicación:

- CreateAssetUsingStructParam
- CreateAsset
- ReadAsset
- UpdateAsset
- DeleteAsset
- TransferAsset

Tenga en cuenta que la transferencia de activos implementada por el contrato inteligente es un escenario simplificado, sin validación de propiedad, destinado solo a demostrar cómo invocar transacciones.

## Ejecutando la app

1. Debe ejecutar la red `fabric-testnet-nano-without-syschannel` (leer fabric-testnet-nano-without-syschannel/README.md)

2. Ejecutar la aplicación
   ```bash
   cd application-go
   go mod vendor
   go build -o app-mycc

   ./app-mycc
   ```
