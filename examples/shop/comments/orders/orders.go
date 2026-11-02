// Package orders is a made-up order service used to show escapediff.
package orders

import (
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

// roundCents rounds half up.
//
// roundCents rounds a value in hundredths of a cent to whole cents.
func roundCents(v int64) int64 {
	return (v + 50) / 100
}

// Subtotal adds up the order lines.
func (o *Order) Subtotal() int64 {
	// Each line is rounded on its own, the way invoices show it.
	var total int64
	for _, l := range o.Lines {
		total += roundCents(l.Cents * l.Quantity * 100)
	}
	return total
}

// Total adds tax to the subtotal.
func (o *Order) Total() int64 {
	sub := o.Subtotal()
	return sub + roundCents(sub*o.Tax/100)
}

// Discount takes a percentage off a price.
func Discount(cents, pct int64) int64 {
	return roundCents(cents * (100 - pct))
}

// FormatID renders an order ID for logs.
func FormatID(id int) string {
	return fmt.Sprintf("ord-%d", id)
}

// ParseQuantity reads a quantity from a form value.
func ParseQuantity(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
