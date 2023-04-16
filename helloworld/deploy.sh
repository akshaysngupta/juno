#!/bin/bash

az group create -n myResourceGroup -l eastus

az vm create --resource-group myResourceGroup --name myVM --image UbuntuLTS --admin-username azureuser --generate-ssh-keys