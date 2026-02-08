package services

import (
	"LotterySystem/internal/models"
	"LotterySystem/internal/storage"
	"LotterySystem/internal/utils"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type LotteryService struct {
	users   *storage.UserRepository
	draws   *storage.DrawRepository
	tickets *storage.TicketRepository
	prizes  *storage.PrizeRepository
	rng     *rand.Rand
}

func NewLotteryService(
	users *storage.UserRepository,
	draws *storage.DrawRepository,
	tickets *storage.TicketRepository,
	prizes *storage.PrizeRepository,
) *LotteryService {
	return &LotteryService{
		users:   users,
		draws:   draws,
		tickets: tickets,
		prizes:  prizes,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *LotteryService) RegisterUser(username, password string) (models.User, error) {
	if _, err := s.users.GetByUsername(username); err == nil {
		return models.User{}, errors.New("username already exists")
	}

	user := models.User{
		ID:        s.generateID(),
		Username:  username,
		Password:  password,
		Balance:   10000,
		CreatedAt: time.Now(),
	}

	if err := s.users.Save(user); err != nil {
		return models.User{}, err
	}

	user.Password = ""
	return user, nil
}

func (s *LotteryService) LoginUser(username, password string) (models.User, error) {
	user, err := s.users.GetByUsername(username)
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	if user.Password != password {
		return models.User{}, errors.New("invalid username or password")
	}

	user.Password = "" // Don't return password
	return user, nil
}

func (s *LotteryService) GetUser(userID string) (models.User, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return models.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (s *LotteryService) CreateDraw() (models.Draw, error) {
	if _, err := s.draws.GetPending(); err == nil {
		return models.Draw{}, errors.New("there is already an active draw")
	}

	draw := models.Draw{
		ID:             s.generateID(),
		WinningNumbers: []int{},
		Status:         "pending",
		DrawDate:       time.Now(),
		CreatedAt:      time.Now(),
	}

	if err := s.draws.Save(draw); err != nil {
		return models.Draw{}, err
	}

	return draw, nil
}

func (s *LotteryService) ExecuteDraw(drawID string) (models.Draw, error) {
	draw, err := s.draws.GetByID(drawID)
	if err != nil {
		return models.Draw{}, err
	}

	if draw.Status == "completed" {
		return models.Draw{}, errors.New("draw already completed")
	}

	draw.WinningNumbers = utils.GenerateWinningNumbers()
	draw.Status = "completed"

	if err := s.draws.Update(draw); err != nil {
		return models.Draw{}, err
	}

	if err := s.processDrawResults(draw); err != nil {
		return models.Draw{}, err
	}

	return draw, nil
}

func (s *LotteryService) GetDraw(drawID string) (models.Draw, error) {
	return s.draws.GetByID(drawID)
}

func (s *LotteryService) ListDraws() []models.Draw {
	return s.draws.List()
}

func (s *LotteryService) GetPendingDraw() (models.Draw, error) {
	return s.draws.GetPending()
}

func (s *LotteryService) CreateTicket(userID, drawID string, numbers []int) (models.Ticket, error) {
	// Validate numbers
	if !utils.ValidateNumbers(numbers) {
		return models.Ticket{}, errors.New("invalid numbers: must be 6 unique numbers between 1 and 49")
	}

	draw, err := s.draws.GetByID(drawID)
	if err != nil {
		return models.Ticket{}, errors.New("draw not found")
	}

	if draw.Status != "pending" {
		return models.Ticket{}, errors.New("draw is not accepting tickets")
	}

	user, err := s.users.GetByID(userID)
	if err != nil {
		return models.Ticket{}, errors.New("user not found")
	}

	ticketCost := 100
	if user.Balance < ticketCost {
		return models.Ticket{}, errors.New("insufficient balance")
	}

	user.Balance -= ticketCost
	if err := s.users.Update(user); err != nil {
		return models.Ticket{}, err
	}

	ticket := models.Ticket{
		ID:        s.generateID(),
		UserID:    userID,
		DrawID:    drawID,
		Numbers:   numbers,
		Matches:   0,
		CreatedAt: time.Now(),
	}

	if err := s.tickets.Save(ticket); err != nil {
		return models.Ticket{}, err
	}

	return ticket, nil
}

func (s *LotteryService) GetUserTickets(userID string) []models.Ticket {
	return s.tickets.GetByUserID(userID)
}

func (s *LotteryService) GetTicket(ticketID string) (models.Ticket, error) {
	return s.tickets.GetByID(ticketID)
}

func (s *LotteryService) processDrawResults(draw models.Draw) error {
	tickets := s.tickets.GetByDrawID(draw.ID)

	for _, ticket := range tickets {
		matches := utils.CountMatches(ticket.Numbers, draw.WinningNumbers)
		ticket.Matches = matches

		if matches < 1 {
			s.tickets.Update(ticket)
			continue
		}

		prize, err := s.awardPrize(ticket, matches)
		if err != nil {
			s.tickets.Update(ticket)
			continue
		}

		ticket.PrizeID = prize.ID

		if prize.Type == models.Money {
			user, err := s.users.GetByID(ticket.UserID)
			if err == nil {
				user.Balance += prize.Value
				s.users.Update(user)
			}
		}

		s.tickets.Update(ticket)
	}

	return nil
}

func (s *LotteryService) awardPrize(ticket models.Ticket, matches int) (models.Prize, error) {
	prizeDefs, exists := models.PrizeDefinitions[matches]
	if !exists || len(prizeDefs) == 0 {
		return models.Prize{}, errors.New("no prize for this match count")
	}

	selectedPrize := prizeDefs[s.rng.Intn(len(prizeDefs))]

	prize := models.Prize{
		ID:           s.generateID(),
		TicketID:     ticket.ID,
		UserID:       ticket.UserID, // 🔥 ВАЖНО
		Type:         selectedPrize.Type,
		Name:         selectedPrize.Name,
		Value:        selectedPrize.Value,
		MatchesCount: matches,
	}

	if err := s.prizes.Save(prize); err != nil {
		return models.Prize{}, err
	}

	return prize, nil
}

func (s *LotteryService) GetPrizeByTicket(ticketID string) (models.Prize, error) {
	return s.prizes.GetByTicketID(ticketID)
}

func (s *LotteryService) GetAllPrizes() []models.Prize {
	return s.prizes.List()
}

func (s *LotteryService) GetStats() map[string]interface{} {
	prizes := s.prizes.List()

	stats := map[string]int{
		"money":  0,
		"travel": 0,
		"gift":   0,
	}

	totalValue := 0
	for _, prize := range prizes {
		stats[string(prize.Type)]++
		totalValue += prize.Value
	}

	return map[string]interface{}{
		"total_prizes":   len(prizes),
		"prizes_by_type": stats,
		"total_value":    totalValue,
		"total_users":    len(s.users.List()),
		"total_tickets":  len(s.tickets.List()),
		"total_draws":    len(s.draws.List()),
	}
}

func (s *LotteryService) generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
