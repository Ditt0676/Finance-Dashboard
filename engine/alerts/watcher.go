package alerts

import (
	"finance-dashboard/engine/models"
	"fmt"
	"log"
	"sync"
	"time"
)

type Watcher struct {
	threshold float64
	previous  map[string]float64
	Alerts    chan models.Alert
	mu        sync.Mutex
}

func NewWatcher(thresholdPct float64) *Watcher {
	return &Watcher{
		Alerts:    make(chan models.Alert, 100),
		threshold: thresholdPct,
		previous:  make(map[string]float64),
	}
}

func (w *Watcher) CheckStock(stock models.Stock) {
	w.check(stock.Symbol, stock.CurrentPrice)
}

func (w *Watcher) CheckCrypto(crypto models.Crypto) {
	w.check(crypto.ID, crypto.CurrentPrice)
}

func (w *Watcher) check(id string, current float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	prev, exists := w.previous[id]
	if !exists {
		w.previous[id] = current
		return
	}

	if prev == 0 {
		w.previous[id] = current
		return
	}

	changePct := ((current - prev) / prev) * 100

	var alertType, direction string
	switch {
	case changePct > w.threshold:
		alertType = "price_up"
		direction = "subió"
	case changePct < -w.threshold:
		alertType = "price_down"
		direction = "bajó"
	default:
		w.previous[id] = current
		return
	}

	select {
	case w.Alerts <- models.Alert{
		Symbol:    id,
		Type:      alertType,
		Message:   fmt.Sprintf("%s %s %.2f%% (de $%.2f a $%.2f)", id, direction, changePct, prev, current),
		Threshold: w.threshold,
		Triggered: current,
		Time:      time.Now().Format("2006-01-02 15:04:05"),
	}:
	default:
		log.Printf("Alerts channel full, dropping alert for %s", id)
	}

	w.previous[id] = current
}

func (w *Watcher) DrainAlerts() []models.Alert {
	result := []models.Alert{}
	for {
		select {
		case alert := <-w.Alerts:
			result = append(result, alert)
		default:
			return result
		}
	}
}
