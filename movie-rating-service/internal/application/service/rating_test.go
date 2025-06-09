//go:build unit_test

package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/domain"
	"movie-rating-service/mocks"
	"testing"
)

type RatingServiceTest struct {
	suite.Suite
	service ratingService
	r       *mocks.RatingRepository
	m       *mocks.MovieRepository
}

func (r *RatingServiceTest) SetupTest() {
	r.r = new(mocks.RatingRepository)
	r.m = new(mocks.MovieRepository)

	r.service = ratingService{ratingRepository: r.r, movieRepository: r.m}
}

func Test_RunRatingServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RatingServiceTest))
}

func (r *RatingServiceTest) TestPromotionService_Create_Success() {
	t := r.T()

	ctx := context.TODO()

	req := request.CreateRating{
		UserID:  1,
		MovieID: 42,
		Score:   4.7,
		Review:  "good",
	}

	rating := &domain.Rating{
		UserID:  req.UserID,
		MovieID: req.MovieID,
		Score:   req.Score,
		Review:  req.Review,
	}

	r.r.On("Create", ctx, mock.MatchedBy(func(r domain.Rating) bool {
		return r.UserID == req.UserID && r.MovieID == req.MovieID && r.Score == req.Score && r.Review == req.Review
	})).Return(rating, nil).Once()

	r.m.On("AddRating", ctx, req.MovieID, req.Score).Return(nil).Once()

	result, err := r.service.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	r.r.AssertExpectations(t)
	r.m.AssertExpectations(t)
}

func (r *RatingServiceTest) TestPromotionService_Create_Error_Failed_To_Create_Rating() {
	t := r.T()

	ctx := context.TODO()

	req := request.CreateRating{
		UserID:  1,
		MovieID: 42,
		Score:   4.7,
		Review:  "good",
	}

	r.r.On("Create", ctx, mock.MatchedBy(func(r domain.Rating) bool {
		return r.UserID == req.UserID && r.MovieID == req.MovieID && r.Score == req.Score && r.Review == req.Review
	})).Return(nil, errors.New("there is an error")).Once()

	result, err := r.service.Create(ctx, req)

	assert.ErrorContains(t, err, "failed to rate movie: there is an error")
	assert.Nil(t, result)

	r.r.AssertExpectations(t)
	r.m.AssertExpectations(t)
}

func (r *RatingServiceTest) TestPromotionService_Create_Error_Failed_To_Update_Movie() {
	t := r.T()

	ctx := context.TODO()

	req := request.CreateRating{
		UserID:  1,
		MovieID: 42,
		Score:   4.7,
		Review:  "good",
	}

	rating := &domain.Rating{
		UserID:  req.UserID,
		MovieID: req.MovieID,
		Score:   req.Score,
		Review:  req.Review,
	}

	r.r.On("Create", ctx, mock.MatchedBy(func(r domain.Rating) bool {
		return r.UserID == req.UserID && r.MovieID == req.MovieID && r.Score == req.Score && r.Review == req.Review
	})).Return(rating, nil).Once()

	r.m.On("AddRating", ctx, req.MovieID, req.Score).Return(errors.New("there is an error")).Once()

	result, err := r.service.Create(ctx, req)

	assert.ErrorContains(t, err, "failed to update rating: there is an error")
	assert.Nil(t, result)

	r.r.AssertExpectations(t)
	r.m.AssertExpectations(t)
}

func (r *RatingServiceTest) TestPromotionService_GetRatingsByUserID_Success() {
	t := r.T()

	ctx := context.TODO()

	req := request.GetUserRatings{
		UserID: 1,
	}

	// Prepare fake ratings
	domainRatings := []domain.Rating{
		{
			UserID:  req.UserID,
			MovieID: 42,
			Score:   4.7,
			Review:  "good",
		},
		{
			UserID:  req.UserID,
			MovieID: 43,
			Score:   3.2,
			Review:  "average",
		},
	}

	r.r.On("GetByUserID", ctx, req.UserID).Return(domainRatings, nil).Once()

	result, err := r.service.GetRatingsByUserID(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Ratings, 2)
	assert.Equal(t, domainRatings[0].Movie.Title, result.Ratings[0].RatedMovie.Title)
	assert.Equal(t, domainRatings[1].Score, result.Ratings[1].Rating.Score)

	r.r.AssertExpectations(t)
}

func (r *RatingServiceTest) TestPromotionService_GetRatingsByUserID_Error_Failed_To_Get_User_Ratings() {
	t := r.T()

	ctx := context.TODO()

	req := request.GetUserRatings{
		UserID: 1,
	}

	r.r.On("GetByUserID", ctx, req.UserID).Return(nil, errors.New("there is an error")).Once()

	result, err := r.service.GetRatingsByUserID(ctx, req)

	assert.ErrorContains(t, err, "error occurred while getting user's ratings: there is an error")
	assert.Nil(t, result)

	r.r.AssertExpectations(t)
}
