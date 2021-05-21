#!/bin/bash -u

BASE_DIR=$(cd $(dirname $0); pwd)
TARGET=$1

${BASE_DIR}/undeploy.sh ${TARGET}
${BASE_DIR}/pvc-delete.sh ${TARGET}
${BASE_DIR}/pv-delete.sh ${TARGET}
