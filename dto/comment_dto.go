package dto

import "time"

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
	PostID  uint   `json:"post_id" binding:"required"`
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int64             `json:"total"`
}
