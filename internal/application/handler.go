package application

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/onel1nfavxx/YANDEX-GO_CALCULATOR_PROJECT/pkg/calculation"
)

type Request struct {
	Expression string `json:"expression"`
}

type ResponseSuccess struct {
	Result string `json:"result"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ResponseError{Error: "Method not allowed"})
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseError{Error: "Expression is not valid"})
		log.Println("Error decoding request body:", err)
		return
	}

	result, err := calculation.Calc(req.Expression)
	if err != nil {
		switch err.Error() {
		case "invalid expression", "division by zero":
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ResponseError{Error: "Expression is not valid"})
			log.Println("Calculation error:", err)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponseError{Error: "Internal server error"})
			log.Println("Internal server error:", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseSuccess{Result: floatToString(result)})
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
