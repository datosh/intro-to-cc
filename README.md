# Demystifying Confidential Computing

## A Practical Introduction for Cloud Native Engineers

This is the accompanying repository for the talk and contains all the code
and resources used in the presentation.

### Configuration

Create a file called `.env` in the root of the repository and add the following
variables:

```bash
PROJECT_ID=<your-project-id>
REGION=<your-region>
ZONE=<your-zone>
```

### Login to GCP

```bash
gcloud auth login
gcloud auth application-default login
```

```bash
./configure.sh
```

### Find a suitable image

https://cloud.google.com/confidential-computing/confidential-vm/docs/supported-configurations#supported-os

```bash
› gcloud compute images list \
    --filter="guestOsFeatures[].type:(SEV_SNP_CAPABLE)"
NAME                                       PROJECT              FAMILY                         DEPRECATED  STATUS
centos-stream-9-v20240815                  centos-cloud         centos-stream-9                            READY
cos-101-17162-528-16                       cos-cloud            cos-101-lts                                READY
cos-105-17412-448-16                       cos-cloud            cos-105-lts                                READY
cos-109-17800-309-24                       cos-cloud            cos-109-lts                                READY
cos-113-18244-151-27                       cos-cloud            cos-113-lts                                READY
rhel-8-v20240815                           rhel-cloud           rhel-8                                     READY
rhel-9-v20240815                           rhel-cloud           rhel-9                                     READY
rhel-9-2-sap-v20240815                     rhel-sap-cloud       rhel-9-2-sap-ha                            READY
rhel-9-4-sap-v20240815                     rhel-sap-cloud       rhel-9-4-sap-ha                            READY
rocky-linux-8-optimized-gcp-v20240815      rocky-linux-cloud    rocky-linux-8-optimized-gcp                READY
rocky-linux-8-v20240815                    rocky-linux-cloud    rocky-linux-8                              READY
rocky-linux-9-optimized-gcp-v20240815      rocky-linux-cloud    rocky-linux-9-optimized-gcp                READY
rocky-linux-9-v20240815                    rocky-linux-cloud    rocky-linux-9                              READY
sles-15-sp5-v20240821-x86-64               suse-cloud           sles-15-sp5                                READY
sles-15-sp6-v20240807-x86-64               suse-cloud           sles-15                                    READY
sles-15-sp4-sap-v20240808-x86-64           suse-sap-cloud       sles-15-sp4-sap                            READY
sles-15-sp5-sap-v20240821-x86-64           suse-sap-cloud       sles-15-sp5-sap                            READY
sles-15-sp6-sap-v20240807-x86-64           suse-sap-cloud       sles-15-sp6-sap                            READY
sles-sap-15-sp4-hardened-v20240808-x86-64  suse-sap-cloud       sles-sap-15-sp4-hardened                   READY
sles-sap-15-sp5-hardened-v20240821-x86-64  suse-sap-cloud       sles-sap-15-sp5-hardened                   READY
sles-sap-15-sp6-hardened-v20240807-x86-64  suse-sap-cloud       sles-sap-15-sp6-hardened                   READY
ubuntu-pro-2004-focal-v20240829            ubuntu-os-pro-cloud  ubuntu-pro-2004-lts                        READY
ubuntu-pro-2204-jammy-v20240829            ubuntu-os-pro-cloud  ubuntu-pro-2204-lts                        READY
ubuntu-pro-2404-noble-amd64-v20240830      ubuntu-os-pro-cloud  ubuntu-pro-2404-lts-amd64                  READY
ubuntu-2004-focal-v20240829                ubuntu-os-cloud      ubuntu-2004-lts                            READY
ubuntu-2204-jammy-v20240829                ubuntu-os-cloud      ubuntu-2204-lts                            READY
cos-beta-117-18613-0-10                    cos-cloud            cos-beta                                   READY
cos-dev-121-18632-0-0                      cos-cloud            cos-dev                                    READY
ubuntu-2404-noble-amd64-v20240830          ubuntu-os-cloud      ubuntu-2404-lts-amd64                      READY
cos-stable-113-18244-151-27                cos-cloud            cos-stable                                 READY
ubuntu-minimal-2004-focal-v20240829        ubuntu-os-cloud      ubuntu-minimal-2004-lts                    READY
ubuntu-minimal-2204-jammy-v20240829        ubuntu-os-cloud      ubuntu-minimal-2204-lts                    READY
ubuntu-minimal-2404-noble-amd64-v20240829  ubuntu-os-cloud      ubuntu-minimal-2404-lts-amd64              READY
fedora-coreos-40-20240808-3-0-gcp-aarch64  fedora-coreos-cloud  fedora-coreos-stable-arm64                 READY
fedora-coreos-40-20240808-3-0-gcp-x86-64   fedora-coreos-cloud  fedora-coreos-stable                       READY
fedora-coreos-40-20240825-1-0-gcp-aarch64  fedora-coreos-cloud  fedora-coreos-next-arm64                   READY
fedora-coreos-40-20240825-1-0-gcp-x86-64   fedora-coreos-cloud  fedora-coreos-next                         READY
fedora-coreos-40-20240825-2-0-gcp-aarch64  fedora-coreos-cloud  fedora-coreos-testing-arm64                READY
fedora-coreos-40-20240825-2-0-gcp-x86-64   fedora-coreos-cloud  fedora-coreos-testing                      READY
```

### No resources available

```bash
google_compute_instance.default: Still creating... [10s elapsed]
╷
│ Error: Error waiting for instance to create: The zone 'projects/kammel-dev/zones/europe-west3-c' does not have enough resources available to fulfill the request.  '(resource type:compute)'.
│
│
│ A n2d-standard-2 VM instance with 1 local SSD(s) is not supported in the europe-west3-c zone with your configuration.  Retry your request in a zone that supports the machine type and configuration you require. For more information, view the troubleshooting documentation.
│
│
│ Troubleshooting documentation
│ https://cloud.google.com/compute/docs/resource-error
│
│
│   with google_compute_instance.default,
│   on vm.tf line 6, in resource "google_compute_instance" "default":
│    6: resource "google_compute_instance" "default" {
│
╵
```

### Verify

```bash
gcloud compute instances describe cvm --zone=europe-west4-a --format="yaml(confidentialInstanceConfig)"
```

### Connect

```bash
ssh datosh@$(terraform output -raw public_ip)
```

### Build and Push

```bash
CGO_ENABLED=0 go build -o gcp-snp ./cmd/gcp-snp
scp gcp-snp datosh@$(terraform output -raw public_ip):gcp-snp
ssh datosh@$(terraform output -raw public_ip)
```
