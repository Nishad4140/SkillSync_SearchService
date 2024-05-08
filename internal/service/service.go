package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/Nishad4140/SkillSync_SearchService/internal/adapter"
	helperstruct "github.com/Nishad4140/SkillSync_SearchService/internal/helperStruct"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SearchService struct {
	UserConn pb.UserServiceClient
	adapters adapter.AdapterInterface
	pb.UnimplementedSearchServiceServer
}

func NewSearchService(adapters adapter.AdapterInterface, useraddr string) *SearchService {
	userConn, _ := grpc.Dial(useraddr, grpc.WithInsecure())
	return &SearchService{
		adapters: adapters,
		UserConn: pb.NewUserServiceClient(userConn),
	}
}

func (review *SearchService) UserAddReview(ctx context.Context, req *pb.UserReviewRequest) (*emptypb.Empty, error) {
	user, err := review.UserConn.GetClientById(context.Background(), &pb.GetUserById{
		Id: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	check, err := review.adapters.GetReviewCheck(req.UserId, req.FreelancerId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	if check {
		return nil, fmt.Errorf("this user has already entered a review for the given project please update the existing one or add a new one")
	}
	reqEntity := helperstruct.ReviewHelper{
		UserId:       req.UserId,
		FreelancerId: req.FreelancerId,
		ProjectId:    req.ProjectId,
		Rating:       int(req.Rating),
		Username:     user.Name,
		Description:  req.Description,
	}
	if err := review.adapters.UserAddReview(reqEntity); err != nil {
		return nil, err
	}

	avgRating, err := review.adapters.GetAverageRatingOfFreelancer(req.FreelancerId)
	if err != nil {
		log.Print("error while getting average rating", err)
	}

	_, err = review.UserConn.UpdateAverageRatingOfFreelancer(context.Background(), &pb.UpdateRatingRequest{
		FreelancerId: req.FreelancerId,
		AvgRating:    float32(avgRating),
	})
	if err != nil {
		log.Print("error while updating freelancer rating", err)
	}

	return nil, nil
}

func (r *SearchService) GetReview(req *pb.ReviewById, srv pb.SearchService_GetReviewServer) error {
	reviews, err := r.adapters.GetReviewsByFreelancer(req.Id)
	if err != nil {
		return err
	}
	for _, review := range reviews {
		userId, ok := review["userId"]
		if !ok {
			return fmt.Errorf("userId field not present")
		}
		userName, ok := review["username"]
		if !ok {
			return fmt.Errorf("username not present")
		}
		freelancerId, ok := review["freelancerId"]
		if !ok {
			return fmt.Errorf("username not present")
		}
		projectId, ok := review["projectId"]
		if !ok {
			return fmt.Errorf("username not present")
		}
		rating, ok := review["rating"]
		if !ok {
			return fmt.Errorf("rating not present")
		}
		des, ok := review["description"]
		if !ok {
			return fmt.Errorf("desription not found")
		}
		res := &pb.ReviewResponse{
			UserId:       userId.(string),
			FreelancerId: freelancerId.(string),
			Projectid:    projectId.(string),
			Description:  des.(string),
			Username:     userName.(string),
			Rating:       rating.(int32),
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (review *SearchService) RemoveReview(ctx context.Context, req *pb.UserReviewRequest) (*emptypb.Empty, error) {
	if err := review.adapters.UserDeleteReview(req.UserId, req.FreelancerId, req.ProjectId); err != nil {
		return nil, err
	}
	avgRating, err := review.adapters.GetAverageRatingOfFreelancer(req.FreelancerId)
	if err != nil {
		log.Print("error while getting average rating ", err)
	}
	if _, err := review.UserConn.UpdateAverageRatingOfFreelancer(context.Background(), &pb.UpdateRatingRequest{
		FreelancerId: req.FreelancerId,
		AvgRating: float32(avgRating),
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
