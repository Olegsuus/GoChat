package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

func (s *ServicesUser) HandleGoogleUser(ctx context.Context, userInfo models.GoogleUserInfo) (*models.User, error) {
	const op = "services.HandleGoogleUser"

	user, err := s.sP.Get(ctx, userInfo.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			newUser := &models.User{
				Name:      userInfo.GivenName,
				LastName:  userInfo.FamilyName,
				Email:     userInfo.Email,
				GoogleID:  userInfo.ID,
				Picture:   userInfo.Picture,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			userID, err := s.sP.Add(ctx, newUser)
			if err != nil {
				s.l.Error("Ошибка при создании пользователя", err.Error(), slog.String("op", op))
				return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
			}
			newUser.ID = userID
			s.l.Info("Пользователь успешно создан через Google", slog.String("email", newUser.Email), slog.String("op", op))
			return newUser, nil
		}
		s.l.Error("Ошибка при получении пользователя", err.Error(), slog.String("op", op))
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	updated := false
	if user.GoogleID == "" {
		user.GoogleID = userInfo.ID
		updated = true
	}

	if user.Picture != user.Picture {
		user.Picture = userInfo.Picture
		updated = true
	}

	if updated {
		user.UpdatedAt = time.Now()
		updateData := map[string]interface{}{
			"google_id":  user.GoogleID,
			"picture":    user.Picture,
			"updated_at": user.UpdatedAt,
		}
		err = s.sP.UpdateProfile(ctx, user.ID, updateData)
		if err != nil {
			s.l.Error("Ошибка при обновлении пользователя", err.Error(), slog.String("op", op))
			return nil, fmt.Errorf("ошибка при обновлении пользователя: %w", err)
		}
	}

	s.l.Info("Пользователь найден", slog.String("email", user.Email), slog.String("op", op))
	return user, nil
}
