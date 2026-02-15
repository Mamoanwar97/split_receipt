package service

import (
	"context"
	"fmt"
	"math"
"github.com/jackc/pgx/v5/pgtype"
	"github.com/mamoanwar97/split-receipt/internal/database"
)

type SettlementResult struct {
	FriendID       string             `json:"friend_id"`
	FriendName     string             `json:"friend_name"`
	MealTotal      float64            `json:"meal_total"`
	FixedFees      []FeeBreakdown     `json:"fixed_fees"`
	PercentageFees []FeeBreakdown     `json:"percentage_fees"`
	TotalDue       float64            `json:"total_due"`
}

type FeeBreakdown struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type SettlementService struct {
	queries *database.Queries
}

func NewSettlementService(q *database.Queries) *SettlementService {
	return &SettlementService{queries: q}
}

func (s *SettlementService) Calculate(ctx context.Context, receiptID, friendID pgtype.UUID) (*SettlementResult, error) {
	// Get the friend
	friend, err := s.queries.GetFriend(ctx, friendID)
	if err != nil {
		return nil, fmt.Errorf("friend not found: %w", err)
	}

	// Get meals this friend participates in for this receipt
	meals, err := s.queries.ListMealsByReceiptAndFriend(ctx, database.ListMealsByReceiptAndFriendParams{
		ReceiptID: receiptID,
		FriendID:  friendID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	// Calculate meal shares
	var mealTotal float64
	for _, meal := range meals {
		friendCount, err := s.queries.CountFriendsByMeal(ctx, meal.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to count friends for meal: %w", err)
		}
		if friendCount == 0 {
			continue
		}
		mealPrice := numericToFloat(meal.TotalPrice)
		mealTotal += mealPrice / float64(friendCount)
	}
	mealTotal = roundTo2(mealTotal)

	// Get unique friend count for fixed fee splitting
	uniqueFriends, err := s.queries.CountUniqueFriendsByReceipt(ctx, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to count unique friends: %w", err)
	}

	// Calculate fixed fee shares
	fixedFees, err := s.queries.ListFixedFeesByReceipt(ctx, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed fees: %w", err)
	}

	var fixedFeeBreakdowns []FeeBreakdown
	var fixedFeeTotal float64
	for _, fee := range fixedFees {
		if uniqueFriends == 0 {
			continue
		}
		share := roundTo2(numericToFloat(fee.Amount) / float64(uniqueFriends))
		fixedFeeBreakdowns = append(fixedFeeBreakdowns, FeeBreakdown{
			Name:   fee.Name,
			Amount: share,
		})
		fixedFeeTotal += share
	}

	// Calculate percentage fee shares
	percentageFees, err := s.queries.ListPercentageFeesByReceipt(ctx, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get percentage fees: %w", err)
	}

	var percentageFeeBreakdowns []FeeBreakdown
	var percentageFeeTotal float64
	for _, fee := range percentageFees {
		pct := numericToFloat(fee.Percentage) / 100.0
		amount := roundTo2(pct * mealTotal)

		// Apply cap if set
		if fee.CapAmount.Valid {
			cap := numericToFloat(fee.CapAmount)
			if amount > cap {
				amount = cap
			}
		}

		percentageFeeBreakdowns = append(percentageFeeBreakdowns, FeeBreakdown{
			Name:   fee.Name,
			Amount: amount,
		})
		percentageFeeTotal += amount
	}

	totalDue := roundTo2(mealTotal + fixedFeeTotal + percentageFeeTotal)

	return &SettlementResult{
		FriendID:       uuidToString(friendID),
		FriendName:     friend.Name,
		MealTotal:      mealTotal,
		FixedFees:      fixedFeeBreakdowns,
		PercentageFees: percentageFeeBreakdowns,
		TotalDue:       totalDue,
	}, nil
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	// Convert pgtype.Numeric to float64
	f, _ := n.Float64Value()
	return f.Float64
}

func roundTo2(f float64) float64 {
	return math.Round(f*100) / 100
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]),
		uint16(b[4])<<8|uint16(b[5]),
		uint16(b[6])<<8|uint16(b[7]),
		uint16(b[8])<<8|uint16(b[9]),
		uint64(b[10])<<40|uint64(b[11])<<32|uint64(b[12])<<24|uint64(b[13])<<16|uint64(b[14])<<8|uint64(b[15]),
	)
}

