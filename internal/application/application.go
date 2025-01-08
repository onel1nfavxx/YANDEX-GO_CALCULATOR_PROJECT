package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/onel1nfavxx/YANDEX-GO_CALCULATOR_PROJECT/custom_errors"
	"github.com/onel1nfavxx/YANDEX-GO_CALCULATOR_PROJECT/pkg/calculation"
)

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

type Application struct {
	config *Config
	logger *zap.Logger
}

func New() *Application {
	lg, _ := zap.NewDevelopment()
	return &Application{
		config: ConfigFromEnv(),
		logger: lg,
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()
	req := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Sugar().Fatalf("Bad Request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := calculation.Calc(req.Expression)
	if err != nil {
		if errors.Is(err, custom_errors.ErrInvInputs) {
			logger.Sugar().Errorf("Error: %s, Status: %d, Message: %s", req.Expression, http.StatusUnprocessableEntity, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Expression is not valid"}`))
		} else if errors.Is(err, custom_errors.ErrDivisionByZero) {
			logger.Sugar().Errorf("Error: %s, Status: %d, Message: %s", req.Expression, http.StatusUnprocessableEntity, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "Division by zero"}`))
		} else {
			logger.Sugar().Errorf("Error: %s, Status: %d, Message: %s", req.Expression, http.StatusInternalServerError, err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
	} else {
		logger.Sugar().Infof("Successful calculation: %s = %.2f", req.Expression, res)
		fmt.Fprintf(w, `{"result": "%.2f"}`, res)
	}
}
func (a *Application) RunServer() error {
	a.logger.Info("Starting server on port " + a.config.Addr)
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
