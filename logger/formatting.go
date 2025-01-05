package logger

import "fmt"

func PrintIndented(items []string, indent string, lineLen int) {
	divided := len(items) / lineLen

	var n int
	for i := range int(divided) {
		if i != 0 && divided > 0 {
			fmt.Print(indent)
		}

		num := n + lineLen
		for ; n < num; n++ {
			fmt.Print(items[n] + " ")
		}

		fmt.Println()
	}

	if divided > 0 {
		fmt.Print(indent)
	}

	for ; n < len(items); n++ {
		fmt.Print(items[n] + " ")
	}

	fmt.Println()
}
