#! /usr/bin/env bash 

KEYWORD='mingregister'
# KEYWORD='mingregister-配置服务路由'
# KEYWORD='mingregister-InteractiveWithEtcd'
# KEYWORD='mingregister-AccessControl'
# KEYWORD='mingregister-iptables'
# KEYWORD='mingregister-options'
# KEYWORK='More info:'

find ./ -type f -iname *.go | xargs fgrep "${KEYWORD}"
