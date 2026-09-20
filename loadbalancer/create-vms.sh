#!/bin/bash

multipass launch 24.04 \
  --name lb \
  --cpus 1 \
  --memory 1G \
  --disk 5G

multipass launch 24.04 \
  --name app1 \
  --cpus 1 \
  --memory 1G \
  --disk 5G

multipass launch 24.04 \
  --name app2 \
  --cpus 1 \
  --memory 1G \
  --disk 5G

echo "VMs created"
multipass list