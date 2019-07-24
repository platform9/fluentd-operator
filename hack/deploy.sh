#!/bin/bash
basepath=$(dirname $0)
kubectl create ns logging
kubectl create ns pf9-operators
kubectl apply -f ${basepath}/../deploy/crds/logging_v1alpha1_output_crd.yaml
kubectl apply -n logging -f ${basepath}/../deploy/fluent
kubectl apply -n pf9-operators -f ${basepath}/../deploy
