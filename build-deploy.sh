#!/usr/bin/env bash

CGO_ENABLED=0 go build -o gcp-snp ./cmd/gcp-snp
scp gcp-snp datosh@$(terraform output -raw public_ip):gcp-snp
ssh datosh@$(terraform output -raw public_ip)
