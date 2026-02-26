package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type ticketReviewService struct {
	ctx *ctx.Context
}

type InterTicketReviewService interface {
	// 评审操作
	AssignReviewers(req interface{}) (interface{}, interface{})
	SubmitReview(req interface{}) (interface{}, interface{})
	GetReview(req interface{}) (interface{}, interface{})
	ListReviews(req interface{}) (interface{}, interface{})

	// 评委操作
	CreateReviewer(req interface{}) (interface{}, interface{})
	UpdateReviewer(req interface{}) (interface{}, interface{})
	DeleteReviewer(req interface{}) (interface{}, interface{})
	GetReviewer(req interface{}) (interface{}, interface{})
	ListReviewers(req interface{}) (interface{}, interface{})
}

func newInterTicketReviewService(ctx *ctx.Context) InterTicketReviewService {
	return &ticketReviewService{ctx}
}

// AssignReviewers 分配评委
func (s ticketReviewService) AssignReviewers(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewAssign)

	// 验证工单是否存在
	ticket, err := s.ctx.DB.Ticket().Get(r.TenantId, r.TicketId)
	if err != nil {
		return nil, fmt.Errorf("工单不存在")
	}

	// 验证工单状态
	if ticket.Status != models.TicketStatusVerifying {
		return nil, fmt.Errorf("工单状态为 %s，不能分配评委。工单必须处于验证状态", ticket.Status)
	}

	// 为每个评委创建评审记录
	for _, reviewerId := range r.ReviewerIds {
		review := models.TicketReview{
			ReviewId:   "review-" + tools.RandId(),
			TenantId:   r.TenantId,
			TicketId:   r.TicketId,
			ReviewerId: reviewerId,
			Status:     models.ReviewStatusPending,
			CreatedBy:  r.AssignedBy,
			CreatedAt:  time.Now().Unix(),
		}

		err := s.ctx.DB.TicketReview().CreateReview(review)
		if err != nil {
			return nil, fmt.Errorf("创建评审记录失败: %v", err)
		}
	}

	return nil, nil
}

// SubmitReview 提交评审
func (s ticketReviewService) SubmitReview(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewSubmit)

	// 获取评审记录
	review, err := s.ctx.DB.TicketReview().GetReview(r.ReviewId)
	if err != nil {
		return nil, fmt.Errorf("评审记录不存在")
	}

	// 更新评审
	review.Rating = r.Rating
	review.WorkHours = r.WorkHours
	review.Comment = r.Comment
	review.Status = models.ReviewStatusCompleted
	review.CompletedAt = time.Now().Unix()

	err = s.ctx.DB.TicketReview().UpdateReview(review)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetReview 获取评审详情
func (s ticketReviewService) GetReview(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewQuery)

	review, err := s.ctx.DB.TicketReview().GetReview(r.TicketId)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// ListReviews 获取评审列表
func (s ticketReviewService) ListReviews(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewQuery)

	reviews, total, err := s.ctx.DB.TicketReview().ListReviews(r.TicketId, r.ReviewerId, r.Status, r.Page, r.Size)
	if err != nil {
		return nil, err
	}

	return types.ResponseTicketReviewList{
		List:  reviews,
		Total: total,
	}, nil
}

// CreateReviewer 创建评委
func (s ticketReviewService) CreateReviewer(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewerCreate)

	reviewer := models.TicketReviewer{
		TenantId:   r.TenantId,
		ReviewerId: r.ReviewerId,
		UserName:   r.UserName,
		Email:      r.Email,
		Phone:      r.Phone,
		Department: r.Department,
		Specialty:  r.Specialty,
		IsActive:   true,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	err := s.ctx.DB.TicketReview().CreateReviewer(reviewer)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateReviewer 更新评委
func (s ticketReviewService) UpdateReviewer(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewerUpdate)

	reviewer, err := s.ctx.DB.TicketReview().GetReviewer(r.TenantId, r.ReviewerId)
	if err != nil {
		return nil, fmt.Errorf("评委不存在")
	}

	if r.UserName != "" {
		reviewer.UserName = r.UserName
	}
	if r.Email != "" {
		reviewer.Email = r.Email
	}
	if r.Phone != "" {
		reviewer.Phone = r.Phone
	}
	if r.Department != "" {
		reviewer.Department = r.Department
	}
	if r.Specialty != "" {
		reviewer.Specialty = r.Specialty
	}
	if r.IsActive != nil {
		reviewer.IsActive = *r.IsActive
	}
	reviewer.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.TicketReview().UpdateReviewer(reviewer)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteReviewer 删除评委
func (s ticketReviewService) DeleteReviewer(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewerUpdate)

	err := s.ctx.DB.TicketReview().DeleteReviewer(r.TenantId, r.ReviewerId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetReviewer 获取评委详情
func (s ticketReviewService) GetReviewer(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewerQuery)

	reviewer, err := s.ctx.DB.TicketReview().GetReviewer(r.TenantId, r.ReviewerId)
	if err != nil {
		return nil, err
	}

	return reviewer, nil
}

// ListReviewers 获取评委列表
func (s ticketReviewService) ListReviewers(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestTicketReviewerQuery)

	reviewers, total, err := s.ctx.DB.TicketReview().ListReviewers(r.TenantId, r.Department, r.Specialty, r.IsActive, r.Page, r.Size)
	if err != nil {
		return nil, err
	}

	return types.ResponseTicketReviewerList{
		List:  reviewers,
		Total: total,
	}, nil
}
