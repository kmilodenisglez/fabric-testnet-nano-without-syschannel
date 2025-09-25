#!/usr/bin/env sh
set -eu

### Base variables
COUCHDB_USER="admin"
COUCHDB_PASSWORD="password"
COUCHDB_PORT="5984"
COUCHDB_IMAGE="couchdb:3.4.2"
COUCHDB_CONTAINER="worldstate"

### Logging helpers
log() {
  echo "[INFO] $*"
}

error() {
  echo "[ERROR] $*" >&2
  exit 1
}

### Detect host system to set chaincode address
case "$(uname)" in
  Linux*)   CCADDR="127.0.0.1" ;;
  Darwin*)  CCADDR="host.docker.internal" ;;
  *)        CCADDR="127.0.0.1"; log "Unknown system, defaulting CCADDR to $CCADDR" ;;
esac

### Check required dependencies
for cmd in peer docker curl; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    error "Required command '$cmd' is not installed or not in PATH"
  fi
done

### If peer binary is not available in PATH, look in ../bin
if ! command -v peer >/dev/null 2>&1; then
  log "'peer' not found, adding ../bin to PATH"
  export PATH="${PWD}/../bin:$PATH"
fi

### Fabric environment variables
export FABRIC_CFG_PATH="${PWD}/../config"
export FABRIC_LOGGING_SPEC="debug:cauthdsl,policies,msp,grpc,peer.gossip.mcs,gossip,leveldbhelper=info"

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE="${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt"
export CORE_PEER_TLS_KEY_FILE="${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key"
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"

export CORE_PEER_ID="peer0.org1.example.com"
export CORE_PEER_ADDRESS="127.0.0.1:7051"
export CORE_PEER_LISTENADDRESS="127.0.0.1:7051"
export CORE_PEER_CHAINCODEADDRESS="${CCADDR}:7052"
export CORE_PEER_CHAINCODELISTENADDRESS="127.0.0.1:7052"
export CORE_PEER_GOSSIP_BOOTSTRAP="127.0.0.1:7053"
export CORE_PEER_GOSSIP_EXTERNALENDPOINT="127.0.0.1:7051"
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH="${PWD}/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp"
export CORE_OPERATIONS_LISTENADDRESS="127.0.0.1:8446"
export CORE_PEER_FILESYSTEMPATH="${PWD}/data/peer0.org1.example.com"
export CORE_LEDGER_SNAPSHOTS_ROOTDIR="${PWD}/data/peer0.org1.example.com/snapshots"

### Function to start CouchDB container and configure peer to use it
start_couchdb() {
  # Export environment variables for CouchDB state database
  export CORE_LEDGER_STATE_STATEDATABASE="CouchDB"
  export CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS="127.0.0.1:${COUCHDB_PORT}"
  export CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME="${COUCHDB_USER}"
  export CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD="${COUCHDB_PASSWORD}"

  log "Removing existing CouchDB container (if any)..."
  if docker ps -aqf "name=${COUCHDB_CONTAINER}" | grep -q .; then
    docker rm -f "${COUCHDB_CONTAINER}" >/dev/null
    log "Container ${COUCHDB_CONTAINER} removed."
  else
    log "No previous CouchDB container found."
  fi

  log "Starting CouchDB on port ${COUCHDB_PORT}..."
  docker run -d \
    -p "${COUCHDB_PORT}:${COUCHDB_PORT}" \
    -e COUCHDB_USER="${COUCHDB_USER}" \
    -e COUCHDB_PASSWORD="${COUCHDB_PASSWORD}" \
    --name "${COUCHDB_CONTAINER}" \
    "${COUCHDB_IMAGE}" >/dev/null

  log "Waiting for CouchDB to become available..."
  for i in $(seq 1 10); do
    if curl -s "http://${COUCHDB_USER}:${COUCHDB_PASSWORD}@127.0.0.1:${COUCHDB_PORT}/" >/dev/null; then
      log "CouchDB is ready ✅"
      break
    fi
    log "CouchDB not responding yet... retrying ($i/10)"
    sleep 2
  done
}

### --- MAIN EXECUTION ---

# Uncomment the line below if you want to use CouchDB
#start_couchdb

log "Starting peer node..."
exec peer node start
