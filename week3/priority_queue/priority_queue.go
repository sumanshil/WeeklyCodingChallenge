package priority_queue

type PriorityQueue struct {
	nodes []PriorityQueueNode
}

type PriorityQueueNode interface {
	Value() int
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		nodes: make([]PriorityQueueNode, 0),
	}
}

func (queue *PriorityQueue) Size() int {
	return len(queue.nodes)
}

func (queue *PriorityQueue) Insert(node PriorityQueueNode) {
	queue.nodes = append(queue.nodes, node)
	queue.upHeap(len(queue.nodes) - 1)
}

func (queue *PriorityQueue) upHeap(index int) {
	if index == 0 {
		return
	}
	parentIndex := index / 2
	if queue.nodes[parentIndex].Value() > queue.nodes[index].Value() {
		queue.swap(index, parentIndex)
	}
	queue.upHeap(parentIndex)
}

func (queue *PriorityQueue) swap(index int, index2 int) {
	tempNode := queue.nodes[index]
	queue.nodes[index] = queue.nodes[index2]
	queue.nodes[index2] = tempNode
}

func (queue *PriorityQueue) getTop() PriorityQueueNode {
	if len(queue.nodes) == 0 {
		panic("queue is empty")
	}
	return queue.nodes[0]
}

func (queue *PriorityQueue) Delete() PriorityQueueNode {
	queue.swap(0, len(queue.nodes)-1)
	queue.downHeap(0)
	tempIndex := len(queue.nodes) - 1
	retVal := queue.nodes[tempIndex]
	queue.nodes = queue.nodes[:len(queue.nodes)-1]
	return retVal
}

func (queue *PriorityQueue) downHeap(index int) {
	if index == len(queue.nodes)-1 {
		return
	}
	smallerIndex := index
	leftChildIndex := (2 * index) + 1
	rightChildIndex := (2 * index) + 2
	if leftChildIndex < len(queue.nodes) && queue.nodes[leftChildIndex].Value() < queue.nodes[smallerIndex].Value() {
		smallerIndex = leftChildIndex
	}

	if rightChildIndex < len(queue.nodes) && queue.nodes[rightChildIndex].Value() < queue.nodes[smallerIndex].Value() {
		smallerIndex = rightChildIndex
	}
	if smallerIndex != index {
		queue.swap(index, smallerIndex)
		queue.downHeap(smallerIndex)
	}
}
