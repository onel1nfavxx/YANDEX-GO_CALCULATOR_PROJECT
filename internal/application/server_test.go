package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	handler := http.HandlerFunc(CalcHandler)

	tests := []struct {
		name         string
		method       string
		requestBody  interface{}
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:         "Valid Expression",
			method:       http.MethodPost,
			requestBody:  Request{Expression: "2 + 3 * 4"},
			expectedCode: http.StatusOK,
			expectedBody: ResponseSuccess{Result: "14", Code: 200},
		},
		{
			name:         "Invalid Method",
			method:       http.MethodGet,
			requestBody:  nil,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: ResponseError{Error: "Method not allowed", Code: 405},
		},
		{
			name:         "Invalid Expression - Division by Zero",
			method:       http.MethodPost,
			requestBody:  Request{Expression: "10 / 0"},
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: ResponseError{Error: "Expression is not valid", Code: 422},
		},
		{
			name:         "Invalid Expression - Syntax Error",
			method:       http.MethodPost,
			requestBody:  Request{Expression: "5 + * 2"},
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: ResponseError{Error: "Expression is not valid", Code: 422},
		},
		{
			name:         "Malformed JSON",
			method:       http.MethodPost,
			requestBody:  "invalid json",
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: ResponseError{Error: "Expression is not valid", Code: 422},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			switch body := tt.requestBody.(type) {
			case string:
				reqBody = []byte(body)
			case Request:
				reqBody, err = json.Marshal(body)
				if err != nil {
					t.Fatalf("Не удалось сериализовать тело запроса: %v", err)
				}
			case nil:
				reqBody = nil
			default:
				t.Fatalf("Неизвестный тип тела запроса: %T", body)
			}

			req, err := http.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Не удалось создать запрос: %v", err)
			}

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedCode, rr.Code)
			}

			switch expected := tt.expectedBody.(type) {
			case ResponseSuccess:
				var resp ResponseSuccess
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Errorf("Не удалось декодировать ответ: %v", err)
				}
				if resp != expected {
					t.Errorf("Ожидался ответ %+v, получен %+v", expected, resp)
				}
			case ResponseError:
				var resp ResponseError
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Errorf("Не удалось декодировать ответ: %v", err)
				}
				if resp != expected {
					t.Errorf("Ожидался ответ %+v, получен %+v", expected, resp)
				}
			default:
				t.Errorf("Неизвестный тип ожидаемого тела ответа: %T", expected)
			}
		})
	}
}

func TestServer(t *testing.T) {
	os.Setenv("PORT", "0")

	testServer := httptest.NewServer(http.HandlerFunc(CalcHandler))
	defer testServer.Close()

	tests := []struct {
		name         string
		method       string
		expression   string
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:         "Server Valid Expression",
			method:       http.MethodPost,
			expression:   "10 - 2 * 3",
			expectedCode: http.StatusOK,
			expectedBody: ResponseSuccess{Result: "4", Code: 200},
		},
		{
			name:         "Server Invalid Expression",
			method:       http.MethodPost,
			expression:   "10 / (5 - 5)",
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: ResponseError{Error: "Expression is not valid", Code: 422},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(Request{Expression: tt.expression})
			if err != nil {
				t.Fatalf("Не удалось сериализовать тело запроса: %v", err)
			}

			resp, err := http.Post(testServer.URL, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Не удалось отправить запрос: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedCode, resp.StatusCode)
			}

			switch expected := tt.expectedBody.(type) {
			case ResponseSuccess:
				var respBody ResponseSuccess
				if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Errorf("Не удалось декодировать ответ: %v", err)
				}
				if respBody != expected {
					t.Errorf("Ожидался ответ %+v, получен %+v", expected, respBody)
				}
			case ResponseError:
				var respBody ResponseError
				if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Errorf("Не удалось декодировать ответ: %v", err)
				}
				if respBody != expected {
					t.Errorf("Ожидался ответ %+v, получен %+v", expected, respBody)
				}
			default:
				t.Errorf("Неизвестный тип ожидаемого тела ответа: %T", expected)
			}
		})
	}
}
