#!/bin/bash

if [ "$1" = "choas" ]; then
  #制造混乱，看看其他服务的反应程度
  kubectl delete deploy/nginx -n dis-obj-storage
  echo "1)删除API数据服务......"
  kubectl delete deploy/apinode -n dis-obj-storage
  echo "2)删除数据节点服务......"
  kubectl delete deploy/datanode -n dis-obj-storage
  echo "3)删除元数据节点服务......"
  kubectl delete deploy/metanode -n dis-obj-storage
  echo "4)删除namespace......"
  kubectl delete namespace dis-obj-storage
  kubectl delete serviceaccount my-service-account
  kubectl delete clusterrole my-cluster-role
  kubectl delete clusterrolebinding my-cluster-role-binding
  echo "5)删除ingress......"
  kubectl delete service traefik-ingress-service -n kube-system
  kubectl delete sa traefik-ingress-controller -n kube-system
  kubectl delete clusterrole traefik-ingress-controller -n kube-system
  kubectl delete clusterrolebinding traefik-ingress-controller -n kube-system
  kubectl delete ingress traefik-web-ui -n kube-system
  kubectl delete ds traefik-ingress -n kube-system
else
  echo "1)删除nginx对外提供服务"
  kubectl delete deploy/nginx -n dis-obj-storage
  echo "2)删除API数据服务......"
  kubectl delete deploy/apinode -n dis-obj-storage
  echo "3)删除数据节点服务......"
  kubectl delete deploy/datanode -n dis-obj-storage
  echo "4)删除元数据节点服务......"
  kubectl delete deploy/metanode -n dis-obj-storage
  echo "5)删除namespace......"
  kubectl delete namespace dis-obj-storage
  kubectl delete serviceaccount my-service-account
  kubectl delete clusterrole my-cluster-role
  kubectl delete clusterrolebinding my-cluster-role-binding
  echo "6)删除ingress......"
  kubectl delete service traefik-ingress-service -n kube-system
  kubectl delete sa traefik-ingress-controller -n kube-system
  kubectl delete clusterrole traefik-ingress-controller -n kube-system
  kubectl delete clusterrolebinding traefik-ingress-controller -n kube-system
  kubectl delete ingress traefik-web-ui -n kube-system
  kubectl delete ds traefik-ingress -n kube-system
  echo "end.删除完毕."
fi
