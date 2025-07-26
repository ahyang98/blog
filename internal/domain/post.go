package domain

type Post struct {
	Id           uint
	Title        string
	Content      string
	Author       Author
	CommentCount uint
}

type Author struct {
	Id   uint
	Name string
}
