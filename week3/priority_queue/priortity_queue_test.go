package priority_queue

import "testing"

type IntNode int

func (i IntNode) Value() int {
	return int(i)
}

func TestInsertSortedInput(t *testing.T) {
	queue := NewPriorityQueue()
	queue.Insert(IntNode(4))
	queue.Insert(IntNode(5))
	if queue.getTop().Value() != 4 {
		t.Fail()
	}
	node := queue.Delete()
	if node.Value() != 4 {
		t.Fail()
	}
	if queue.getTop().Value() != 5 {
		t.Fail()
	}
	node = queue.Delete()
	if node.Value() != 5 {
		t.Fail()
	}
}

func TestInsertUnSortedInput(t *testing.T) {
	queue := NewPriorityQueue()
	queue.Insert(IntNode(5))
	queue.Insert(IntNode(4))
	if queue.getTop().Value() != 4 {
		t.Fail()
	}
	node := queue.Delete()
	if node.Value() != 4 {
		t.Fail()
	}
	if queue.getTop().Value() != 5 {
		t.Fail()
	}
	node = queue.Delete()
	if node.Value() != 5 {
		t.Fail()
	}
}

func main() {

}
