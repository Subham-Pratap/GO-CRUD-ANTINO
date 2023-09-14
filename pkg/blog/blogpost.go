package blog

// BlogPost is the struct for a blog post
type BlogPost struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
