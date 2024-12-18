package processing

import (
	"fmt"

	"github.com/raganrrlaw/server/internal/types/product"
)

func UnitConverter(from product.UnitOfMeasure, to product.UnitOfMeasure) (float64, error) {
	switch from {
	case product.UnitOfMeasureGrams:
		switch to {
		case product.UnitOfMeasureGrams:
			return 1.0, nil
		case product.UnitOfMeasureKg:
			return 1000.0, nil
		default:
			return -1, fmt.Errorf("from: %s cannot be converted to: %s", from, to)
		}
	case product.UnitOfMeasureKg:
		switch to {
		case product.UnitOfMeasureGrams:
			return 1000.0, nil
		case product.UnitOfMeasureKg:
			return 1.0, nil
		default:
			return -1, fmt.Errorf("from: %s cannot be converted to: %s", from, to)
		}
	case product.UnitOfMeasureMl:
		switch to {
		case product.UnitOfMeasureMl:
			return 1.0, nil
		case product.UnitOfMeasureL:
			return 1000.0, nil
		default:
			return -1, fmt.Errorf("from: %s cannot be converted to: %s", from, to)
		}
	case product.UnitOfMeasureL:
		switch to {
		case product.UnitOfMeasureMl:
			return 1000.0, nil
		case product.UnitOfMeasureL:
			return 1.0, nil
		default:
			return -1, fmt.Errorf("from: %s cannot be converted to: %s", from, to)
		}
	case product.UnitOfMeasurePcs:
		switch to {
		case product.UnitOfMeasurePcs:
			return 1.0, nil
		default:
			return -1, fmt.Errorf("from: %s cannot be converted to: %s", from, to)
		}
	default:
		return -1, fmt.Errorf("from: invalid unit of measure")
	}
}
