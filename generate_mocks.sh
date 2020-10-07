#!/bin/bash

rm -drf mocks
mkdir -p mocks/mock_io && mockgen io Writer > mocks/mock_io/api.go