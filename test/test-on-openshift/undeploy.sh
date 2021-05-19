BASE_DIR=$(cd $(dirname $0); pwd)

kubectl delete -f ${BASE_DIR}/ita2.yml
kubectl delete -f ${BASE_DIR}/ita1.yml
kubectl delete -f ${BASE_DIR}/pv-pvc.yml
