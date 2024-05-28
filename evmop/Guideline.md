# Guideline

## Requirement
```bash
# create the project
mkdir evmop
cd evmop
operator-sdk init --domain gocrazy.com --repo evmop

# create an api
operator-sdk create api --group learn --version v1alpha1 --kind Blockchain --resource --controller

# login dockerhub (for docker push)
docker login -u <username>

# update controller
make update

# create application
kubectl create namespace ethereum
kubectl apply -f config/samples/learn_v1alpha1_blockchain.yaml

# forward port
kubectl port-forward blockchain-sample-0 8545:8545 -n ethereum
# test json-rpc api
curl http://localhost:8545/ \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}'
```
