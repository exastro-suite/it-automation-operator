
#!/bin/bash -u

BASE_DIR=$(cd $(dirname $0); pwd)
TARGET=$1

oc delete -f ${BASE_DIR}/${TARGET}/pv.yml
