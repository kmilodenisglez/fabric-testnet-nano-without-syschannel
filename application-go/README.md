# Asset Transfer Basic Sample

Este ejemplo demuestra un **flujo completo de transferencia de activos** en Hyperledger Fabric usando una aplicación cliente en Go.

El propósito de este proyecto es que los estudiantes aprendan:

* Cómo **conectar una aplicación cliente** a una red blockchain de Hyperledger Fabric.
* Cómo **enviar transacciones** al chaincode para modificar el `world state`.
* Cómo **consultar el `world state`** usando transacciones de evaluación.
* Cómo **manejar errores** en la invocación de transacciones.

---

## Aplicación Cliente (Go)

La aplicación cliente (`application-go/assetTransfer.go`) sigue un flujo claro que los estudiantes deben observar. Cada línea y sección tiene un propósito educativo:

1. **Inicio de la aplicación**

   ```go
   log.Println("============ application-golang starts ============")
   ```

   Marca el inicio de la ejecución en la consola.

2. **Configuración del descubrimiento de peers**

   ```go
   os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
   ```

   Indica al SDK que todos los peers son accesibles desde `localhost` (útil para redes de prueba locales).

3. **Creación y manejo del wallet**

   ```go
   wallet, _ := gateway.NewFileSystemWallet("wallet")
   ```

   * La wallet almacena identidades digitales (certificados X.509 y claves privadas).
   * Si la identidad `appUser` no existe, la función `populateWallet` la crea usando los archivos de la carpeta `crypto-config`.

4. **Conexión al Gateway**

   ```go
   gw, err := gateway.Connect(
       gateway.WithConfig(config.FromFile("ccp.yaml")),
       gateway.WithIdentity(wallet, "appUser"),
   )
   ```

   * Conecta la aplicación a la red usando el archivo `ccp.yaml` (Common Connection Profile).
   * `appUser` es la identidad que firmará las transacciones.

5. **Obtención del Network y del Contract**

   ```go
   network := gw.GetNetwork(channelName)
   contract := network.GetContract(contractName)
   ```

   * `network` representa el canal de Fabric donde se encuentra el ledger.
   * `contract` permite invocar funciones del chaincode (`basic` en este caso).

6. **Transacciones de ejemplo**

   * **InitLedger**: Inicializa el ledger con activos de ejemplo.

     ```go
     contract.SubmitTransaction("InitLedger")
     ```
   * **GetAllAssets**: Consulta todos los activos.

     ```go
     contract.EvaluateTransaction("GetAllAssets")
     ```
   * **CreateAssetUsingStructParam**: Crea un activo a partir de una estructura en Go.
   * **CreateAsset**: Crea un activo con parámetros individuales.
   * **ReadAsset**: Obtiene los detalles de un activo específico.
   * **AssetExists**: Verifica si un activo existe.
   * **TransferAsset**: Cambia el dueño de un activo existente.

7. **Manejo de errores**
   Cada invocación verifica si ocurre un error y lo reporta:

   ```go
   if err != nil {
       log.Fatalf("Failed to Submit transaction: %v", err)
   }
   ```

   Esto permite a los estudiantes ver cómo capturar problemas al invocar transacciones.

---

## Chaincode

El chaincode (`cc-assettransfer-go`) implementa funciones que corresponden a las transacciones invocadas por la aplicación:

* `CreateAssetUsingStructParam`
* `CreateAsset`
* `ReadAsset`
* `UpdateAsset`
* `DeleteAsset`
* `TransferAsset`

> Nota: La lógica de transferencia de activos está simplificada. No hay validación de propiedad ni reglas complejas, solo se busca **demostrar cómo invocar transacciones desde una aplicación cliente**.

---

## Ejecutando la Aplicación

1. **Levantar la red de prueba**

   ```bash
   cd fabric-testnet-nano-without-syschannel
   ./network.sh up
   ```

   > Consulte el README de `fabric-testnet-nano-without-syschannel` para detalles de configuración.

2. **Ejecutar la aplicación cliente**

   ```bash
   cd application-go
   go mod vendor
   go build -o app-mycc
   ./app-mycc
   ```

   * Observa la salida de la consola:

     * `--> Submit Transaction` indica que se está enviando una transacción al ledger.
     * `--> Evaluate Transaction` indica que se está consultando el ledger sin modificarlo.
     * `*** Result` muestra los resultados devueltos por el chaincode.
