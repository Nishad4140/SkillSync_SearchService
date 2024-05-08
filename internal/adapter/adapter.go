package adapter

import (
	"context"
	"fmt"
	"time"

	helperstruct "github.com/Nishad4140/SkillSync_SearchService/internal/helperStruct"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Adapter struct {
	Mongodb *mongo.Database
}

func NewSearchAdapter(mongodb *mongo.Database) *Adapter {
	return &Adapter{
		Mongodb: mongodb,
	}
}

func (r *Adapter) GetAverageRatingOfFreelancer(freelancerId string) (float64, error) {
	collection := r.Mongodb.Collection("userreview")
	if collection == nil {
		return 0, fmt.Errorf("collection is empty")
	}

	pipeline := bson.A{
		bson.M{"$match": bson.M{"freelancerId": freelancerId}},
		bson.M{"$group": bson.M{"_id": "$freelancerId", "averageRating": bson.M{"$avg": "$rating"}}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(context.Background())
	if !cursor.Next(context.Background()) {
		return 0, fmt.Errorf("no reviews found for freelancer with ID: %s", freelancerId)
	}

	var result bson.M
	if err := cursor.Decode(&result); err != nil {
		return 0, err
	}

	averageRating, ok := result["averageRating"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid average rating format")
	}

	return averageRating, nil
}

func (review *Adapter) GetReviewCheck(userId, freelancerId, projectId string) (bool, error) {
	collection := review.Mongodb.Collection("userreview")
	if collection == nil {
		return false, fmt.Errorf("collection is empty")
	}
	filter := bson.M{"userId": userId, "freelancerId": freelancerId, "projectId": projectId}
	res := collection.FindOne(context.Background(), filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, res.Err()
	}
	return true, nil
}

func (review *Adapter) GetReviewsByFreelancer(freelancerId string) ([]bson.M, error) {
	collection := review.Mongodb.Collection("userreview")
	if collection == nil {
		return nil, fmt.Errorf("collection not found")
	}
	filter := bson.M{"freelancerId": freelancerId}
	options := options.Find().SetSort(bson.D{{"timestamp", -1}})
	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var reviews []bson.M
	for cursor.Next(context.Background()) {
		var review bson.M
		err = cursor.Decode(&review)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (review *Adapter) UserAddReview(req helperstruct.ReviewHelper) error {
	collection := review.Mongodb.Collection("userreview")
	if collection == nil {
		err := review.Mongodb.CreateCollection(context.Background(), "userreview")
		if err != nil {
			return err
		}
		collection = review.Mongodb.Collection("userreview")
	}
	reviewDoc := bson.M{
		"userId":       req.UserId,
		"freelancerId": req.FreelancerId,
		"projectId":    req.ProjectId,
		"rating":       req.Rating,
		"username":     req.Username,
		"description":  req.Description,
		"timestamp":    time.Now(),
	}
	_, err := collection.InsertOne(context.Background(), reviewDoc)
	if err != nil {
		return err
	}
	return nil
}

func (review *Adapter) UserDeleteReview(userId, freelancerId, projectId string) error {
	collection := review.Mongodb.Collection("userreview")
	if collection == nil {
		return fmt.Errorf("collection is empty")
	}
	filter := bson.M{"userId": userId, "freelancerId": freelancerId, "projectId": projectId}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
