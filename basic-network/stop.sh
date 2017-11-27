#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev
#docker commit "`docker ps | grep "peer node start" | awk ' { print $1} '`" hyperledger/fabric-peer:latest

# Shut down the Docker containers that might be currently running.
docker-compose -f docker-compose.yml stop
