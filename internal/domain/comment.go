package domain

type Comment struct {
	Id      uint
	Content string
	Author  Author
	PostId  uint
}
