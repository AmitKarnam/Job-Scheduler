package heap

type Heap interface {
	Parent(index int) int
	LeftChild(index int) int
	RightChild(index int) int
	Insert()
	DeleteMin()
	Delete()
	Peak()
	HeapifyDown()
	HeapifyUp()
}

// Data type definition
type heapImp struct {
	arr []int
}

func NewHeap() Heap {
	return &heapImp{
		arr: make([]int, 0),
	}
}

// Data type methods
func (h *heapImp) Parent(index int) int {
	return ((index - 1) / 2)
}

func (h *heapImp) LeftChild(index int) int {
	return ((index * 2) + 1)
}

func (h *heapImp) RightChild(index int) int {
	return ((index * 2) + 2)
}

func (h *heapImp) Insert() {

}

func (h *heapImp) DeleteMin() {

}

func (h *heapImp) Delete() {

}

func (h *heapImp) Peak() {

}

func (h *heapImp) HeapifyUp() {

}

func (h *heapImp) HeapifyDown() {

}
