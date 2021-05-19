BASE_DIR=$(cd $(dirname $0); pwd)

kubectl apply -f ${BASE_DIR}/pv-pvc.yml
kubectl apply -f ${BASE_DIR}/ita1.yml
kubectl apply -f ${BASE_DIR}/ita2.yml
