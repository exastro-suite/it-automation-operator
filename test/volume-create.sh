#!/bin/bash -u

BASE_DIR=$(cd $(dirname $0); pwd)
TARGET=$1

${BASE_DIR}/pv-create.sh $TARGET
${BASE_DIR}/pvc-create.sh $TARGET
kubectl apply -f ${BASE_DIR}/job-pv-initializer.yml
kubectl wait --for=condition=complete --timeout=120s job/pv-initializer
kubectl delete -f ${BASE_DIR}/job-pv-initializer.yml
