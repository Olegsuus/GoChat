package services

import (
	"context"
	"fmt"
	"github.com/Olegsuus/Auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"time"
)

func (s *ServicesUser) ResetPassword(ctx context.Context, email, secretWord, newPassword string) error {
	const op = "services_reset_password"
	s.l.With(slog.String("op", op))

	user, err := s.sP.Get(ctx, email)
	if err != nil {
		s.l.Error("Пользователь с такой почтой не найден", slog.String("email", email), slog.String("op", op))
		return fmt.Errorf("пользователь не найден")
	}

	if user.SecretWord == "" {
		s.l.Error("Секретное слово не задано", slog.String("email", email), slog.String("op", op))
		return fmt.Errorf("секретное слово не задано, восстановление пароля невозможно")
	}

	if user.SecretWord != secretWord {
		s.l.Error("Неверное секретное слово", slog.String("email", email), slog.String("op", op))
		return fmt.Errorf("неверное секретное слово")
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		s.l.Error("Ошибка при хешировании пароля", err.Error(), slog.String("op", op))
		return fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	err = s.sP.UpdatePassword(ctx, user.ID, hashedPassword)
	if err != nil {
		s.l.Error("Ошибка при обновлении пароля", err.Error(), slog.String("op", op))
		return fmt.Errorf("ошибка при обновлении пароля: %w", err)
	}

	s.l.Info("Пароль успешно обновлен", slog.String("email", email), slog.String("op", op))
	return nil
}

func (s *ServicesUser) UpdateProfile(ctx context.Context, id primitive.ObjectID, dto models.UpdateUserDTO) error {
	const op = "services_update_profile"
	s.l.With(slog.String("op", op))

	updateData := bson.M{}

	if dto.Name != nil {
		updateData["name"] = *dto.Name
	}
	if dto.LastName != nil {
		updateData["last name"] = *dto.LastName
	}
	if dto.Email != nil {
		updateData["email"] = *dto.Email
	}
	if dto.Password != nil {
		hashedPassword, err := HashPassword(*dto.Password)
		if err != nil {
			s.l.Error("Ошибка при хешировании пароля", err.Error(), slog.String("op", op))
			return fmt.Errorf("ошибка при хешировании пароля: %w", err)
		}
		updateData["password"] = hashedPassword
	}
	if dto.PhoneNumber != nil {
		updateData["phone_number"] = *dto.PhoneNumber
	}
	if dto.Country != nil {
		updateData["country"] = *dto.Country
	}
	if dto.City != nil {
		updateData["city"] = *dto.City
	}
	if len(updateData) == 0 {
		return fmt.Errorf("нет данных для обновления")
	}
	updateData["updated_at"] = time.Now()

	err := s.sP.UpdateProfile(ctx, id, updateData)
	if err != nil {
		s.l.Error("Ошибка при обновлении профиля пользователя", err.Error(), slog.String("op", op))
		return fmt.Errorf("ошибка при обновлении профиля пользователя: %w", err)
	}

	s.l.Info("Профиль пользователя успешно обновлен", slog.String("userID", id.Hex()), slog.String("op", op))
	return nil
}
