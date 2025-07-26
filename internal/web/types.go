package web

const (
	ErrSystemError = "system error"
)

type PostVo struct {
	Id           uint   `json:"id,omitempty"`
	Title        string `json:"title,omitempty"`
	Content      string `json:"content,omitempty"`
	AuthorId     uint   `json:"author_id,omitempty"`
	AuthorName   string `json:"author_name,omitempty"`
	CommentCount uint   `json:"comment_count,omitempty"`
}
