package handlers

import (
	"blog_system/database"
	"blog_system/dto"
	"blog_system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	userID := c.GetUint("userID")

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := database.DB.First(&post, req.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  req.PostID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	response := dto.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		CreatedAt: comment.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func GetCommentsByPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var comments []models.Comment
	if err := database.DB.Where("post_id = ?", postID).Order("created_at asc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	var commentResponses []dto.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, dto.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			PostID:    comment.PostID,
			CreatedAt: comment.CreatedAt,
		})
	}

	response := dto.CommentListResponse{
		Comments: commentResponses,
		Total:    int64(len(commentResponses)),
	}

	c.JSON(http.StatusOK, response)
}
