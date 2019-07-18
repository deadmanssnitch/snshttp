package snshttp

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var hostPattern = regexp.MustCompile(`^sns\.[a-zA-Z0-9\-]{3,}\.amazonaws\.com(\.cn)?$`)

// Verify will verify that a payload came from SNS
func Verify(signingSignature, signingCertURL, calSignature string) error {
	payloadSignature, err := base64.StdEncoding.DecodeString(signingSignature)
	if err != nil {
		return err
	}

	certURL, err := url.Parse(signingCertURL)
	if err != nil {
		return err
	}

	if !hostPattern.Match([]byte(certURL.Host)) {
		return fmt.Errorf("certificate is located on an invalid domain")
	}

	resp, err := http.Get(signingCertURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	decodedPem, _ := pem.Decode(body)
	if decodedPem == nil {
		return errors.New("the decoded PEM file was empty")
	}

	parsedCertificate, err := x509.ParseCertificate(decodedPem.Bytes)
	if err != nil {
		return err
	}

	return parsedCertificate.CheckSignature(x509.SHA1WithRSA, []byte(calSignature), payloadSignature)
}
