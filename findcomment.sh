#! /usr/bin/env bash 

KEYWORD='mingregister'
# KEYWORD='mingregister-配置服务路由'
# KEYWORD='mingregister-InteractiveWithEtcd'
# KEYWORD='mingregister-AccessControl'

find ./ -type f -iname *.go | xargs fgrep "${KEYWORD}"
