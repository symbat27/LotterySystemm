package main

import (
	"LotterySystem/internal/handlers"
	"LotterySystem/internal/services"
	"LotterySystem/internal/storage"
	"log"
	"net/http"
)

func main() {
	userRepo := storage.NewUserRepository()
	drawRepo := storage.NewDrawRepository()
	ticketRepo := storage.NewTicketRepository()
	prizeRepo := storage.NewPrizeRepository()

	service := services.NewLotteryService(
		userRepo,
		drawRepo,
		ticketRepo,
		prizeRepo,
	)

	fs := http.FileServer(http.Dir("./internal/frontend"))
	mux.Handle("/", fs)

	log.Println("Lottery Imitation System")
	log.Println("User panel:  http://localhost:8080/index.html")
	log.Println("Admin panel: http://localhost:8080/admin.html")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
