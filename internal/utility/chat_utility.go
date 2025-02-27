package utility

import (
	"fmt"
	"strconv"
	"strings"
)

func FormatNumber(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits = len(in) - 1
	}
	if numOfDigits <= 3 {
		return in
	}

	var sb strings.Builder
	var commas int
	if n < 0 {
		sb.WriteByte('-')
	}

	// Tentukan berapa koma yang kita butuhkan
	commas = (numOfDigits - 1) / 3

	// Tentukan berapa digit pada grup pertama
	firstGroupLen := numOfDigits - commas*3

	// Tambahkan digit grup pertama
	if n < 0 {
		sb.WriteString(in[1 : firstGroupLen+1])
	} else {
		sb.WriteString(in[:firstGroupLen])
	}

	// Tambahkan digit grup-grup lainnya
	var nextDigit int
	if n < 0 {
		nextDigit = firstGroupLen + 1
	} else {
		nextDigit = firstGroupLen
	}
	for i := 0; i < commas; i++ {
		sb.WriteByte('.')
		sb.WriteString(in[nextDigit : nextDigit+3])
		nextDigit += 3
	}

	return sb.String()
}

func GetPriceRange(products []map[string]interface{}) (min, max float64) {
	if len(products) == 0 {
		return 0, 0
	}

	min = products[0]["price"].(float64)
	max = min

	for _, p := range products {
		price := p["price"].(float64)
		if price < min {
			min = price
		}
		if price > max {
			max = price
		}
	}

	return min, max
}

func FormatPriceRange(min, max float64) string {
	return fmt.Sprintf("Rp%s - Rp%s",
		FormatNumber(int64(min)),
		FormatNumber(int64(max)))
}
