package util

import (
	"math"
	"strings"
)

var (
	charset = "abcdefghijklmnopqrstuvwxyzA0123456789BCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 进制
	system int64 = 62
)

func Base62(incr int64) (string, error) {
	var shortUri []byte

	for incr > 0 {
		var r byte     // 下标指向的字符
		var tmp []byte // 64进制字符数组下标
		number := incr % system

		r = charset[number]

		tmp = append(tmp, r)
		shortUri = append(tmp, shortUri...)
		incr = incr / system
	}

	return string(shortUri), nil
}

func Decode62(s string) (int64, error) {
	var result int64

	s = ReverseString(s)
	for index, char := range s {
		i := int64(strings.Index(charset, string(rune(char))))
		result += i * int64(math.Pow(float64(system), float64(index)))

	}
	return result, nil
}

func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
