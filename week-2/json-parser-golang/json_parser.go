package main

// var reservedCharacters  map[rune]rune
var reservedCharacters = map[rune]rune{
	',': ',',
	':': ':',
	'"': '"',
	'[': '[',
	']': ']',
	'{': '{',
	'}': '}',
}

var ignoredCharacters = map[rune]rune{
	' ':  ' ',
	'\n': '\n',
}

type StateMachineNode struct {
	char        rune
	transitions map[rune]*StateMachineNode
	funcs       []func(inputStr rune)
}

type StateMachineMethods interface {
	AddTransitions(inputStr rune, node *StateMachineNode)
	IsValidTransition(inputStr rune) bool
	Transition(inputStr rune) *StateMachineNode
	normalizeChar(inputStr rune) rune
}

func NewStateMachineNode(inputStr rune, funcs []func(input rune)) *StateMachineNode {
	return &StateMachineNode{
		char:        inputStr,
		transitions: map[rune]*StateMachineNode{},
		funcs:       funcs,
	}
}

func (node *StateMachineNode) AddTransitions(inputStr rune, newNode *StateMachineNode) {
	node.transitions[inputStr] = newNode
}

func (node *StateMachineNode) normalizeChar(input rune) rune {
	if node.char == '*' && !isSpecialChar(input) {
		return '*'
	}
	return input
}

func (node *StateMachineNode) IsValidTransition(inputStr rune) bool {
	_, ok := node.transitions['*']
	_, transitionExists := node.transitions[inputStr]

	if !isSpecialChar(inputStr) && ok {
		return true
	} else if node.char == '*' && !isSpecialChar(inputStr) {
		return true
	}
	return transitionExists
}

func isSpecialChar(input rune) bool {
	_, ok := reservedCharacters[input]
	if ok {
		return true
	} else {
		return false
	}
}

func isIgnoredCharacters(input rune) bool {
	_, ok := ignoredCharacters[input]
	if ok {
		return true
	} else {
		return false
	}
}

func (node *StateMachineNode) Transition(inputStr rune) *StateMachineNode {
	if isIgnoredCharacters(inputStr) {
		return node
	}
	if !node.IsValidTransition(inputStr) {
		return nil
	} else {
		newTransition := node.Transition(inputStr)
		if newTransition == nil {
			return nil
		}
		for i := 0; i < len(newTransition.funcs); i++ {

			newTransition.funcs[i](inputStr)
		}
	}

}
