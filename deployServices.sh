#!/bin/bash

echo "1)部署元数据节点服务......"
kubectl apply -f namespace-deployment.yaml

echo "2)部署service-account......"
kubectl apply -f service-account.yaml

echo "3)部署cluster-role......"
kubectl apply -f cluster-role.yaml

echo "4)部署cluster-binding......"
kubectl apply -f cluster-binding.yaml

echo "5)部署元数据节点服务......"
kubectl apply -f metanode/metanode-deployment.yaml

echo "6)部署数据节点服务......"
kubectl apply -f datanode/datanode-deployment.yaml

echo "7)部署API数据服务......"
kubectl apply -f apinode/apinode-deployment.yaml

echo "8)部署nginx对外提供服务"
kubectl apply -f traefik/svc.yaml
kubectl apply -f traefik/rbac.yaml
kubectl apply -f traefik/ds.yaml
kubectl apply -f traefik/ingress.yaml

#kubectl apply -f nginx-deployment.yaml

echo "9)部署完成."
