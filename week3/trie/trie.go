package trie

type TrieNode struct {
	c        rune
	children map[rune]*TrieNode
	isLeaf   bool
	data     interface{}
}

func (trieNode *TrieNode) add(input string, data string) {
	currentNode := trieNode
	for i := 0; i < len(input); i++ {
		c := input[i]
		if currentNode.children == nil {
			currentNode.children = make(map[rune]*TrieNode)
		}
		v, ok := currentNode.children[rune(c)]
		if ok {
			currentNode = v
		} else {
			currentNode.children[rune(c)] = &TrieNode{c: rune(c)}
		}
		currentNode = currentNode.children[rune(c)]
	}
	currentNode.isLeaf = true
	currentNode.data = data
}

func (trieNode *TrieNode) get(input string) interface{} {
	currentNode := trieNode
	for i := 0; i < len(input); i++ {
		c := input[i]
		v, ok := currentNode.children[rune(c)]
		if !ok {
			return nil
		}
		currentNode = v
	}

	if currentNode.isLeaf {
		return currentNode.data
	} else {
		return nil
	}
}
