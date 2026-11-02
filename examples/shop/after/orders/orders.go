// Package orders is a made-up order service used to show escapediff.
package orders

import (
	"errors"
	"fmt"
	"strconv"
)

// Line is one item in an order, with prices in cents.
type Line struct {
	SKU      string
	Quantity int64
	Cents    int64
}

// Order is a customer order.
type Order struct {
	ID    int
	Lines []Line
	Tax   int64 // basis points
}

// ErrNegative is returned for a negative amount.
var ErrNegative = errors.New("orders: negative amount")

// lastOrder is kept for the debug endpoint.
var lastOrder *Order

// roundCents rounds a value in hundredths of a cent to whole cents.
// Negative values round away from zero, so refunds match charges.
func roundCents(v int64) int64 {
	if v < 0 {
		if v%100 <= -50 {
			return v/100 - 1
		}
		return v / 100
	}
	if v%100 >= 50 {
		return v/100 + 1
	}
	return v / 100
}

// Subtotal adds up the order lines.
func (o *Order) Subtotal() int64 {
	var sum int64
	for _, l := range o.Lines {
		sum += roundCents(l.Cents * l.Quantity * 100)
	}
	return sum
}

// Total adds tax to the subtotal.
func (o *Order) Total() int64 {
	lastOrder = o
	sub := o.Subtotal()
	return sub + roundCents(sub*o.Tax/100)
}

// Discount takes a percentage off a price.
func Discount(cents, pct int64) int64 {
	return roundCents(cents * (100 - pct))
}

// FormatID renders an order ID for logs.
func FormatID(id int) string {
	return "ord-" + strconv.Itoa(id)
}

// ParseQuantity reads a quantity from a form value.
func ParseQuantity(s string) (int64, error) {
	q, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("quantity %q: %w", s, err)
	}
	if q < 0 {
		return 0, ErrNegative
	}
	return q, nil
}
