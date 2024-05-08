package initializer

import (
	"github.com/Nishad4140/SkillSync_SearchService/internal/adapter"
	"github.com/Nishad4140/SkillSync_SearchService/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func Initializer(mongoDB *mongo.Database) *service.SearchService {
	adapter := adapter.NewSearchAdapter(mongoDB)
	service := service.NewSearchService(adapter, "ss-user-service:4001")
	return service
}
