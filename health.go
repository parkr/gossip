package gossip

import "net/http"

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "no database configured", http.StatusInternalServerError)
		return
	}

	if h.DB.GetConnection() == nil {
		http.Error(w, "no database connection", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(`healthy`))
}
