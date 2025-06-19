package server

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	model "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/prediction_v1/go"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/domain"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/utils"
)

type PredictionHandler struct {
	r *mux.Router
	logger *slog.Logger
	client model.PredictionServiceClient
}

func (h *PredictionHandler) setupgRPC(conn *grpc.ClientConn) {
	h.client = model.NewPredictionServiceClient(conn)
}

func (h *PredictionHandler) setupRoutes() {
	h.r.HandleFunc("/", h.PredictionHandler).Methods("POST")
}

func (h *PredictionHandler) PredictionHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("operation", "get car price prediction"),
	)

	params := new(domain.PredictionRequest)
	if err := utils.ParseJson(r.Body, params); err != nil {
		http.Error(w, "failed to decode prediction request body: " + err.Error(), http.StatusBadRequest)
	}

	log.Info(
		"Send given params to the prediction service...",
		slog.Any("params", params),
	)

	response, err := h.client.Predict(r.Context(), &model.PredictRequest{
		Make: params.Make,
		Model: params.Model,
		Year: int32(params.Year),
		Hp: int32(params.Hp),
		Body: params.Body,
		Yearsell: int32(params.YearSell),
		Odometer: int32(params.Odometer),
		Color: params.Color,
	})

	if err != nil {
		utils.HandleResponseErr(w, h.logger, "prediction operation failed - ", err)
		return
	}

	var buffer bytes.Buffer

	mimeWriter := multipart.NewWriter(&buffer)
	defer mimeWriter.Close()

	_ = mimeWriter.WriteField("price", strconv.Itoa(int(response.GetPrice())))
	_ = mimeWriter.WriteField("sell_count", strconv.Itoa(int(response.GetSellCount())))
	for i, url := range response.PhotoUrls {
		_ = mimeWriter.WriteField(fmt.Sprintf("url_%d", i), url)
	}

	imageWriter, err := mimeWriter.CreateFormFile("image", "prediction.png")
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(imageWriter, bytes.NewReader(response.GraphPng))
	if err != nil {
		http.Error(w, "Failed to write image data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mimeWriter.FormDataContentType())

	w.Write(buffer.Bytes())
}
