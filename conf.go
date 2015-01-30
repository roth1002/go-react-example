package main

import "github.com/olebedev/config"

// Base configuration
var conf, _ = config.ParseYaml(`
env: development
app:
  name: go react example
api:
  prefix: /api/v1
duktape:
  poolSize: 5
debug: true
port: 5000
title: Go React Example
`)
