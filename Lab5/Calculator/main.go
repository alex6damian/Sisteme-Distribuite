package main

import (
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// helper function to perform calculations
func calc(a float64, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("împărțire la 0")
		}
		return a / b, nil
	case "^":
		return math.Pow(a, b), nil
	default:
		return 0, fmt.Errorf("operator necunoscut: %s", op)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Calculator")

	entry := widget.NewEntry()
	entry.Disable()
	entry.SetText("0")

	typing := false

	reset := false

	btn := func(label string, tapped func()) *widget.Button {
		return widget.NewButton(label, tapped)
	}

	addEntry := func(text string) {
		if !typing || entry.Text == "0" || reset {
			entry.SetText(text)
			typing = true
			reset = false
			return
		}
		entry.SetText(entry.Text + text)
	}

	operatorButton := func(op string) {
		switch op {
		case "+", "-", "*", "/", "^":
			entry.SetText(entry.Text + " " + op + " ")
		case "=":
			reset = true
			fmt.Println("Calculating result for:", entry.Text)
			var numbers []float64
			var operators []string

			var currentNum string
			for _, char := range entry.Text {
				if char == ' ' {
					if currentNum != "" {
						var num float64
						fmt.Sscanf(currentNum, "%f", &num)
						numbers = append(numbers, num)
						currentNum = ""
					}
				} else if char == '+' || char == '-' || char == '*' || char == '/' || char == '^' {
					operators = append(operators, string(char))
				} else {
					currentNum += string(char)
				}
			}
			if currentNum != "" {
				var num float64
				fmt.Sscanf(currentNum, "%f", &num)
				numbers = append(numbers, num)
			}

			if len(numbers) == 0 {
				return
			}

			// first pass: ^, *, /
			for i := 0; i < len(operators); {
				op := operators[i]
				if op == "^" || op == "*" || op == "/" {
					res, err := calc(numbers[i], numbers[i+1], op)
					if err != nil {
						entry.SetText("Error")
						return
					}
					numbers[i] = res
					// remove numbers[i+1] and operators[i]
					numbers = append(numbers[:i+1], numbers[i+2:]...)
					operators = append(operators[:i], operators[i+1:]...)
				} else {
					i++
				}
			}

			// second pass: + and -
			result := numbers[0]
			for i, op := range operators {
				res, err := calc(result, numbers[i+1], op)
				if err != nil {
					entry.SetText("Error")
					return
				}
				result = res
			}

			entry.SetText(fmt.Sprintf("%g", result))
		}
	}

	clearButton := func() {
		entry.SetText("0")
		typing = false
	}

	grid := container.NewGridWithColumns(4,
		btn("7", func() { addEntry("7") }),
		btn("8", func() { addEntry("8") }),
		btn("9", func() { addEntry("9") }),
		btn("/", func() { operatorButton("/") }),

		btn("4", func() { addEntry("4") }),
		btn("5", func() { addEntry("5") }),
		btn("6", func() { addEntry("6") }),
		btn("*", func() { operatorButton("*") }),

		btn("1", func() { addEntry("1") }),
		btn("2", func() { addEntry("2") }),
		btn("3", func() { addEntry("3") }),
		btn("-", func() { operatorButton("-") }),

		btn(".", func() {
			if !typing {
				addEntry("0.")
				return
			}
			if !contains(entry.Text, ".") {
				addEntry(".")
			}
		}),
		btn("0", func() { addEntry("0") }),
		btn("^", func() { operatorButton("^") }),
		btn("+", func() { operatorButton("+") }),

		btn("C", func() { clearButton() }),
		btn("=", func() { operatorButton("=") }),
	)

	content := container.NewVBox(entry, grid)
	w.SetContent(content)
	w.Resize(fyne.NewSize(360, 300))
	w.ShowAndRun()
}

func contains(s, sub string) bool {
	return indexOf(s, sub) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
