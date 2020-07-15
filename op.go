package main

import "fmt"

type binOpType = int

const (
	binOpPlus binOpType = iota
	binOpMinus
	binOpTimes
	binOpDivide
)

func opName(op binOpType) (string, error) {
	switch op {
	case binOpPlus:
		return "+", nil
	case binOpMinus:
		return "-", nil
	case binOpTimes:
		return "*", nil
	case binOpDivide:
		return "/", nil
	default:
		return "??", fmt.Errorf("Unsupported operator: %d", op)
	}
}

func opAction(op binOpType) (string, error) {
	switch op {
	case binOpPlus:
		return "plus", nil
	case binOpMinus:
		return "minus", nil
	case binOpTimes:
		return "times", nil
	case binOpDivide:
		return "divide", nil
	default:
		return "??", fmt.Errorf("Unsupported operator: %d", op)
	}
}
