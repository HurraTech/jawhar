#!/bin/bash
cd zahif
bash gen-grpc.sh
cp internal/server/proto/zahif.pb.go ../internal/zahif/proto/
