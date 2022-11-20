# Off-chain data store sample

The off-chain data store sample demonstrates:

- Receiving block events in a client application.
- Using a checkpointer to resume event listening after a failure or application restart.
- Extracting ledger updates from block events in order to build an off-chain data store.

## About the sample

This sample shows how to replicate the data in your blockchain network to an off-chain data store. Using an off-chain data store allows you to analyze the data from your network or build a dashboard without degrading the performance of your application.

### Application

The client application provides several "commands" that can be invoked using the command-line:

- **listen**: Listen for block events, and use them to replicate ledger updates in an off-chain data store. See:
  - TypeScript: [application-typescript/src/listen.ts](application-typescript/src/listen.ts)

To keep the sample code concise, the **listen** command writes ledger updates to an output file named `store.log` in the current working directory (which for the Java sample is the `application-java/app` directory). A real implementation could write ledger updates directly to an off-chain data store of choice. You can inspect the information captured in this file as you run the sample.

Note that the **listen** command is is restartable and will resume event listening after the last successfully processed block / transaction. This is achieved using a checkpointer to persist the current listening position. Checkpoint state is persisted to a file named `checkpoint.json` in the current working directory. If no checkpoint state is present, event listening begins from the start of the ledger (block number zero).

## Running the sample

The Fabric test network is used to deploy and run this sample. Follow these steps in order:

1. Create and run the test network.

2. Deploy one smart contract.

3. Run listener

   ```bash
   # To run the application
   cd application-typescript
2. Deploy one smart contract in network.

4. Build app
```bash
cd application-typescript
npm run build
```

5. Run listener
To run the application
   ```bash 
   npm start listen
   ```

> Note: Interrupt the listener process using **Control-C**.

## Clean up

The persisted event checkpoint position can be removed by deleting the `checkpoint.json` file while the listener is stopped.

The recorded ledger updates can be removed by deleting the `store.log` file.

When you are finished, you can bring down the test network. Be sure to remove the `checkpoint.json` and `store.log` files before attempting to run the application with a new network.