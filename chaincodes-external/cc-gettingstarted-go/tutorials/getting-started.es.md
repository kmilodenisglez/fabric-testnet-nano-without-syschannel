# Getting started - Empezando

## Contenido del tutorial
- [Requisitos previos](#Requisitos-previos)
- [Gestión interna](#Gestión-interna)
- [Declarando el contrato](#Declarando-el-contrato)
- [Escribiendo funciones del contrato](#escribiendo-funciones-del-contrato)
- [Usando contratos en un chaincode](#usando-contratos-en-un-chaincode)
- [Probando el chaincode como desarrollador](#Probando-el-chaincode-como-desarrollador)
- [Qué hacer a continuación?](#Que-hacer-a-continuación?)

## Requisitos previos
Este tutorial asume que tienes:
- [Go 1.17.x](https://golang.org/doc/install)
- Un clon de [fabric-testnet-nano-devmode](https://github.com/kmilodenisglez/fabric-testnet-nano-devmode)

## Gestión interna
Dado que este tutorial hará uso de la configuración ' `fabric-testnet-nano-devmode`, debe desarrollar dentro de `fabric-testnet-nano-devmode/chaincodes`. Cree una carpeta dentro de `fabric-testnet-nano-devmode/chaincodes` llamada `cc-gettingstarted-go` y abra ahí su editor preferido.

Ejecuta en su terminal el comando

```
go mod init github.com/kmilodenisglez/cc-gettingstarted-go
```

para configurar el módulo go, debes ejecutar:
 
```
go get -u github.com/hyperledger/fabric-contract-api-go
```
para obtener la última versión de fabric-contract-api-go para usar en su chaincode.

Recuerde ejecutar `go mod vendor`para construir un directorio denominado `vendor` que va a contener los paquetes necesario para compilar y el comando `go build` para construir el binario del chaincode.

```
go mod vendor
```

## Declarando el contrato
El contractapi genera un chaincode tomando uno o más "contratos" agrupados en un chaincode en ejecución. Lo primero que haremos aquí es declarar un contrato para usar en nuestro chaincode. Este contrato será simple, manejando la lectura y escritura de strings hacia y desde el world-state.

Todos los contratos para usar en un chaincode deben implementar [contractapi.ContractInterface](https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#ContractInterface). La forma más fácil de hacer esto es insertar la estructura `contractapi.Contract` dentro del contrato, que proporcionará la funcionalidad predeterminada.

Comenzamos el contrato creando un nuevo archivo `simple-contract.go` dentro de la carpeta `cc-gettingstarted-go`. Dentro del archivo, cree una estructura llamada `SimpleContract` que implementa la estructura `contractapi.Contract`. Este será nuestro contrato para administrar datos hacia y desde el world-state.

```
package main

import (
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// contrato SimpleContract para manejar escritura y lectura desde el world-state
type SimpleContract struct {
    contractapi.Contract
}
```

## Escribiendo funciones del contrato
Por defecto, todas las funciones públicas de una estructura se pueden llamar a través del chaincode; pero deben contar con un conjunto de reglas.

Si una función pública de un contrato en un chaincode no cumple con algunas de estas reglas, se devolverá un error en la creación del chaincode. Las reglas son las siguientes:

- La función de los contratos sólo podrá tomar los siguientes tipos de parámetros:
    - string
    - bool
    - int (incluido int8, int16, int32 y int64)
    - uint (incluido uint8, uint16, uint32 y uint64)
    - float32
    - float64
    - time.Time
    - Arrays/slices de cualquier tipo permitido
    - Structs (con los campos públicos de los tipos permitidos u otra estructura)
    - Punteros a structs
    - Maps con llave de tipo string y valores de cualquiera de los tipos permitidos
    - interface{} (Solo se permite cuando se retorna directamente, cuando se invoque a través de una transacción recibirá un string)
- Las funciones de los contratos solo pueden devolver cero, uno o dos valores
    - Si la función está definida para no retornar valores, se devolverá una respuesta de éxito para todas las llamadas a esa función de contrato.
    - Si la función está definida para retornar un unico valor, ese valor puede ser `error` o cualquiera de los tipos permitidos (excepto `interfaz{}`).
    - Si la función está definida para retornar dos valores, el primero puede ser cualquiera de los tipos permitidos enumerados para los parámetros (excepto `interfaz{}`) y el segundo debe ser `error`
- Las funciones de los contratos también pueden tomar el contexto de una transacción siempre que:
    - Sea el primer parámetro
    - Cualquiera de los siguientes casos
      - Es del tipo *contractapi.TransactionContext o un contexto de transacción personalizado definido en el chaincode que se usará para el contrato.
      - Es una interfaz que cumple el tipo de contexto de transacción en uso para el contrato, p. [contractapi.TransactionContextInterface](https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#TransactionContextInterface)

La primera función a escribir para `simple-contract.go` es `Crear`. Esta función agrega un nuevo par llave-valor al world-state utilizando una llave y un valor proporcionados por el usuario. Para interactúa con el world-state, necesitamos pasar el contexto de la transacción como argumento a la función.

Podemos tomar el contexto de transacción predeterminado proporcionado por contractapi (`contractapi.TransactionContext`) ya que proporciona todas las funciones necesarias para interactuar con el world-state. Sin embargo, tomar directamente `contractapi.TransactionContext` plantea algunos problemas,

¿qué pasaría si tuviéramos que escribir pruebas unitarias para nuestro contrato?

Tendríamos que crear una instancia de ese tipo que luego requeriría una instancia de [stub](https://godoc.org/github.com/hyperledger/fabric-chaincode-go/shim#ChaincodeStub) y terminaría haciendo nuestro pruebas complejas.

En cambio, lo que podemos hacer es tomar una interfaz que se encuentre con el contexto de la transacción; afortunadamente, el paquete contractapi define uno: `contractapi.TransactionContextInterface`. Esto significa que si tuviéramos que escribir pruebas unitarias, podríamos enviar un contexto de transacción simulado que luego podría usarse para rastrear llamadas o simplemente simplificar nuestra configuración de prueba. Como la función está destinada a escribir en lugar de devolver datos, solo retornará el tipo error.

```
// Crear agrega una nueva llave con valor al world-state
func (sc *SimpleContract) Crear(ctx contractapi.TransactionContextInterface, key string, value string) error {
    activoActual, err := ctx.GetStub().GetState(key)

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    if activoActual != nil {
        return fmt.Errorf("No se puede almacenar en el world state. La llave %s ya existe", key)
    }

    err = ctx.GetStub().PutState(key, []byte(value))

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    return nil
}
```

La función usa el stub del contexto de la transacción ([shim.ChaincodeStubInterface](https://godoc.org/github.com/hyperledger/fabric-chaincode-go/shim#ChaincodeStubInterface)) para leer primero el world-state, verificando que no existe ningún valor con la llave proporcionada, y luego coloca un nuevo valor en el world-state, convirtiendo el valor pasado en un arreglo de bytes según sea necesario.

La segunda función a agregar al contrato es `Actualizar`, funcionará similar a la función Crear, sin embargo, en lugar de generar un error si la llave existe en el world-state, se va a generar un error si no existe.

```
// Actualizar cambia el valor con llave en el world state
func (sc *SimpleContract) Actualizar(ctx contractapi.TransactionContextInterface, key string, value string) error {
    activoActual, err := ctx.GetStub().GetState(key)

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    if activoActual == nil {
        return fmt.Errorf("No se puede actualizar el world state. La llave %s no existe", key)
    }

    err = ctx.GetStub().PutState(key, []byte(value))

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    return nil
}
```

La tercera y última función para agregar al contrato es `Leer`. Esta recibe una llave y retornará el valor del world-state. Por lo tanto, retornará un string (el tipo de valor antes de convertir a bytes para el world-state) y también retornará un tipo error.

```
// Leer devuelve el valor en clave en el world-state
func (sc *SimpleContract) Leer(ctx contractapi.TransactionContextInterface, key string) (string, error) {	
    activoActual, err := ctx.GetStub().GetState(key)

    if err != nil {
        return "", errors.New("No se puede interactuar con el world state")
    }

    if activoActual == nil {
	return "", fmt.Errorf("No se puede leer en el world state. La llave %s no existe", key)
    }

    return string(activoActual), nil
}
```

El contrato final se verá así:

```
package main

import (
    "errors"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// contrato SimpleContract para manejar escritura y lectura desde el world-state
type SimpleContract struct {
    contractapi.Contract
}

// Crear agrega una nueva llave con valor al world-state
func (sc *SimpleContract) Crear(ctx contractapi.TransactionContextInterface, key string, value string) error {
    activoActual, err := ctx.GetStub().GetState(key)

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    if activoActual != nil {
        return fmt.Errorf("No se puede crear en el world state. La llave %s ya existe", key)
    }

    err = ctx.GetStub().PutState(key, []byte(value))

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    return nil
}

// Actualizar cambia el valor con llave en el world state
func (sc *SimpleContract) Actualizar(ctx contractapi.TransactionContextInterface, key string, value string) error {
    activoActual, err := ctx.GetStub().GetState(key)

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    if activoActual == nil {
        return fmt.Errorf("No se puede actualizar el world state. La llave %s no existe", key)
    }

    err = ctx.GetStub().PutState(key, []byte(value))

    if err != nil {
        return errors.New("No se puede interactuar con el world state")
    }

    return nil
}

// Leer devuelve el valor en clave en el world-state
func (sc *SimpleContract) Leer(ctx contractapi.TransactionContextInterface, key string) (string, error) {
    activoActual, err := ctx.GetStub().GetState(key)

    if err != nil {
        return "", errors.New("No se puede interactuar con el world state")
    }

    if activoActual == nil {
        return "", fmt.Errorf("No se puede leer en el world state. La llave %s no existe", key)
    }

    return string(activoActual), nil
}
```

## Usando contratos en un chaincode
En la misma carpeta que su archivo `simple-contract.go`, cree un archivo llamado `main.go`. Aquí agregue una función `main`. Esto se llamará cuando se ejecute su programa go.

```
package main

import (
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
}
```

Hasta ahora, ha creado un contrato; sin embargo, Fabric utiliza un chaincode que debe cumplir con la interfaz [shim.Chaincode](https://godoc.org/github.com/hyperledger/fabric-chaincode-go/shim#Chaincode). La interfaz de chaincode necesita dos funciones Init e Invoke. Afortunadamente, no necesita escribirlos, ya que contractapi proporciona una forma de generar un chaincode a partir de uno o más contratos. Para crear un chaincode, agregue lo siguiente a su función `main`:

```
    simpleContract := new(SimpleContract)

    cc, err := contractapi.NewChaincode(simpleContract)

    if err != nil {
        panic(err.Error())
    }
```

Una vez que tenga su chaincode, debe iniciarlo para que se pueda llamar a través de transacciones. Para hacer esto, agregue lo siguiente a continuación donde crea su chaincode:

```
    if err := cc.Start(); err != nil {
        panic(err.Error())
    }
```

Your `main.go` file should now look like this:

```
package main

import (
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
    simpleContract := new(SimpleContract)

    cc, err := contractapi.NewChaincode(simpleContract)

    if err != nil {
        panic(err.Error())
    }

    if err := cc.Start(); err != nil {
        panic(err.Error())
    }
}
```

## Probando el chaincode como desarrollador
Antes debe seguir los pasos del [readme.md](https://github.com/kmilodenisglez/fabric-testnet-nano-devmode#instructions-for-starting-network) para levantar la red de desarrollo
### Running the chaincode

```
cd cc-gettingstarted-go
```

sigue los pasos del [readme.md](https://github.com/kmilodenisglez/fabric-testnet-nano-devmode#1-build-the-chaincode) para construir chaincode

Ahora ejecuta el chaincode:

> Nota: debe mantenerse escuchando

Recuerde exportar las variables de entornos antes de operar con el chacinode usando el script `env.sh`

```
source ./env.sh
```

sigue los pasos del [readme.md](https://github.com/kmilodenisglez/fabric-testnet-nano-devmode#1-start-the-chaincode) para construir chaincode

Una vez que se crea una instancia del chaincode, se puede emitir transacciones para llamar a las funciones del contrato dentro del chaincode. Primero use el `invoke` para crear un nuevo par llave-valor en el world-state:

```
peer chaincode invoke -o 127.0.0.1:7050 -n mycc -c '{"Args":["Crear", "KEY_1", "VALUE_1"]}' -C mychannel
```

El primer argumento de la invocación es la función que desea llamar. Como solo tiene un contrato en su chaincode, simplemente puede pasar el nombre de la función. Los siguientes argumentos conforman los valores que se enviarán a la función. Los argumentos en fabric enviados a un chaincode siempre son string, sin embargo, como se describió anteriormente, una función de contrato puede tomar tipos que no sean string. El chaincode generado por contractapi maneja la conversión de estos valores (aunque en este caso nuestra función toma string); puede obtener más información sobre esto en tutoriales posteriores. Tenga en cuenta que no tiene que especificar el contexto de la transacción a pesar de que la función `Crear` toma uno, este se genera para ti.

Ahora que ha creado su par de llave-valor, puede usar la función Actualizar de su contrato para cambiar el valor. Esto nuevamente se puede hacer emitiendo un comando de invocación en el contenedor CLI:
```
peer chaincode invoke -o 127.0.0.1:7050 -c '{"Args":["Actualizar", "KEY_1", "VALUE_2"]}' -n mycc -C mychannel
```

Luego puede leer el valor almacenado para una llave emitiendo un comando de consulta contra la función de Leer del contrato:

```
peer chaincode query -o 127.0.0.1:7050 -c '{"Args":["Leer", "KEY_1"]}' -n mycc -C mychannel
```

Debería retornar el valor "VALUE_2".

## Plus

```
func (sc *SimpleContract) MostrarInfoStub(ctx contractapi.TransactionContextInterface) error {
    log.Println("**********************")
    log.Println("")
    log.Println("[Channel ID] ", ctx.GetStub().GetChannelID())
    log.Println("")
    log.Println("**********************")
    log.Println("")
    log.Println("[Transaction ID] ", ctx.GetStub().GetTxID())
    log.Println("")
    log.Println("**********************")

	// obteniendo el ID del cliente que invoka la TX
	identityID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return errors.New("Obteniendo el ID del cliente")
	}

	// obteniendo el MSP-ID del cliente que invoka la TX
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return errors.New("Obteniendo el MSP-ID del cliente")
	}

	log.Println("[Client identity] ", identityID)
	log.Println("")
	log.Println("**********************")
	log.Println("")
	log.Println("[Client MSP-ID] ", mspID)
	log.Println("")
	log.Println("**********************")
	log.Println("")
	return nil
}
```
La función MostrarInfoStub es un ejemplo de como trabajar con el Stub y ClientIdentity, para obtener el identificador del canal, ID de la transacción y ID del cliente que inicia la Tx. Asi como la Organizacion (MSPID).

## ¿Qué es lo próximo?
Siga el tutorial [Uso de funciones avanzadas](./using-advanced-features.md).
