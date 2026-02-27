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

	var err error
	ctx, err = s.store.Begin(ctx)
	defer s.store.Rollback(ctx)

	var user *models.User
	user, err = s.store.GetUserByTID(ctx, s.cfg.Root.TID)
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

	files := []struct {
		id   int
		uuid uuid.UUID
	}{
		{uuid: uuid.MustParse("f3aaf3df-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab0dd5-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab14fe-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab1996-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab1e48-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab240c-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab2978-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab2cd4-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab308d-11af-11f1-b4af-c87f54a92045")},
	}
	for i, f := range files {
		var file *models.File
		file, err = s.store.GetFileByUUID(ctx, f.uuid)
		if err != nil {
			if !errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("failed to get file by id: %w", err)
			}
			file = &models.File{
				UUID:       f.uuid,
				Content:    fmt.Sprintf("/files/%d.jpg", i+1),
				Size:       0,
				MimeType:   "image/jpeg",
				OriginName: fmt.Sprintf("%d.jpg", i+1),
				UserID:     user.ID,
				CreatedAt:  time.Now(),
			}
			err = s.store.CreateFile(ctx, file)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
		}
		files[i].id = file.ID
	}

	categories := []struct {
		uuid uuid.UUID
		name string
	}{
		{uuid: uuid.MustParse("f3ab3f6c-11af-11f1-b4af-c87f54a92045"), name: "Детские праздники"},
		{uuid: uuid.MustParse("f3ab46b4-11af-11f1-b4af-c87f54a92045"), name: "Корпоратив"},
		{uuid: uuid.MustParse("f3ab4a71-11af-11f1-b4af-c87f54a92045"), name: "День рождение"},
		{uuid: uuid.MustParse("f3ab4e03-11af-11f1-b4af-c87f54a92045"), name: "Свадьба"},
		{uuid: uuid.MustParse("f3ab5261-11af-11f1-b4af-c87f54a92045"), name: "Выписка"},
		{uuid: uuid.MustParse("f3ab56ce-11af-11f1-b4af-c87f54a92045"), name: "Гендер пати"},
		{uuid: uuid.MustParse("f3ab5a58-11af-11f1-b4af-c87f54a92045"), name: "8 марта"},
		{uuid: uuid.MustParse("f3ab5e39-11af-11f1-b4af-c87f54a92045"), name: "23 февраля"},
		{uuid: uuid.MustParse("f3ab61f8-11af-11f1-b4af-c87f54a92045"), name: "Без повода"},
	}
	for _, category := range categories {
		_, err = s.store.GetCategoryByUUID(ctx, category.uuid)
		if err != nil {
			if !errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("failed to get category by name: %w", err)
			}
			err = s.store.CreateCategory(ctx, &models.Category{
				UUID:      category.uuid,
				Name:      category.name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	products := []struct {
		uuid        uuid.UUID
		t           enums.ProductType
		file        int
		from        int
		to          *int
		name        string
		description string
	}{
		{uuid: uuid.MustParse("671838e9-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeProduct, file: files[0].id, from: 3500, to: nil, name: "Фонтан из воздушных шаров", description: "Композиция по индивидуальному дизайну для любого события"},
		{uuid: uuid.MustParse("671853bf-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeProduct, file: files[1].id, from: 7000, to: nil, name: "Оформление помещения / фотозона", description: "Декорирование любого помещения по индивидуальному дизайну"},
		{uuid: uuid.MustParse("67185b7b-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeProduct, file: files[2].id, from: 5000, to: nil, name: "Коробка - сюрприз", description: "Подарочный бокс с композицией из шаров для любого события"},
		{uuid: uuid.MustParse("6bb9ab38-11b1-11f1-a6a2-c87f54a92045"), t: enums.ProductTypeProduct, file: files[3].id, from: 3000, to: nil, name: "Бабл бокс", description: "Креативная упаковка для небольшого подарка с шаром баблс"},
		{uuid: uuid.MustParse("671868fb-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[0].id, from: 3500, to: nil, name: "Фонтан из воздушных шаров", description: "Композиция по индивидуальному дизайну для любого события"},
		{uuid: uuid.MustParse("67186e98-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[1].id, from: 7000, to: nil, name: "Оформление помещения / фотозона", description: "Декорирование любого помещения по индивидуальному дизайну"},
		{uuid: uuid.MustParse("671874da-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[2].id, from: 5000, to: nil, name: "Коробка - сюрприз", description: "Подарочный бокс с композицией из шаров для любого события"},
		{uuid: uuid.MustParse("67186286-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[3].id, from: 3000, to: nil, name: "Бабл бокс", description: "Креативная упаковка для небольшого подарка с шаром баблс"},
	}
	for _, p := range products {
		_, err = s.store.GetProductByUUID(ctx, p.uuid)
		if err != nil {
			if !errors.Is(err, models.ErrNotFound) {
				return fmt.Errorf("failed to get product by id: %w", err)
			}
			err = s.store.CreateProduct(ctx, &models.Product{
				UUID:        p.uuid,
				Name:        p.name,
				Description: p.description,
				PriceFrom:   p.from,
				PriceTo:     p.to,
				Type:        p.t,
				FileID:      p.file,
				UserID:      user.ID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to create product: %w", err)
			}
		}
	}

	//reviews := []struct {
	//	uuid    uuid.UUID
	//	name    string
	//	content string
	//}{
	//	{uuid: uuid.MustParse("f3ab6d17-11af-11f1-b4af-c87f54a92045"), name: "Елена", content: "Огромное спасибо за шарики! Именинник был в восторге и доставка порадовала, все вовремя) Не пожалела, что обратилась именно к вам!"},
	//	{uuid: uuid.MustParse("f3ab78e7-11af-11f1-b4af-c87f54a92045"), name: "Евгений", content: "Обратился за бабл боксом с украшением внутри, хотелось поздравить девушку с годовщиной. Она очень удивилась такой креативной идее. Я тоже не видел ничего подобного в нашем городе до этого)) Желаю вам дальнейшего развития, вы классные!"},
	//	{uuid: uuid.MustParse("f3ab7f35-11af-11f1-b4af-c87f54a92045"), name: "Александр", content: "Выражаю благодарность студии за профессиональное оформление нашего корпоратива. Требовался лаконичный декор в цветах компании. Результат превзошел ожидания: композиции у входа и фотозона были выполнены безупречно, с вниманием к деталям. Отдельно отмечу пунктуальность, четкое соблюдение сроков и договоренностей."},
	//	{uuid: uuid.MustParse("f3ab84ad-11af-11f1-b4af-c87f54a92045"), name: "Любовь", content: "Благодарю персонал за то, что взялись за очень срочный заказ, выполнили и доставили максимально быстро. Посоветую вас друзьям и сама обращусь еще не раз."},
	//	{uuid: uuid.MustParse("f3ab8aa7-11af-11f1-b4af-c87f54a92045"), name: "Анна", content: "Как человек, который ценит визуал, долго искала в Екатеринбурге студию, которая умеет в тренды. Команда Ola предложила крутые цветовые сочетания для моей вечеринки. Шары держались несколько дней, не теряя вид! Ни одного лопнувшего. Это показатель. Если вам важен дизайн, атмосфера и стойкость - выбор очевиден."},
	//}
	//for i, r := range reviews {
	//	if err != nil {
	//		return fmt.Errorf("failed to create uuid: %w", err)
	//	}
	//	_, err = s.store.GetReviewByUUID(ctx, r.uuid)
	//	if err != nil {
	//		if !errors.Is(err, models.ErrNotFound) {
	//			return fmt.Errorf("failed to get review by id: %w", err)
	//		}
	//		err = s.store.CreateReview(ctx, &models.Review{
	//			UUID:        r.uuid,
	//			Name:        r.name,
	//			Content:     r.content,
	//			FileID:      files[i+4].id,
	//			UserID:      user.ID,
	//			PublishedAt: time.Now(),
	//			CreatedAt:   time.Now(),
	//			UpdatedAt:   time.Now(),
	//		})
	//		if err != nil {
	//			return fmt.Errorf("failed to create review: %w", err)
	//		}
	//	}
	//}

	s.store.Commit(ctx)
	s.log.Info("Seeder complete the work")
	return nil
}
