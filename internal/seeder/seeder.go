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
		content string
		uuid    uuid.UUID
	}{
		{uuid: uuid.MustParse("f3aaf3df-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab0dd5-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab14fe-11af-11f1-b4af-c87f54a92045")},
		{uuid: uuid.MustParse("f3ab1996-11af-11f1-b4af-c87f54a92045")},
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
		files[i].content = file.Content
	}

	categories := []struct {
		uuid uuid.UUID
		name string
		id   int
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
		item, getErr := s.store.GetCategoryByUUID(ctx, category.uuid)
		if getErr != nil {
			if !errors.Is(getErr, models.ErrNotFound) {
				return fmt.Errorf("failed to get category by name: %w", getErr)
			}
			createErr := s.store.CreateCategory(ctx, &models.Category{
				UUID:      category.uuid,
				Name:      category.name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
			if createErr != nil {
				return fmt.Errorf("failed to create category: %w", createErr)
			}

			item, getErr = s.store.GetCategoryByUUID(ctx, category.uuid)
			if getErr != nil {
				return fmt.Errorf("failed to reload category: %w", getErr)
			}
		}
		for i := range categories {
			if categories[i].uuid == category.uuid {
				categories[i].id = item.ID
				break
			}
		}
	}

	products := []struct {
		uuid        uuid.UUID
		t           enums.ProductType
		file        string
		from        int
		to          *int
		name        string
		description string
		categories  []int
	}{
		{uuid: uuid.MustParse("671838e9-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeProduct, file: files[0].content, from: 3500, to: nil, name: "Фонтан из воздушных шаров", description: "Композиция по индивидуальному дизайну для любого события", categories: []int{categories[0].id, categories[2].id, categories[6].id, categories[7].id, categories[8].id}},
		{uuid: uuid.MustParse("671853bf-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeProduct, file: files[1].content, from: 7000, to: nil, name: "Оформление помещения / фотозона", description: "Декорирование любого помещения по индивидуальному дизайну", categories: []int{categories[1].id, categories[2].id, categories[3].id, categories[8].id}},
		{uuid: uuid.MustParse("67185b7b-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeProduct, file: files[2].content, from: 5000, to: nil, name: "Коробка - сюрприз", description: "Подарочный бокс с композицией из шаров для любого события", categories: []int{categories[2].id, categories[4].id, categories[5].id}},
		{uuid: uuid.MustParse("6bb9ab38-11b1-11f1-a6a2-c87f54a92045"), t: enums.ProductTypeProduct, file: files[3].content, from: 3000, to: nil, name: "Бабл бокс", description: "Креативная упаковка для небольшого подарка с шаром баблс", categories: []int{categories[2].id, categories[3].id, categories[5].id, categories[8].id}},
		{uuid: uuid.MustParse("671868fb-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[0].content, from: 3500, to: nil, name: "Фонтан из воздушных шаров", description: "Композиция по индивидуальному дизайну для любого события", categories: []int{categories[0].id, categories[2].id, categories[6].id, categories[7].id, categories[8].id}},
		{uuid: uuid.MustParse("67186e98-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[1].content, from: 7000, to: nil, name: "Оформление помещения / фотозона", description: "Декорирование любого помещения по индивидуальному дизайну", categories: []int{categories[1].id, categories[2].id, categories[3].id, categories[8].id}},
		{uuid: uuid.MustParse("671874da-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[2].content, from: 5000, to: nil, name: "Коробка - сюрприз", description: "Подарочный бокс с композицией из шаров для любого события", categories: []int{categories[2].id, categories[4].id, categories[5].id}},
		{uuid: uuid.MustParse("67186286-11b1-11f1-afae-c87f54a92045"), t: enums.ProductTypeService, file: files[3].content, from: 3000, to: nil, name: "Бабл бокс", description: "Креативная упаковка для небольшого подарка с шаром баблс", categories: []int{categories[2].id, categories[3].id, categories[5].id, categories[8].id}},
	}
	for _, p := range products {
		product, getErr := s.store.GetProductByUUID(ctx, p.uuid)
		if getErr != nil {
			if !errors.Is(getErr, models.ErrNotFound) {
				return fmt.Errorf("failed to get product by id: %w", getErr)
			}
			createErr := s.store.CreateProduct(ctx, &models.Product{
				UUID:        p.uuid,
				Name:        p.name,
				Description: p.description,
				PriceFrom:   p.from,
				PriceTo:     p.to,
				Type:        p.t,
				IsMain:      true,
				FileContent: p.file,
				UserID:      user.ID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
			if createErr != nil {
				return fmt.Errorf("failed to create product: %w", createErr)
			}

			product, getErr = s.store.GetProductByUUID(ctx, p.uuid)
			if getErr != nil {
				return fmt.Errorf("failed to reload product: %w", getErr)
			}
		}

		err = s.store.ReplaceProductCategories(ctx, product.ID, p.categories)
		if err != nil {
			return fmt.Errorf("failed to create product categories: %w", err)
		}
	}

	s.store.Commit(ctx)
	s.log.Info("Seeder complete the work")
	return nil
}
