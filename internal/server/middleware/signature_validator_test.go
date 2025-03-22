package middleware

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/frolmr/metrics.git/pkg/signer"
)

func TestWithSignature(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	signKey := "test-key"

	body := bytes.NewBufferString("test body")
	//nolint:noctx // No need for context in tests
	req, err := http.NewRequest("POST", "/test", body)
	if err != nil {
		t.Fatal(err)
	}

	reqSignature := signer.SignPayloadWithKey(body.Bytes(), []byte(signKey))
	req.Header.Set(domain.SignatureHeader, hex.EncodeToString(reqSignature))

	rr := httptest.NewRecorder()

	handler := WithSignature(signKey)(testHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := "test response"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}

	expectedSignature := "2711cc23e9ab1b8a9bc0fe991238da92671624a9ebdaf1c1abec06e7e9a14f9b"
	actualSignature := rr.Header().Get(domain.SignatureHeader)
	if actualSignature != expectedSignature {
		t.Errorf("handler returned unexpected signature: got %v want %v", actualSignature, expectedSignature)
	}
}

func TestWithSignature_InvalidSignature(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	signKey := "test-key"

	body := bytes.NewBufferString("test body")
	//nolint:noctx // No need for context in tests
	req, err := http.NewRequest("POST", "/test", body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(domain.SignatureHeader, "invalid-signature")

	rr := httptest.NewRecorder()

	handler := WithSignature(signKey)(testHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadGateway {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadGateway)
	}

	expectedBody := "invalid signature\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}

func TestWithSignature_NoSignature(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	signKey := "test-key"

	body := bytes.NewBufferString("test body")
	//nolint:noctx // No need for context in tests
	req, err := http.NewRequest("POST", "/test", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := WithSignature(signKey)(testHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := "test response"
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}

	if rr.Header().Get(domain.SignatureHeader) != "" {
		t.Errorf("handler returned unexpected signature header: got %v want %v", rr.Header().Get(domain.SignatureHeader), "")
	}
}
