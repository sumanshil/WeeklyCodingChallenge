package trie

import "testing"

func TestInsert(t *testing.T) {
	trie := &TrieNode{}
	trie.add("abcdef", "abc")
	data := trie.get("abcdef")
	v, ok := data.(string)
	if !ok {
		t.Fail()
	}
	if v != "abc" {
		t.Fail()
	}
}

func main() {

}
