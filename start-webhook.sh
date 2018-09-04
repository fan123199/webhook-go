#!/usr/bin/env bash
pkill webhook-go
go build
nohup ./webhook-go 1> ~/log/webhook-go.out 2> ~/log/webhook-go.err &
