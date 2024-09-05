#!/usr/bin/env bash

source .env

gcloud auth login
gcloud auth application-default login

echo "Setting up gcloud to use project ${PROJECT_ID}"
gcloud config set project ${PROJECT_ID}
