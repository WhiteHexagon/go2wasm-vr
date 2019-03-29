#!/bin/bash
GOOS=js GOARCH=wasm go build -ldflags='-s -w' -o website/app.wasm

