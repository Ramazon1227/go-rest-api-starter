package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/email"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/utils"
	"github.com/Ramazon1227/go-rest-api-starter/storage"
)

const collectionName = "user"

type userRepoImpl struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) storage.UserRepoImpl {
	return &userRepoImpl{db: db}
}

func (r *userRepoImpl) coll() *mongo.Collection {
	return r.db.Collection(collectionName)
}

func (r *userRepoImpl) Add(ctx context.Context, entity *models.UserCreateModel) (*models.PrimaryKey, error) {
	coll := r.coll()

	// Check if a user with this email already exists (possibly soft-deleted).
	var existing models.User
	err := coll.FindOne(ctx, bson.M{"email": entity.Email}).Decode(&existing)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == nil && existing.DeletedAt != nil {
		// Reactivate the soft-deleted user.
		now := time.Now()
		_, err = coll.UpdateOne(ctx,
			bson.M{"_id": existing.Id},
			bson.M{"$set": bson.M{
				"name":       entity.Name,
				"role":       entity.Role,
				"phone":      entity.Phone,
				"deleted_at": nil,
				"updated_at": now,
			}},
		)
		if err != nil {
			return nil, err
		}
		return &models.PrimaryKey{Id: existing.Id}, nil
	}

	plainPassword, err := utils.GenerateRandomPassword(8)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	now := time.Now()
	expires := time.Date(2072, 5, 1, 11, 21, 59, 0, time.UTC)

	doc := models.User{
		Id:        id,
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  hashedPassword,
		Role:      entity.Role,
		Phone:     entity.Phone,
		Active:    1,
		ExpiresAt: &expires,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	_, err = coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	if err2 := email.SendWelcomeEmail(entity.Email, entity.Name, plainPassword); err2 != nil {
		log.Printf("Failed to send welcome email: %v", err2)
	}

	return &models.PrimaryKey{Id: id}, nil
}

func (r *userRepoImpl) UpdateProfile(ctx context.Context, entity *models.UpdateUserProfileModel) error {
	now := time.Now()
	result, err := r.coll().UpdateOne(ctx,
		bson.M{"_id": entity.Id, "deleted_at": nil},
		bson.M{"$set": bson.M{
			"name":       entity.Name,
			"email":      entity.Email,
			"updated_at": now,
		}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return storage.ErrorNotFound
	}
	return nil
}

func (r *userRepoImpl) GetById(ctx context.Context, pKey *models.PrimaryKey) (*models.User, error) {
	var user models.User
	err := r.coll().FindOne(ctx, bson.M{"_id": pKey.Id, "deleted_at": nil}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, storage.ErrorNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepoImpl) GetList(ctx context.Context, queryParam *models.QueryParam) (*models.GetUserListModel, error) {
	filter := bson.M{"deleted_at": nil}

	count, err := r.coll().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(queryParam.Limit)).
		SetSkip(int64(queryParam.Offset))

	cursor, err := r.coll().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users = make([]*models.User, 0,queryParam.Limit)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &models.GetUserListModel{
		Users: users,
		Count: int(count),
	}, nil
}

func (r *userRepoImpl) Delete(ctx context.Context, pKey *models.PrimaryKey) error {
	now := time.Now()
	result, err := r.coll().UpdateOne(ctx,
		bson.M{"_id": pKey.Id, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return storage.ErrorNotFound
	}
	return nil
}

func (r *userRepoImpl) GetByEmail(ctx context.Context, emailAddr string) (*models.User, error) {
	var user models.User
	err := r.coll().FindOne(ctx, bson.M{"email": emailAddr}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, storage.ErrorNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepoImpl) UpdateUserProfile(ctx context.Context, userId string, req *models.UpdateProfileRequest) error {
	set := bson.M{"updated_at": time.Now()}

	if req.Name != "" {
		set["name"] = req.Name
	}
	if req.Phone != "" {
		set["phone"] = req.Phone
	}
	if req.Email != "" {
		set["email"] = req.Email
	}

	if len(set) == 1 { // only updated_at — nothing to do
		return nil
	}

	result, err := r.coll().UpdateOne(ctx,
		bson.M{"_id": userId},
		bson.M{"$set": set},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return storage.ErrorNotFound
	}
	return nil
}

func (r *userRepoImpl) UpdatePassword(ctx context.Context, userId string, currentPassword, newPassword string) error {
	user, err := r.GetById(ctx, &models.PrimaryKey{Id: userId})
	if err != nil {
		return err
	}

	if !utils.CheckPassword(user.Password, currentPassword) {
		return fmt.Errorf("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	result, err := r.coll().UpdateOne(ctx,
		bson.M{"_id": userId},
		bson.M{"$set": bson.M{
			"password":   hashedPassword,
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return storage.ErrorNotFound
	}
	return nil
}

