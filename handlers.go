package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type handler struct {
	sharedSecret string
	payloadTtl   time.Duration
	logger       *CustomLogger
}

func newHandler(sharedSecret string, payloadTtl time.Duration, l *CustomLogger) *handler {
	return &handler{
		sharedSecret: sharedSecret,
		payloadTtl:   payloadTtl,
		logger:       l,
	}
}

func (h *handler) TestHandler1(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("TestHandler1 running")
	response := map[string]string{"msg": "TestHandler1 200"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *handler) TestGroup1Handler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("TestGroup1Handler running")
	response := map[string]string{"msg": "TestGroup1Handler 200"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// Обработчик для маршрута /test
func (h *handler) testHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"message": "ok"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		HttpRespErrRFC9457("testHandler", "Failed to encode JSON response", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}
}

func (h *handler) panicHandler(w http.ResponseWriter, r *http.Request) {

	// Вызываем панику для теста
	panic("something went wrong")

	response := map[string]string{"message": "ok"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		HttpRespErrRFC9457("testHandler", "Failed to encode JSON response", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}
}

func (h *handler) payloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload, err := generatePayload(h.sharedSecret, h.payloadTtl)
	if err != nil {
		HttpRespErrRFC9457("payloadHandler", "Failed to generatePayload", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}

	response := map[string]string{"payload": payload}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		HttpRespErrRFC9457("testHandler", "Failed to encode JSON response", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}
}

type Domain struct {
	LengthBytes uint32 `json:"lengthBytes"`
	Value       string `json:"value"`
}

type MessageInfo struct {
	Timestamp int64  `json:"timestamp"`
	Domain    Domain `json:"domain"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
	StateInit string `json:"state_init"`
}

type TonProof struct {
	Address    string      `json:"address"`
	Network    string      `json:"network"`
	Proof      MessageInfo `json:"proof"`
	RefferalID uint        `json:"refferal_id"`
}

func (h *handler) proofHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		HttpRespErrRFC9457("proofHandler", "Failed to Signature", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}
	var tp TonProof
	err = json.Unmarshal(b, &tp)
	if err != nil {
		HttpRespErrRFC9457("proofHandler", "Failed to Signature", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}

	err = checkPayload(tp.Proof.Payload, h.sharedSecret)
	if err != nil {
		HttpRespErrRFC9457("proofHandler", "Failed to Signature", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}

	response := map[string]string{"proof": "proof"}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		HttpRespErrRFC9457("proofHandler", "Failed to Signature", err, http.StatusInternalServerError, w, r, h.logger)
		return
	}
}
