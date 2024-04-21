package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	Addition       = 0
	Subtraction    = 0
	Multiplication = 0
	Division       = 0
)

func CalculateExpr(expression string) (float64, error) {
	tokens := Tokenize(expression)
	postfix, err := InfixToPostfix(tokens)
	if err != nil {
		return 0, err
	}
	result, err := EvaluatePostfix(postfix)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func Tokenize(expression string) []string {
	// Удаляем все пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")

	var tokens []string
	var buffer strings.Builder

	for _, char := range expression {
		if IsOperator(char) {
			// Если встретили оператор, добавляем предыдущее число и оператор в список токенов,
			// а затем сбрасываем буфер
			if buffer.Len() > 0 {
				tokens = append(tokens, buffer.String())
				buffer.Reset()
			}
			tokens = append(tokens, string(char))
		} else if char == '(' || char == ')' {
			// Если встретили скобку, добавляем предыдущее число (если есть) и скобку в список токенов,
			// а затем сбрасываем буфер
			if buffer.Len() > 0 {
				tokens = append(tokens, buffer.String())
				buffer.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			buffer.WriteRune(char)
		}
	}

	// Добавляем последнее число в список токенов, если оно есть
	if buffer.Len() > 0 {
		tokens = append(tokens, buffer.String())
	}

	return tokens
}

func InfixToPostfix(tokens []string) ([]string, error) {
	var postfix []string
	var stack []string

	for _, token := range tokens {
		switch token {
		case "+", "-":
			// Приоритет операций "+" и "-" равен 1,
			// поэтому выталкиваем из стека все операции с большим или равным приоритетом
			for len(stack) > 0 && (stack[len(stack)-1] == "+" || stack[len(stack)-1] == "-" || stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "*", "/":
			// Приоритет операций "*" и "/" равен 2,
			// поэтому выталкиваем из стека все операции с большим или равным приоритетом
			for len(stack) > 0 && (stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			// Выталкиваем все операции из стека в постфиксную форму до тех пор, пока не встретим открывающую скобку "("
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			// Удаляем открывающую скобку "(" из стека
			if len(stack) > 0 && stack[len(stack)-1] == "(" {
				stack = stack[:len(stack)-1]
			} else {
				return nil, fmt.Errorf("Не найдена открывающая скобка")
			}
		default:
			postfix = append(postfix, token)
		}
	}

	// Выталкиваем все оставшиеся операции из стека в постфиксную форму
	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("Не найдена закрывающая скобка")
		}
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return postfix, nil
}

func IsOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

// Функция для вычисления выражеия
func EvaluatePostfix(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token {
		case "+":
			if len(stack) < 2 {
				return 0, fmt.Errorf("Недостаточно операндов для операции сложения")
			}

			//Время выполнения операции
			for i := 0; i < Addition; i++ {
				time.Sleep(time.Second)
			}

			result := stack[len(stack)-2] + stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, result)
		case "-":
			if len(stack) < 2 {
				return 0, fmt.Errorf("Недостаточно операндов для операции вычитания")
			}

			//Время выполнения операции
			for i := 0; i < Subtraction; i++ {
				time.Sleep(time.Second)
			}

			result := stack[len(stack)-2] - stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, result)
		case "*":
			if len(stack) < 2 {
				return 0, fmt.Errorf("Недостаточно операндов для операции умножения")
			}

			//Время выполнения операции
			for i := 0; i < Multiplication; i++ {
				time.Sleep(time.Second)
			}

			result := stack[len(stack)-2] * stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, result)
		case "/":
			if len(stack) < 2 {
				return 0, fmt.Errorf("Недостаточно операндов для операции деления")
			}
			if stack[len(stack)-1] == 0 {
				return 0, fmt.Errorf("Деление на ноль")
			}

			//Время выполнения операции
			for i := 0; i < Division; i++ {
				time.Sleep(time.Second)
			}

			result := stack[len(stack)-2] / stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, result)
		default:
			number, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("Некорректный токен: %s", token)
			}
			stack = append(stack, number)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("Некорректное выражение")
	}

	return stack[0], nil
}
