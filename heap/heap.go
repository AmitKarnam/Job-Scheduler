package heap

import (
	"fmt"

	"github.com/AmitKarnam/Job-Scheduler/models"
)

type Heap interface {
	Parent(index int) int
	LeftChild(index int) int
	RightChild(index int) int
	Insert(job models.Job)
	DeleteMin() models.Job
	Peak() models.Job
	HeapifyDown(index int)
}

// Data type definition
type jobHeap struct {
	arr []models.Job
}

func NewHeap() Heap {
	return &jobHeap{
		arr: make([]models.Job, 0),
	}
}

// Data type methods
func (h *jobHeap) Parent(index int) int {
	return ((index - 1) / 2)
}

func (h *jobHeap) LeftChild(index int) int {
	return ((index * 2) + 1)
}

func (h *jobHeap) RightChild(index int) int {
	return ((index * 2) + 2)
}

func (h *jobHeap) Swap(j1, j2 int) {
	h.arr[j1], h.arr[j2] = h.arr[j2], h.arr[j1]
}

func (h *jobHeap) Insert(job models.Job) {
	h.arr = append(h.arr, job)

	curIndex := len(h.arr) - 1

	for (curIndex > 0) && (job.NextExecutionTime.Before(h.arr[h.Parent(curIndex)].NextExecutionTime)) {
		h.Swap(h.Parent(curIndex), curIndex)
		curIndex = h.Parent(curIndex)
	}
	h.Swap(h.Parent(curIndex), curIndex)
	curIndex = h.Parent(curIndex)
}

func (h *jobHeap) DeleteMin() models.Job {
	n := len(h.arr)

	if n == 0 {
		fmt.Println("Empty Heap, nothing to delete")
		return models.Job{}
	}

	if n == 1 {
		min := h.arr[0]
		h.arr = []models.Job{} // reset heap
		return min
	}

	min := h.arr[0]
	h.arr[0] = h.arr[n-1]
	h.arr = h.arr[:n-1]

	h.HeapifyDown(0)
	return min
}

func (h *jobHeap) HeapifyDown(index int) {
	l := h.LeftChild(index)
	r := h.RightChild(index)
	smallest := index

	if l < len(h.arr) && h.arr[l].NextExecutionTime.Before(h.arr[smallest].NextExecutionTime) {
		smallest = l
	}

	if r < len(h.arr) && h.arr[r].NextExecutionTime.Before(h.arr[smallest].NextExecutionTime) {
		smallest = r
	}

	if smallest != index {
		h.Swap(index, smallest)
		h.HeapifyDown(smallest)
	}
}

func (h *jobHeap) Peak() models.Job {
	if len(h.arr) == 0 {
		return models.Job{}
	}
	return h.arr[0]
}
