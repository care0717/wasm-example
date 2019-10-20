package main

import(
	"errors"
	"regexp"
	"strconv"
	"syscall/js"
)


var stack string

func inputNum(this js.Value, i []js.Value) interface{} {
	// 引数を文字列に変換して連結する
	stack += i[0].String()
	setResult(stack)
	return nil
}

func doClear(this js.Value, i []js.Value) interface{} {
	clear()
	setResult(stack)
	return nil
}

func clear() {
	stack = ""
}

func setResult(res string) {
	js.Global().Get("document").Call("getElementById", "result").Set("textContent", res)
}


func isOpe(b byte) bool {
	switch b {
	case '/','+','-','*':
		return true
	default:
		return false
	}
}

func inputOpe(this js.Value, i []js.Value) interface{} {
	if len(stack) == 0 {
		return nil
	}
	if isOpe(stack[len(stack)-1]) {
		stack = stack[:len(stack)-1]
	}
	stack += i[0].String()
	setResult(stack)
	return nil
}

func calculate(stack string) (int, error) {
	ope := regexp.MustCompile(`/|\*|\+|-`)
	num := regexp.MustCompile(`\d+`)
	nums := ope.Split(stack, -1)
	numbers := make([]int, len(nums))
	for i, v := range(nums) {
		n, _ := strconv.Atoi(v)
		numbers[i] = n
	}
	opes := num.ReplaceAllString(stack, "")
	res, err := calc(numbers, opes)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func calc(numbers []int, opes string) (int, error) {
	if len(numbers) == 1 {
		return numbers[0], nil
	}
	tmp, err := exec(opes[0], numbers[0], numbers[1])
	if err != nil {
		return 0, err
	}
	numbers[1] = tmp
	numbers = numbers[1:]
	opes = opes[1:]
	return calc(numbers, opes)
}
func exec(ope byte, v1 int, v2 int) (int, error) {
	switch ope {
	case '+':
		return v1 + v2, nil
	case '-':
		return v1-v2, nil
	case '/':
		return v1/v2, nil
	case '*':
		return v1*v2, nil
	default:
		return 0, errors.New("undefined operator")
	}
}

func doCalc(this js.Value, i []js.Value) interface{} {
	if len(stack) == 0 || isOpe(stack[len(stack)-1]) {
		return nil
	}

	res, err := calculate(stack)
	if err != nil {
		clear()
		setResult(err.Error())
	}
	stack = strconv.Itoa(res)
	setResult(stack)
	return nil
}

func registerCallbacks() {
	js.Global().Set("inputNum", js.FuncOf(inputNum))
	js.Global().Set("inputOpe", js.FuncOf(inputOpe))
	js.Global().Set("doCalc", js.FuncOf(doCalc))
	js.Global().Set("doClear", js.FuncOf(doClear))
}

func main(){
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}