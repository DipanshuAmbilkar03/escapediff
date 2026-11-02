package orders

import "strconv"

// Encoder writes an order total in some format.
type Encoder interface {
	Encode(buf []byte, cents int64) []byte
}

// PlainEncoder writes cents as a decimal number.
type PlainEncoder struct{}

// Encode appends the number to buf.
func (PlainEncoder) Encode(buf []byte, cents int64) []byte {
	return strconv.AppendInt(buf, cents, 10)
}

// DefaultEncoder is used by Receipt. Tests may replace it.
var DefaultEncoder Encoder = PlainEncoder{}

// Receipt renders the order total.
func Receipt(o *Order) []byte {
	buf := make([]byte, 0, 32)
	return DefaultEncoder.Encode(buf, o.Total())
}

// Batch totals many orders.
func Batch(orders []*Order) int64 {
	var sum int64
	for _, o := range orders {
		sum += o.Total()
	}
	return sum
}
