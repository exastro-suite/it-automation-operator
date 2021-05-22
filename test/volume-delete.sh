#!/bin/bash -u

BASE_DIR=$(cd $(dirname $0); pwd)
TARGET=$1

kubectl delete -f ${BASE_DIR}/job-pv-initializer.yml
${BASE_DIR}/pvc-delete.sh $TARGET
${BASE_DIR}/pv-delete.sh $TARGET
