package handlers

import (
	"LotterySystem/internal/services"
	"encoding/json"
	"net/http"
)

type AdminHandler struct {
	service *services.LotteryService
}

func NewAdminHandler(s *services.LotteryService) *AdminHandler {
	return &AdminHandler{service: s}
}

func (h *AdminHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/admin/draws", h.handleDraws)
	mux.HandleFunc("/api/admin/draws/execute", h.executeDraw)
	mux.HandleFunc("/api/admin/draws/pending", h.getPendingDraw)
	mux.HandleFunc("/api/admin/stats", h.getStats)
	mux.HandleFunc("/api/admin/prizes", h.getPrizes)
}

func (h *AdminHandler) handleDraws(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listDraws(w, r)
	case http.MethodPost:
		h.createDraw(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AdminHandler) createDraw(w http.ResponseWriter, _ *http.Request) {
	draw, err := h.service.CreateDraw()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"draw":    draw,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AdminHandler) listDraws(w http.ResponseWriter, _ *http.Request) {
	draws := h.service.ListDraws()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(draws); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AdminHandler) getPendingDraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	draw, err := h.service.GetPendingDraw()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draw)
}

func (h *AdminHandler) executeDraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		DrawID string `json:"draw_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	draw, err := h.service.ExecuteDraw(req.DrawID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"draw":    draw,
		"message": "Draw executed successfully",
	})
}

func (h *AdminHandler) getStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := h.service.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *AdminHandler) getPrizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prizes := h.service.GetAllPrizes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prizes)
}
