package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	TicketReviewRepo struct {
		entryRepo
	}

	InterTicketReviewRepo interface {
		// 评审操作
		CreateReview(review models.TicketReview) error
		UpdateReview(review models.TicketReview) error
		DeleteReview(reviewId string) error
		GetReview(reviewId string) (models.TicketReview, error)
		ListReviews(ticketId string, reviewerId string, status models.TicketReviewStatus, page, size int) ([]models.TicketReview, int64, error)
	}
)

func newTicketReviewInterface(db *gorm.DB, g InterGormDBCli) InterTicketReviewRepo {
	return &TicketReviewRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

// CreateReview 创建评审
func (tr TicketReviewRepo) CreateReview(review models.TicketReview) error {
	return tr.g.Create(&models.TicketReview{}, &review)
}

// UpdateReview 更新评审
func (tr TicketReviewRepo) UpdateReview(review models.TicketReview) error {
	return tr.g.Updates(Updates{
		Table:   &models.TicketReview{},
		Where:   map[string]interface{}{"review_id": review.ReviewId},
		Updates: review,
	})
}

// DeleteReview 删除评审
func (tr TicketReviewRepo) DeleteReview(reviewId string) error {
	return tr.g.Delete(Delete{
		Table: &models.TicketReview{},
		Where: map[string]interface{}{"review_id": reviewId},
	})
}

// GetReview 获取评审详情
func (tr TicketReviewRepo) GetReview(reviewId string) (models.TicketReview, error) {
	var review models.TicketReview
	db := tr.db.Model(&models.TicketReview{})
	db.Where("review_id = ?", reviewId)
	err := db.First(&review).Error
	if err != nil {
		return review, err
	}
	return review, nil
}

// ListReviews 获取评审列表
func (tr TicketReviewRepo) ListReviews(ticketId string, reviewerId string, status models.TicketReviewStatus, page, size int) ([]models.TicketReview, int64, error) {
	var (
		reviews []models.TicketReview
		count   int64
	)

	db := tr.db.Model(&models.TicketReview{})

	if ticketId != "" {
		db.Where("ticket_id = ?", ticketId)
	}
	if reviewerId != "" {
		db.Where("reviewer_id = ?", reviewerId)
	}
	if status != "" {
		db.Where("status = ?", status)
	}

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("created_at DESC")

	err := db.Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, count, nil
}
