package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"test/test_app/app/model/response"
	"test/test_app/app/service/logger"
	"test/test_app/app/service/util"
)

type ISupabaseService interface {
	AddNewUsers(ctx context.Context, url string, email, password string) (interface{}, error)
	GetUser(ctx context.Context, url string) (response.Userdata, error)
	AddNewFile(ctx context.Context, url string, values map[string]interface{}) (interface{}, error)
	GetFile(ctx context.Context, url string) (response.Filedata, error)
}

type SupabaseService struct{}

var errorValue = make(map[string]interface{})
var successValue = make(map[string]interface{})
var req *http.Request

func NewSupabaseService() ISupabaseService {
	return &SupabaseService{}
}

func (s *SupabaseService) AddNewUsers(ctx context.Context, url string, email, password string) (interface{}, error) {
	baseUrl := util.GetEnvWithKey("BASE_URL_SUPABASE")
	log := logger.Logger(ctx)
	jwt := util.GetEnvWithKey("SUPABASE_JWT")
	client := &http.Client{}
	var values = make(map[string]interface{})
	values["email"] = email
	values["password"] = password
	marshalJson, _ := json.Marshal(values)
	responseBody := bytes.NewBuffer(marshalJson)
	req, err := http.NewRequest(http.MethodPost, baseUrl+url, responseBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", jwt)
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < 300
	if !statusOK {
		log.Errorf("unknown, status code: %d", resp.StatusCode)
		return nil, err
	} else if resp.StatusCode != http.StatusNoContent {
		return nil, nil
	}
	return nil, nil
}

func (s *SupabaseService) GetUser(ctx context.Context, url string) (response.Userdata, error) {
	baseUrl := util.GetEnvWithKey("BASE_URL_SUPABASE")
	log := logger.Logger(ctx)
	jwt := util.GetEnvWithKey("SUPABASE_JWT")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, baseUrl+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", jwt)
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < 300
	if !statusOK {
		log.Errorf("unknown, status code: %d", resp.StatusCode)
		return nil, err
	} else if resp.StatusCode != http.StatusNoContent {
		Userdata := response.Userdata{}
		if err = json.NewDecoder(resp.Body).Decode(&Userdata); err == nil {
			return Userdata, nil
		}
	}
	return nil, err
}


func (s *SupabaseService) AddNewFile(ctx context.Context, url string, values map[string]interface{}) (interface{}, error) {
	baseUrl := util.GetEnvWithKey("BASE_URL_SUPABASE")
	log := logger.Logger(ctx)
	jwt := util.GetEnvWithKey("SUPABASE_JWT")
	client := &http.Client{}
	marshalJson, _ := json.Marshal(values)
	responseBody := bytes.NewBuffer(marshalJson)
	req, err := http.NewRequest(http.MethodPost, baseUrl+url, responseBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", jwt)
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < 300
	if !statusOK {
		log.Errorf("unknown, status code: %d", resp.StatusCode)
		return nil, err
	} else if resp.StatusCode != http.StatusNoContent {
		return nil, nil
	}
	return nil, nil
}



func (s *SupabaseService) GetFile(ctx context.Context, url string) (response.Filedata, error) {
	baseUrl := util.GetEnvWithKey("BASE_URL_SUPABASE")
	log := logger.Logger(ctx)
	jwt := util.GetEnvWithKey("SUPABASE_JWT")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, baseUrl+url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", jwt)
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < 300
	if !statusOK {
		log.Errorf("unknown, status code: %d", resp.StatusCode)
		return nil, err
	} else if resp.StatusCode != http.StatusNoContent {
		Filedata := response.Filedata{}
		if err = json.NewDecoder(resp.Body).Decode(&Filedata); err == nil {
			return Filedata, nil
		}
	}
	return nil, err
}