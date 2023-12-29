package model

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	httptransport "github.com/go-kit/kit/transport/http"
)

var ErrBadRequest = errors.New("bad request")
var ErrInvalidId = errors.New("invalid id")

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request AddRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, ErrBadRequest
	}
	return request, nil
}

func decodeRemoveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "ID"))

	if err != nil {
		return nil, ErrInvalidId
	}
	return RemoveRequest{
		ID: id,
	}, nil
}

func decodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func encodeAddResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(*AddReponse)
	return json.NewEncoder(w).Encode(res.err)
}

func encodeRemoveResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(*RemoveResponse)
	return json.NewEncoder(w).Encode(res.err)
}

func encodeGetAllResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(*GetAllResponse)
	if res.err != nil {
		return json.NewEncoder(w).Encode(res.err)
	}
	return json.NewEncoder(w).Encode(res.payload)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case ErrInvalidId:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func MakeHandler(s Service) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	addHandler := httptransport.NewServer(
		MakeAddEndpoint(s),
		decodeAddRequest,
		encodeAddResponse,
		options...,
	)

	removeHandler := httptransport.NewServer(
		MakeRemoveEndpoint(s),
		decodeRemoveRequest,
		encodeRemoveResponse,
		options...,
	)

	getAllHandler := httptransport.NewServer(
		MakeGetAllEndpoint(s),
		decodeGetAllRequest,
		encodeGetAllResponse,
		options...,
	)

	r := chi.NewRouter()
	r.Route("/items", func(r chi.Router) {
		r.Get("/", getAllHandler.ServeHTTP)
		r.Post("/add", addHandler.ServeHTTP)
		r.Get("/remove/{ID}", removeHandler.ServeHTTP)
	})

	return r
}
