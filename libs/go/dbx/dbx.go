package dbx

import (
	"database/sql"
	"log"
	"os"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func New() IORM {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return &orm{db: db}
}

func IsRecordNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

type IORM interface {
	List(dest any, conds ...any) error
	ListWithOrder(dest any, order string, conds ...any) error
	ListWithOrderAndLimit(dest any, order string, limit int, conds ...any) error
	ListIDs(table string, dest any, conds ...any) error
	ListForAppID(appId uint, dest any, limit *int) error
	ListForPlanID(appID, planID uint, dest any, preload string) error
	GetForPublicID(appID uint, publicID string, dest any) error
	FindForID(id uint, dest any) error
	FindFirstForAppID(appID uint, dest any) error
	FindFirst(dest any, conds ...any) error
	FindFirstOrNil(dest any, conds ...any) error
	Create(value any) error
	Save(dest any) error // Save is used to update an existing record
	Update(dest any, conds ...any) error
	UpdateForPublicID(appID uint, publicID string, dest any) error
	Upsert(dest any, query any, args ...any) error
	UpsertForPublicID(appID uint, publicID string, value any) error
	RemoveForPublicID(appID uint, publicID string, value any) error
	Delete(value, query any, args ...any) error
	DeleteForID(dest any) error
	Count(count *int64, dest any, query string, args ...any) error
	IncrementCount(dest any, field string, amount int64) error

	AutoMigrate(dst ...any) error
}

type orm struct {
	db *gorm.DB
}

func (o *orm) AutoMigrate(dst ...any) error {
	return o.db.AutoMigrate(
		&models.Account{},
		&models.App{},
		&models.Address{},
		&models.Company{},
		&models.Feature{},
		&models.Secret{},
		&models.User{},
		&models.Subscription{},
		&models.SubscriptionItem{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.Connection{},
		&models.Plan{},
		&models.PlanItem{},
		&models.PlanFeature{},
		&models.PlanGroup{},
		&models.Quote{},
		&models.CheckoutSession{},
		&models.PaymentMethod{},
		&models.Entitlement{},
		&models.MeterEvent{},
		&models.AppConfig{},
	)
}

func (o *orm) List(dest any, conds ...any) error {
	if err := o.db.Order("created_at DESC").Find(dest, conds...).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		log.Printf("Error listing records with conditions %v: %v", conds, err)

		return err
	}
	return nil
}

func (o *orm) ListWithOrder(dest any, order string, conds ...any) error {
	if err := o.db.Order(order).Find(dest, conds...).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		log.Printf("Error listing records with conditions %v: %v", conds, err)

		return err
	}
	return nil
}
func (o *orm) ListWithOrderAndLimit(dest any, order string, limit int, conds ...any) error {
	if err := o.db.Order(order).Find(dest, conds...).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		log.Printf("Error listing records with conditions %v: %v", conds, err)

		return err
	}
	return nil
}

func (o *orm) ListIDs(table string, dest any, conds ...any) error {
	log.Printf("Listing IDs from table %s with conditions: %v", table, conds)
	if err := o.db.
		Table(table).
		Where(conds[0], conds[1:]...).
		Select("id").
		Find(dest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No records found for conditions:", conds)
			return nil
		}

		log.Printf("Error listing ids with conditions %v: %v", conds, err)

		return err
	}
	return nil
}

func (o *orm) ListForAppID(appID uint, dest any, limit *int) error {
	tx := o.db.Where("app_id = ?", appID)

	if limit != nil {
		tx = tx.Limit(*limit)
	}

	if err := tx.Find(dest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		return err
	}
	return nil
}

func (o *orm) ListForPlanID(appID, planID uint, dest any, preload string) error {
	if err := o.db.
		Where("app_id = ? AND plan_id = ?", appID, planID).
		Preload(preload).
		Find(dest).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}

		log.Printf("Error listing records for appID %d and planID %d: %v", appID, planID, err)

		return err
	}
	return nil
}

func (o *orm) GetForPublicID(appID uint, publicID string, dest any) error {
	if err := o.db.Where("app_id = ? AND public_id = ?", appID, publicID).First(dest).Error; err != nil {
		log.Printf("Error finding record with appID %d and publicID %s: %v", appID, publicID, err)
		return err
	}
	return nil
}

func (o *orm) FindForID(id uint, dest any) error {
	if err := o.db.First(dest, "id = ?", id).Error; err != nil {
		log.Printf("Error finding record with ID %d: %v", id, err)
		return err
	}
	return nil
}

func (o *orm) FindFirst(dest any, conds ...any) error {
	if err := o.db.Where(conds[0], conds[1:]...).First(dest).Error; err != nil {
		log.Printf("Error finding first record with conditions %v: %v", conds, err)
		return err
	}
	return nil
}

func (o *orm) FindFirstOrNil(dest any, conds ...any) error {
	if err := o.db.Where(conds[0], conds[1:]...).First(dest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		log.Printf("Error finding first record with conditions %v: %v", conds, err)
		return err
	}
	return nil
}

func (o *orm) FindFirstForAppID(appID uint, dest any) error {
	if err := o.db.Where("app_id = ?", appID).First(dest).Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) Create(value any) error {
	if err := o.db.Create(value).Error; err != nil {
		log.Printf("Error creating record: %v", err)
		return err
	}
	return nil
}

func (o *orm) Save(dest any) error {
	if err := o.db.Save(dest).Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) Update(dest any, conds ...any) error {
	if err := o.db.
		Model(dest).
		Clauses(clause.Returning{}).
		Where(conds[0], conds[1:]...).
		Updates(dest).
		Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) UpdateForPublicID(appID uint, publicID string, dest any) error {
	if err := o.db.
		Model(dest).
		Clauses(clause.Returning{}).
		Where("app_id = ? AND public_id = ?", appID, publicID).
		Updates(dest).
		Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) Upsert(dest any, query any, args ...any) error {
	if err := o.db.First(dest, query, args).Error; err != nil {
		log.Printf("Error finding record for upsert with query %v: %v", query, err)

		if err != sql.ErrNoRows {
			if err := o.db.Create(dest).Error; err != nil {
				return err
			}
			log.Println("Record created successfully")
		} else {
			return err
		}
	}

	if err := o.db.Save(dest).Error; err != nil {
		return err
	}

	return nil
}

func (o *orm) UpsertForPublicID(appID uint, publicID string, value any) error {
	if err := o.db.
		Where("app_id = ?", appID).
		Where("public_id = ?", publicID).
		FirstOrCreate(value).
		Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) RemoveForPublicID(appID uint, publicID string, value any) error {
	if err := o.db.
		Where("app_id = ? AND public_id = ?", appID, publicID).
		Delete(value).Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) Delete(value, query any, args ...any) error {
	if err := o.db.
		Where(query, args...).
		Delete(value).Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) DeleteForID(dest any) error {
	if err := o.db.Unscoped().Delete(dest).Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) Count(count *int64, dest any, query string, args ...any) error {
	if err := o.db.Model(dest).Where(query, args...).Count(count).Error; err != nil {
		return err
	}
	return nil
}

func (o *orm) IncrementCount(dest any, field string, amount int64) error {
	if err := o.db.
		Model(dest).
		UpdateColumn(field, gorm.Expr(field+" + ?", amount)).
		Scan(&dest).Error; err != nil {
		return err
	}
	return nil
}
