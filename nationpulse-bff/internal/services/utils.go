package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nationpulse-bff/internal/repos"
	. "github.com/nationpulse-bff/internal/utils"
	"github.com/segmentio/kafka-go"
)

type UtilsService struct {
	Configs *Configs
	repo    *repos.UtilsRepo
}

func NewUtilsService(configs *Configs, repo *repos.UtilsRepo) *UtilsService {
	return &UtilsService{
		Configs: configs,
		repo:    repo,
	}
}

func (us *UtilsService) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching permissions...")
	userID := r.Form.Get("userID")
	// fmt.Println("USERID", userID)
	data, err := us.repo.GetPermissions(userID)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data, true, nil)
}

func (us *UtilsService) PublishExportRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting reporting service to generate csv")
	var request ExportApiMessageRequest

	req, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Printf("Error reading from request body: %s", err)
		WriteJSON(w, http.StatusBadRequest, nil, false, err)

	}

	if err := json.Unmarshal(req, &request); err != nil {
		fmt.Println("Error unmarshalling request body", err)
		WriteJSON(w, http.StatusBadRequest, nil, false, err)
	}
	// fmt.Println("REQUEST PAYLOAD", request)

	GetQueryAndHeaders(&request)
	// create new writer
	kw := us.Configs.Kafka.NewWriter("message-log")

	// create a timeout using context
	kwriteCtx, cancel := context.WithTimeout(us.Configs.Context, 30*time.Second)
	defer cancel()
	// marshal the request body to export request payload
	exportPayload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling export payload for kafka message", err)
		WriteJSON(w, http.StatusBadRequest, nil, false, err)
	}

	// publish the message to kafka broker
	publishExportMessage := kafka.Message{
		Key:   []byte(request.ExportID),
		Value: []byte(exportPayload),
	}

	defer kw.Close()
	err = kw.WriteMessages(kwriteCtx, publishExportMessage)

	if err != nil {
		fmt.Println("Error publishing message to kafka broker: ", err)
		WriteJSON(w, http.StatusInternalServerError, nil, false, err)
	}

	var response = struct {
		ExportID   string `json:"exportID"`
		Status     string `json:"status"`
		StatusCode int    `json:"statusCode"`
	}{
		ExportID:   request.ExportID,
		Status:     "processing",
		StatusCode: 0,
	}

	// fmt.Println("Response", response)
	WriteJSON(w, http.StatusOK, response, true, nil)
}

func (us *UtilsService) SubscribeExportResponse(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting reporting service to receive csv")
	// Setting headers
	origin := r.Header.Get("Origin")

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Expose-Headers", "Content-type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	var response ExportApiMessageResponse
	var event = "export-status"
	// create a newreader and read mesages. Upon reading messages create a payload and send the response to the SSE.
	fmt.Println("TOPIC======>", us.Configs.Cfg.KafkaReaderTopic)
	kr := us.Configs.Kafka.NewReader("my-group", us.Configs.Cfg.KafkaReaderTopic)
	defer kr.Close()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(w, ": ping\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}()

	processMessage := func(message kafka.Message) error {
		fmt.Println("PROCESSING MESSAGE & -----", message.Value)

		sb := strings.Builder{}

		if err := json.Unmarshal(message.Value, &response); err != nil {
			log.Println("Error unmarshalling message response: ", err)
			return err
		}

		buff := bytes.NewBuffer([]byte{})

		encoder := json.NewEncoder(buff)

		if err := encoder.Encode(response); err != nil {
			log.Println("Error encoding event data: ", err)
			return err
		}

		sb.WriteString(fmt.Sprintf("event: %s\n", event))
		sb.WriteString(fmt.Sprintf("data: %s\n", buff.String()))
		sb.WriteString("\n")

		fmt.Println("SSE event message", sb.String())

		if _, err := fmt.Fprint(w, sb.String()); err != nil {
			log.Println("Error sending event to w", err)
			return err
		}

		flusher.Flush()

		return nil
	}
	ReadMessages(kr, r.Context(), processMessage)
}
