
#!/bin/bash -u

BASE_DIR=$(cd $(dirname $0); pwd)
TARGET=$1

kubectl apply -f ${BASE_DIR}/${TARGET}/pvc-ita1.yml 
kubectl apply -f ${BASE_DIR}/${TARGET}/pvc-ita2.yml 
