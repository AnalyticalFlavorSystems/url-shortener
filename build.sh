#!/bin/bash

go build -a -tags netgo -installsuffix netgo -ldflags '-w' . &&
docker build -t quay.io/gastrograph/url . &&
docker push quay.io/gastrograph/url
