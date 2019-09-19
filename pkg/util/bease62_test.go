package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestBase62(t *testing.T) {
	r, _ := Decode62("ZZZZ")
	fmt.Println("Decode62", r)

	e, _ := Base62(14776335)
	fmt.Println("Base62", e)

	fmt.Println(string(r) == e)
	a := "https://www.google.com"
	if strings.HasPrefix(a, "https://") {
		fmt.Println(true)
	}

	fmt.Println(BytesToInt64([]byte("1444")))
	//n := 926132836
	//for {
	//	e, _ := Base62(int64(n))
	//	if len(e) > 6 {
	//		fmt.Println(n)
	//		os.Exit(0)
	//	}
	//	n += 1
	//}
}

// 14776335
// 916132832
