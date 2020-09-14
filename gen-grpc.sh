#!/bin/bash
cd zahif
bash gen-grpc.sh
cp internal/proto/zahif.pb.go ../internal/zahif/proto/