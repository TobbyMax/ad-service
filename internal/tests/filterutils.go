package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DateLayout     = "2006-01-02"
	DateTimeLayout = "Mon, 2 Jan 2006 15:04:05 MST"
)

func (tc *testClient) listAdsByStatus(published any) (adsResponse, error) {
	body := map[string]any{
		"published": published,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAdsByUser(userID int64) (adsResponse, error) {
	body := map[string]any{
		"user_id": userID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAdsByDate(date string) (adsResponse, error) {
	body := map[string]any{
		"date": date,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAdsByUserAndDate(userID int64, date string) (adsResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"date":    date,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAdsByOptions(userID int64, date string, published bool, title string) (adsResponse, error) {
	body := map[string]any{
		"user_id":   userID,
		"date":      date,
		"published": published,
		"title":     title,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAdsByTitle(title string) (adsResponse, error) {
	body := map[string]any{
		"title": title,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}
