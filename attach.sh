#!/bin/bash

POD=$(kubectl get pods -n gomud -l "app=app-mud" -o jsonpath="{.items[0].metadata.name}")

kubectl attach -n gomud $POD -i