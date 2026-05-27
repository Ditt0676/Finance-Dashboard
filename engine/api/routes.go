package api

import (
	"context"
	"encoding/json"
	"finance-dashboard/engine/alerts"
	"finance-dashboard/engine/metrics"
	"finance-dashboard/engine/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	dataPath string
	watcher  *alerts.Watcher
	httpSrv  *http.Server
}

func NewServer(dataPath string) *Server {
	return &Server{
		dataPath: dataPath,
		watcher:  alerts.NewWatcher(2.0),
	}
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/stocks", s.handleStocks)
	mux.HandleFunc("/api/crypto", s.handleCrypto)
	mux.HandleFunc("/api/metrics", s.handleMetrics)
	mux.HandleFunc("/api/alerts", s.handleAlerts)

	s.httpSrv = &http.Server{
		Addr:         port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Rutas disponibles:")
	log.Println("   GET /health")
	log.Println("   GET /api/stocks")
	log.Println("   GET /api/crypto")
	log.Println("   GET /api/metrics")
	log.Println("   GET /api/alerts")

	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) handleStocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stocks, err := loadStocks(s.dataPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, stocks)
}

func (s *Server) handleCrypto(w http.ResponseWriter, r *http.Request) {
	crypto, err := loadCrypto(s.dataPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, crypto)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	result := map[string]models.Metrics{}

	stocks, err := loadStocks(s.dataPath)
	if err == nil {
		for symbol, stock := range stocks {
			result[symbol] = metrics.CalculateStockMetrics(stock)
		}
	}

	cryptos, err := loadCrypto(s.dataPath)
	if err == nil {
		for id, crypto := range cryptos {
			result[id] = metrics.CalculateCryptoMetrics(crypto)
		}
	}

	if len(result) == 0 {
		http.Error(w, "no data available", http.StatusInternalServerError)
		return
	}

	writeJSON(w, result)
}

func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	stocks, err := loadStocks(s.dataPath)
	if err == nil {
		for _, stock := range stocks {
			s.watcher.CheckStock(stock)
		}
	}

	cryptos, err := loadCrypto(s.dataPath)
	if err == nil {
		for _, crypto := range cryptos {
			s.watcher.CheckCrypto(crypto)
		}
	}

	writeJSON(w, s.watcher.DrainAlerts())
}

func loadStocks(dataPath string) (map[string]models.Stock, error) {
	data, err := os.ReadFile(filepath.Join(dataPath, "stocks.json"))
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer stocks.json: %v", err)
	}
	var stocks map[string]models.Stock
	if err := json.Unmarshal(data, &stocks); err != nil {
		return nil, fmt.Errorf("error parseando stocks.json: %v", err)
	}
	return stocks, nil
}

func loadCrypto(dataPath string) (map[string]models.Crypto, error) {
	data, err := os.ReadFile(filepath.Join(dataPath, "crypto.json"))
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer crypto.json: %v", err)
	}
	var crypto map[string]models.Crypto
	if err := json.Unmarshal(data, &crypto); err != nil {
		return nil, fmt.Errorf("error parseando crypto.json: %v", err)
	}
	return crypto, nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
