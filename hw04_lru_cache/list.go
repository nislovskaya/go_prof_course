package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	size int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.size == 0 {
		l.head = newItem
		l.tail = newItem
	} else {
		newItem.Next = l.head
		l.head.Prev = newItem
		l.head = newItem
	}

	l.size++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.size == 0 {
		l.head = newItem
		l.tail = newItem
	} else {
		newItem.Prev = l.tail
		l.tail.Next = newItem
		l.tail = newItem
	}

	l.size++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	l.size--

	if l.size == 0 {
		l.head = nil
		l.tail = nil
	}

	i.Next = nil
	i.Prev = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head || i == nil {
		return
	}

	l.Remove(i)

	i.Next = l.head
	i.Prev = nil

	if l.head != nil {
		l.head.Prev = i
	}

	l.head = i

	if l.size == 1 {
		l.tail = i
	}

	l.size++
}
