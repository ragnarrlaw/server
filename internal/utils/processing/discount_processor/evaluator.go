package discountprocessor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/ragnarrlaw/rules/rule_engine/lexer"
	"github.com/ragnarrlaw/rules/rule_engine/parser"
	"github.com/ragnarrlaw/server/internal/types/recommender"
)

func EvaluateRecommendations(recommendations []recommender.Recommendation) error {
	for i := range recommendations {
		for _, discount := range recommendations[i].Discounts {
			for j := range recommendations[i].Items {
				l := lexer.NewLexer(discount)
				p := parser.NewParser(l)
				r, err := p.ParseRule()
				if err != nil {
					return err
				}
				result, err := EvaluateLogicalCondition(r.Condition, &recommendations[i], &recommendations[i].Items[j])
				if err != nil {
					return err
				}
				if result {
					ApplyAction(r.Action, &recommendations[i], &recommendations[i].Items[j])
				}
				recommendations[i].DiscountedCost += recommendations[i].Items[j].DiscountedPrice
			}
		}
	}
	return nil
}

func EvaluateLogicalCondition(r *parser.LogicalCondition, recommendation *recommender.Recommendation, item *recommender.RecommendationItem) (bool, error) {
	leftResult := evaluateSingleCondition(r.Left, recommendation, item)
	if r.Operator != "" {
		rightResult := evaluateSingleCondition(r.Right, recommendation, item)
		switch r.Operator {
		case "AND":
			return leftResult && rightResult, nil
		case "OR":
			return leftResult || rightResult, nil
		default:
			return false, fmt.Errorf("unknown logical operator %s", r.Operator)
		}
	}
	return leftResult, nil
}

func evaluateSingleCondition(c *parser.Condition, recommendation *recommender.Recommendation, item *recommender.RecommendationItem) bool {
	switch strings.ToLower(c.Key) {
	case "product_id":
		return compareUUID(c.Operator, c.Value, item.ProductId)
	case "requested_quantity":
		return compare(c.Operator, c.Value, item.Quantity)
	case "cart_price":
		return compare(c.Operator, c.Value, recommendation.TotalCost)
	default:
		return false
	}
}

func compareUUID(operator string, left interface{}, right uuid.UUID) bool {
	switch operator {
	case "=":
		leftUUID, err := uuid.Parse(left.(string))
		if err != nil {
			return false
		}
		return leftUUID == right
	case "!=":
		leftUUID, err := uuid.Parse(left.(string))
		if err != nil {
			return false
		}
		return leftUUID != right
	case "IN", "in":
		for _, e := range left.([]string) {
			if uuid.MustParse(e) == right {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func compare(operator string, left interface{}, right interface{}) bool {
	switch operator {
	case "=":
		return left == right
	case "!=":
		return left != right
	case ">":
		return left.(float64) > right.(float64)
	case ">=":
		return left.(float64) >= right.(float64)
	case "<":
		return left.(float64) < right.(float64)
	case "<=":
		return left.(float64) <= right.(float64)
	case "IN", "in":
		for _, e := range left.([]string) {
			if e == right.(string) {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func ApplyAction(a *parser.Action, recommendation *recommender.Recommendation, item *recommender.RecommendationItem) error {
	c := strings.ToLower(a.DiscountType)
	f, err := strconv.ParseFloat(a.Value.(string), 64)
	if err != nil {
		return err
	}
	switch c {
	case "product_percentage":
		if item.DiscountedPrice == 0 {
			item.DiscountedPrice = item.Cost * (1 - f/100)
		} else {
			item.DiscountedPrice *= (1 - f/100)
		}
		return nil
	case "product_flat_amount":
		if item.DiscountedPrice == 0 {
			item.DiscountedPrice = item.Cost - f
		} else {
			item.DiscountedPrice -= f
		}
		return nil
	case "cart_percentage":
		if recommendation.DiscountedCost == 0 {
			recommendation.DiscountedCost = recommendation.TotalCost * (1 - f/100)
		} else {
			recommendation.DiscountedCost *= (1 - f/100)
		}
		return nil
	case "cart_flat_amount":
		if recommendation.DiscountedCost == 0 {
			recommendation.DiscountedCost = recommendation.TotalCost - f
		} else {
			recommendation.DiscountedCost -= f
		}
		return nil
	default:
		return fmt.Errorf("unknown discount type %s", c)
	}
}
