package adapter

import (
	helperstruct "github.com/Nishad4140/SkillSync_SearchService/internal/helperStruct"
	"go.mongodb.org/mongo-driver/bson"
)

type AdapterInterface interface {
	UserAddReview(req helperstruct.ReviewHelper) error
	GetReviewsByFreelancer(freelancerId string) ([]bson.M, error)
	UserDeleteReview(userId, freelancerId, projectId string) error
	GetReviewCheck(userId, freelancerId, projectId string) (bool, error)
	GetAverageRatingOfFreelancer(freelancerId string) (float64, error)
}
