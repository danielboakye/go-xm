package repo

import (
	"context"
	"database/sql"

	"github.com/danielboakye/go-xm/models"
)

type Repository struct {
	db *sql.DB
}

type IRepository interface {
	CreateCompany(context.Context, string, string, int64, bool, string) (string, error)
	UpdateCompany(context.Context, string, *string, *string, *int64, *bool, *string) error
	DeleteCompany(context.Context, string) error
	GetCompanyByID(context.Context, string) (company models.Company, err error)
	GetCompanyByName(context.Context, string) (company models.Company, err error)
}

func NewRepository(db *sql.DB) IRepository {
	return &Repository{db: db}
}
func (r *Repository) CreateCompany(
	ctx context.Context,
	name string,
	description string,
	amountOfEmployees int64,
	registered bool,
	companyType string,
) (companyID string, err error) {

	err = r.db.QueryRowContext(ctx, `
			INSERT INTO companies (company_name, description, amount_of_employees, is_registered, company_type)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING company_id
		`,
		name, description, amountOfEmployees, registered, companyType,
	).Scan(&companyID)
	return
}

func (r *Repository) UpdateCompany(
	ctx context.Context,
	companyID string,
	name *string,
	description *string,
	amountOfEmployees *int64,
	registered *bool,
	companyType *string,
) (err error) {

	_, err = r.db.ExecContext(ctx, `
			UPDATE companies 
			SET 
				company_name = coalesce($2, company_name), 
				description = coalesce($3, description), 
				amount_of_employees = coalesce($4, amount_of_employees), 
				is_registered = coalesce($5, is_registered), 
				company_type = coalesce($6, company_type),
				modified_at = now()
			WHERE
				company_id = $1
		`,
		companyID, name, description, amountOfEmployees, registered, companyType,
	)
	return
}

func (r *Repository) DeleteCompany(
	ctx context.Context,
	companyID string,
) (err error) {

	_, err = r.db.ExecContext(ctx, `
			UPDATE companies 
			SET 
				deleted_at = now()
			WHERE
				company_id = $1
		`,
		companyID,
	)
	return
}

func (r *Repository) GetCompanyByID(
	ctx context.Context,
	companyID string,
) (company models.Company, err error) {

	err = r.db.QueryRowContext(ctx, `
			SELECT
				company_id, company_name, description, amount_of_employees, is_registered, company_type
			FROM companies
			WHERE company_id = $1
				AND deleted_at IS NULL
		`,
		companyID,
	).Scan(
		&company.ID,
		&company.Name,
		&company.Description,
		&company.AmountOfEmployees,
		&company.Registered,
		&company.CompanyType,
	)
	return
}

func (r *Repository) GetCompanyByName(
	ctx context.Context,
	companyName string,
) (company models.Company, err error) {

	err = r.db.QueryRowContext(ctx, `
			SELECT
				company_id, company_name, description, amount_of_employees, is_registered, company_type
			FROM companies
			WHERE company_name = $1
				AND deleted_at IS NULL
		`,
		companyName,
	).Scan(
		&company.ID,
		&company.Name,
		&company.Description,
		&company.AmountOfEmployees,
		&company.Registered,
		&company.CompanyType,
	)
	return
}
