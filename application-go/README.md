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
## Que hace?
### `go mod vendor`
El comando **`go mod vendor`** copia todos los paquetes y dependencias que tu proyecto de Go necesita a una carpeta local llamada **`vendor`**.

---

#### ¿Por qué se usa `go mod vendor`?

Su propósito principal es garantizar que tu proyecto pueda ser **construido y ejecutado de manera reproducible y sin conexión a internet**.

Normalmente, el comando `go build` descarga las dependencias del proyecto desde internet y las almacena en un caché global. Sin embargo, si usas `go mod vendor`, la herramienta `go` se configurará para buscar esas dependencias exclusivamente en la carpeta `vendor` de tu proyecto.

Esto es útil para:
* **Asegurar compilaciones consistentes**: Todos los desarrolladores y el sistema de compilación usan exactamente las mismas versiones de las dependencias.
* **Entornos de red restringidos**: Si la máquina de compilación no tiene acceso a internet, el proyecto puede construirse sin problemas.

En resumen, `go mod vendor` "empaqueta" todas las dependencias de tu proyecto en un solo lugar, para que no dependas de una conexión externa para compilarlas.
