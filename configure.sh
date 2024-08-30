#!/usr/bin/env bash

source .env

echo "Setting up gcloud to use project ${PROJECT_ID}"
gcloud config set project ${PROJECT_ID}
