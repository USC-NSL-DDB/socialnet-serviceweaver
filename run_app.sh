#!/bin/bash

go mod tidy
weaver generate .
go run .
