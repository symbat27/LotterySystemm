package models

type PrizeType string

const (
	Money  PrizeType = "money"
	Travel PrizeType = "travel"
	Gift   PrizeType = "gift"
)

type Prize struct {
	ID           string    json:"id"
	TicketID     string    json:"ticket_id"
	UserID       string    json:"user_id"
	Type         PrizeType json:"type"
	Name         string    json:"name"
	Value        int       json:"value"
	MatchesCount int       json:"matches_count"
}

var PrizeDefinitions = map[int][]Prize{
	1: {
		{Type: Money, Name: "Consolation Prize - 100 TG", Value: 100, MatchesCount: 1},
	},
	2: {
		{Type: Gift, Name: "Electric Kettle", Value: 0, MatchesCount: 2},
	},
	3: {
		{Type: Money, Name: "Small Cash Prize", Value: 500, MatchesCount: 3},
	},
	4: {
		{Type: Money, Name: "Medium Cash Prize", Value: 2000, MatchesCount: 4},
	},
	5: {
		{Type: Travel, Name: "Travel Voucher", Value: 0, MatchesCount: 5},
	},
	6: {
		{Type: Money, Name: "Jackpot", Value: 100000, MatchesCount: 6},
	},
}