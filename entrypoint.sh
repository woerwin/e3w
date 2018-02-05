#!/bin/bash
set -e

cat /med/hosts.test >> /etc/hosts

# app 启动脚本
/med/$1 -conf staging.conf -listen :80
