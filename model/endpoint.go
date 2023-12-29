package model

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type AddRequest struct {
	Name string `json:"name"`
}

type RemoveRequest struct {
	ID int `json:"id"`
}

type AddReponse struct {
	err error
}

type RemoveResponse struct {
	err error
}

type GetAllResponse struct {
	payload []model
	err     error
}

func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(AddRequest)
		err := s.Add(input.Name)
		return &AddReponse{
			err: err,
		}, nil
	}
}

func MakeRemoveEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(RemoveRequest)
		err := s.Remove(input.ID)
		return &RemoveResponse{
			err: err,
		}, nil
	}
}

func MakeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetAll()
		return &GetAllResponse{
			payload: res,
			err:     err,
		}, nil
	}
}
