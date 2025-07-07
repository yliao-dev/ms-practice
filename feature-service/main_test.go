package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)


 var TestCases = []struct {
	name string
	feature string
	wantBody FeatureResponse
	errorMsg string
	wantStatus int
} {
	{
		"feature enabled", 
		"new-checkout", 
		FeatureResponse{Feature: "new-checkout", Enabled: true},
		"",
		http.StatusOK,
	},
	{
		"feature disabled", 
		"dark-mode", 
		FeatureResponse{Feature: "dark-mode", Enabled: false},
		"",
		http.StatusOK,
	},
	{
		"feature not found", 
		"legacy-report", 
		FeatureResponse{Feature: "new-checkout", Enabled: true},
		`feature not found: legacy-report`,
		http.StatusNotFound,
	},
}

func TestFeatureFlagHnadler(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/features/"+tc.feature, nil)
			response := httptest.NewRecorder()
			FeatureFlagHandler(response, request)

			if response.Code != tc.wantStatus {
				t.Errorf("wrong status: got %v, want %v", response.Code, tc.wantStatus)
			}
			if tc.wantStatus == http.StatusNotFound {
				// compare the error message
				var errRes ErrorResponse
				err := json.Unmarshal(response.Body.Bytes(), &errRes)
				if err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if errRes.Error != tc.errorMsg {
					t.Errorf("wrong error message: got %v, wanted %v", errRes, tc.errorMsg)
				}
				return
				
			}
			var got FeatureResponse
			err := json.Unmarshal(response.Body.Bytes(), &got)
			if err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if got != tc.wantBody {
				t.Errorf("wrong body: got %v, want %v", got, tc.wantBody)
			}
			
		},
		
		)
	}
}