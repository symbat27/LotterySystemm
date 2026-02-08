package handlers

import (
	"LotterySystem/internal/services"
	"encoding/json"
	"net/http"
)

type TicketHandler struct {
	service *services.LotteryService
}

func NewTicketHandler(s *services.LotteryService) *TicketHandler {
	return &TicketHandler{service: s}
}

func (h *TicketHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/tickets", h.handleTickets)
	mux.HandleFunc("/api/tickets/user", h.getUserTickets)
	mux.HandleFunc("/api/tickets/detail", h.getTicketDetail)
}

func (h *TicketHandler) handleTickets(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.createTicket(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Method not allowed",
	})
}

func (h *TicketHandler) createTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		UserID  string `json:"user_id"`
		DrawID  string `json:"draw_id"`
		Numbers []int  `json:"numbers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	ticket, err := h.service.CreateTicket(req.UserID, req.DrawID, req.Numbers)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"ticket":  ticket,
	})
}

func (h *TicketHandler) getUserTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Method not allowed",
		})
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User ID is required",
		})
		return
	}

	tickets := h.service.GetUserTickets(userID)

	result := make([]map[string]interface{}, 0, len(tickets))
	for _, ticket := range tickets {

		draw, err := h.service.GetDraw(ticket.DrawID)
		drawStatus := "unknown"
		if err == nil {
			drawStatus = draw.Status
		}

		ticketData := map[string]interface{}{
			"ticket":      ticket,
			"prize":       nil,
			"draw_status": drawStatus, // 🔥 ВАЖНО
		}

		if ticket.PrizeID != "" {
			prize, err := h.service.GetPrizeByTicket(ticket.ID)
			if err == nil {
				ticketData["prize"] = prize
			}
		}

		result = append(result, ticketData)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *TicketHandler) getTicketDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Method not allowed",
		})
		return
	}

	ticketID := r.URL.Query().Get("id")
	if ticketID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Ticket ID is required",
		})
		return
	}

	ticket, err := h.service.GetTicket(ticketID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	result := map[string]interface{}{
		"ticket": ticket,
		"prize":  nil,
	}

	if ticket.PrizeID != "" {
		prize, err := h.service.GetPrizeByTicket(ticket.ID)
		if err == nil {
			result["prize"] = prize
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
