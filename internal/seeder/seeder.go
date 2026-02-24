package seeder

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/config"
	"github.com/zagvozdeen/ola/internal/logger"
	"github.com/zagvozdeen/ola/internal/store"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
	cfg   *config.Config
	log   *logger.Logger
	store *store.Store
}

func New(cfg *config.Config, log *logger.Logger, store *store.Store) *Seeder {
	return &Seeder{
		cfg:   cfg,
		log:   log,
		store: store,
	}
}

func (s *Seeder) Run(ctx context.Context) error {
	if !s.cfg.App.RunSeeder {
		s.log.Info("Seeder disabled")
		return nil
	}

	user, err := s.store.GetUserByTID(ctx, s.cfg.Root.TID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return fmt.Errorf("failed to get user by tid: %w", err)
		}
		var password []byte
		password, err = bcrypt.GenerateFromPassword([]byte(s.cfg.Root.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to generate hashed password: %w", err)
		}
		user = &models.User{
			TID:       new(s.cfg.Root.TID),
			UUID:      s.cfg.Root.UUID,
			FirstName: s.cfg.Root.FirstName,
			LastName:  new(s.cfg.Root.LastName),
			Username:  new(s.cfg.Root.Username),
			Email:     new(s.cfg.Root.Email),
			Phone:     new(s.cfg.Root.Phone),
			Password:  new(string(password)),
			Role:      enums.UserRoleAdmin,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.store.CreateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	var uid uuid.UUID
	for i := range 9 {
		_, err = s.store.GetFileByID(ctx, i+1)
		if err != nil {
			if !errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("failed to get file by id: %w", err)
			}
			uid, err = uuid.NewUUID()
			if err != nil {
				return fmt.Errorf("failed to create uuid: %w", err)
			}
			err = s.store.CreateFile(ctx, &models.File{
				UUID:       uid,
				Content:    fmt.Sprintf("/files/%d.jpg", i+1),
				Size:       0,
				MimeType:   "image/jpeg",
				OriginName: fmt.Sprintf("%d.jpg", i+1),
				UserID:     user.ID,
				CreatedAt:  time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
		}
	}

	categories := []string{"Детские праздники", "Корпоратив", "День рождение", "Свадьба", "Выписка", "Гендер пати", "8 марта", "23 февраля", "Без повода"}
	for _, category := range categories {
		_, err = s.store.GetCategoryByName(ctx, category)
		if err != nil {
			if !errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("failed to get category by name: %w", err)
			}
			uid, err = uuid.NewUUID()
			if err != nil {
				return fmt.Errorf("failed to create uuid: %w", err)
			}
			err = s.store.CreateCategory(ctx, &models.Category{
				UUID:      uid,
				Name:      category,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	products := []struct {
		from        int
		to          *int
		name        string
		description string
	}{
		{from: 3500, to: nil, name: "Фонтан из воздушных шаров", description: "Композиция по индивидуальному дизайну для любого события"},
		{from: 7000, to: nil, name: "Оформление помещения / фотозона", description: "Декорирование любого помещения по индивидуальному дизайну"},
		{from: 5000, to: nil, name: "Коробка - сюрприз", description: "Подарочный бокс с композицией из шаров для любого события"},
		{from: 3000, to: nil, name: "Бабл бокс", description: "Креативная упаковка для небольшого подарка с шаром баблс"},
	}
	var counter int
	for _, t := range []enums.ProductType{enums.ProductTypeProduct, enums.ProductTypeService} {
		for i, p := range products {
			counter++
			_, err = s.store.GetProductByID(ctx, counter)
			if err != nil {
				if !errors.Is(err, models.ErrNotFound) {
					return fmt.Errorf("failed to get product by id: %w", err)
				}
				uid, err = uuid.NewUUID()
				if err != nil {
					return fmt.Errorf("failed to create uuid: %w", err)
				}
				err = s.store.CreateProduct(ctx, &models.Product{
					UUID:        uid,
					Name:        p.name,
					Description: p.description,
					PriceFrom:   p.from,
					PriceTo:     p.to,
					Type:        t,
					FileID:      i + 1,
					UserID:      user.ID,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				if err != nil {
					return fmt.Errorf("failed to create product: %w", err)
				}
			}
		}
	}

	reviews := []struct {
		name    string
		content string
	}{
		{name: "Елена", content: "Огромное спасибо за шарики! Именинник был в восторге и доставка порадовала, все вовремя) Не пожалела, что обратилась именно к вам!"},
		{name: "Евгений", content: "Обратился за бабл боксом с украшением внутри, хотелось поздравить девушку с годовщиной. Она очень удивилась такой креативной идее. Я тоже не видел ничего подобного в нашем городе до этого)) Желаю вам дальнейшего развития, вы классные!"},
		{name: "Александр", content: "Выражаю благодарность студии за профессиональное оформление нашего корпоратива. Требовался лаконичный декор в цветах компании. Результат превзошел ожидания: композиции у входа и фотозона были выполнены безупречно, с вниманием к деталям. Отдельно отмечу пунктуальность, четкое соблюдение сроков и договоренностей."},
		{name: "Любовь", content: "Благодарю персонал за то, что взялись за очень срочный заказ, выполнили и доставили максимально быстро. Посоветую вас друзьям и сама обращусь еще не раз."},
		{name: "Анна", content: "Как человек, который ценит визуал, долго искала в Екатеринбурге студию, которая умеет в тренды. Команда Ola предложила крутые цветовые сочетания для моей вечеринки. Шары держались несколько дней, не теряя вид! Ни одного лопнувшего. Это показатель. Если вам важен дизайн, атмосфера и стойкость - выбор очевиден."},
	}
	for i, r := range reviews {
		if err != nil {
			return fmt.Errorf("failed to create uuid: %w", err)
		}
		_, err = s.store.GetReviewByID(ctx, i+1)
		if err != nil {
			if !errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("failed to get review by id: %w", err)
			}
			err = s.store.CreateReview(ctx, &models.Review{
				UUID:        uid,
				Name:        r.name,
				Content:     r.content,
				FileID:      i + 5,
				UserID:      user.ID,
				PublishedAt: time.Now(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to create review: %w", err)
			}
		}
	}

	return nil
}
