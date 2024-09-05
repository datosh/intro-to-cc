package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/go-sev-guest/abi"
	"github.com/google/go-sev-guest/client"
	"github.com/google/go-sev-guest/proto/sevsnp"
	"github.com/google/go-sev-guest/verify"
)

func main() {
	dev, err := client.OpenDevice()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := dev.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	nonce := [64]byte{0x01}
	report, certChain, err := generateReport(dev, nonce)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ChipIP: %x\n", report.ChipId)
	log.Printf("VMPL: %d\n", report.Vmpl)

	log.Printf(
		"Commited TCB: %d.%d.%d\n",
		report.CommittedMajor, report.CommittedMinor, report.CommittedBuild,
	)

	log.Printf("Measurement: %x\n", report.Measurement)

	verifyReport(report, certChain, nonce)
}

func generateReport(dev *client.LinuxDevice, nonce [64]byte) (*sevsnp.Report, []byte, error) {
	rawReport, certChain, err := client.GetRawExtendedReport(dev, nonce)
	if err != nil {
		log.Fatal(err)
	}
	report, err := abi.ReportToProto(rawReport)
	if err != nil {
		log.Fatal(err)
	}
	return report, certChain, nil
}

// verifyReport based on ARK -> ASK -> VCEK -> report
func verifyReport(report *sevsnp.Report, certChain []byte, nonce [64]byte) {
	vcek, err := getVCEKFromCertChain(certChain)
	if err != nil {
		log.Fatal(err)
	}

	ask, ark, err := downloadAskArk()
	if err != nil {
		log.Fatal(err)
	}

	err = verifyCertChain(ask, ark, vcek, report)
	if err != nil {
		log.Fatal(err)
	}

	err = checkNonce(nonce, report)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation Report Verified!")
}

func checkNonce(nonce [64]byte, report *sevsnp.Report) error {
	if len(report.ReportData) != 64 {
		return fmt.Errorf("report data should be 64 byte, but is %d byte", len(report.ReportData))
	}

	for i := 0; i < 64; i++ {
		if nonce[i] != report.ReportData[i] {
			return fmt.Errorf("nonce byte %d does not match %x != %x", i, nonce[i], report.ReportData[i])
		}
	}

	return nil
}

func verifyCertChain(ask, ark, vcek *x509.Certificate, report *sevsnp.Report) error {
	if err := ask.CheckSignatureFrom(ark); err != nil {
		return fmt.Errorf("verifying ARK -> ASK: %w", err)
	}
	if err := vcek.CheckSignatureFrom(ask); err != nil {
		return fmt.Errorf("verifying ASK -> VCEK: %w", err)
	}
	if err := verify.SnpProtoReportSignature(report, vcek); err != nil {
		return fmt.Errorf("verifying VCEK -> report: %w", err)
	}

	return nil
}

func getVCEKFromCertChain(certChain []byte) (*x509.Certificate, error) {
	certs := new(abi.CertTable)
	if err := certs.Unmarshal(certChain); err != nil {
		return nil, fmt.Errorf("unmarshaling cert chain: %w", err)
	}
	vcekCertRaw, err := certs.GetByGUIDString(abi.VcekGUID)
	if err != nil {
		return nil, fmt.Errorf("retrieving VCEK: %w", err)
	}
	vcek, err := x509.ParseCertificate(vcekCertRaw)
	if err != nil {
		return nil, fmt.Errorf("parsing VCEK: %w", err)
	}

	return vcek, nil
}

func downloadAskArk() (*x509.Certificate, *x509.Certificate, error) {
	resp, err := http.Get("https://kdsintf.amd.com/vcek/v1/Milan/cert_chain")
	if err != nil {
		return nil, nil, fmt.Errorf("downloading cert chain: %w", err)
	}
	defer resp.Body.Close()
	cert_chain, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading cert chain: %w", err)
	}

	askPEM, rest := pem.Decode(cert_chain)
	if askPEM == nil {
		return nil, nil, fmt.Errorf("decoding ASK: %w", err)
	}
	arkPEM, rest := pem.Decode(rest)
	if arkPEM == nil || len(rest) != 0 {
		return nil, nil, fmt.Errorf("decoding ARK: %w", err)
	}

	ask, err := x509.ParseCertificate(askPEM.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing ASK: %w", err)
	}
	ark, err := x509.ParseCertificate(arkPEM.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing ARK: %w", err)
	}

	return ask, ark, nil
}
