package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/domain"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/bill/port"
	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/adapter/storage/types"
	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/fp"
)

type billRepo struct {
	db *sql.DB
}

func NewBillRepo(d *sql.DB) port.Repo {
	return &billRepo{db: d}
}

func (r *billRepo) Create(ctx context.Context, b *domain.Bill) (*domain.Bill, error) {
	query := `
	INSERT INTO bills(
		name,
		bill_type,
		bill_id,
		amount,
		due_date,
		image_id,
		apartment_id
	)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (bill_id) DO NOTHING
	RETURNING id;
	`

	args := []any{
		b.Name, b.Type.String(), b.BillNumber,
		b.Amount, b.DueDate, b.ImageID, b.ApartmentID,
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&b.ID)
	if err != nil {
		// If no ID was returned, it means bill_id already exists
		if err == sql.ErrNoRows {
			e := fmt.Errorf("bill with id %d already exists", b.BillNumber)
			return nil, fp.WrapErrors(bill.ErrAlreadyExists, e)
		}
		return nil, err
	}

	return b, nil
}

func (r *billRepo) Read(ctx context.Context, filter *domain.BillFilter) (*domain.Bill, error) {
	query := `SELECT 
		id, name, bill_type, bill_id, amount,  
		due_date, image_id, apartment_id 
		FROM bills WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIdx := 1

	if filter.ID != common.NilID {
		query += fmt.Sprintf(" AND id = $%d", argIdx)
		args = append(args, filter.ID)
		argIdx++
	}

	if filter.ApartmentID != common.NilID {
		query += fmt.Sprintf(" AND apartment_id = $%d", argIdx)
		args = append(args, filter.ApartmentID)
		argIdx++
	}

	if filter.Type != "" {
		query += fmt.Sprintf(" AND bill_type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}

	if filter.BillNumber != 0 {
		query += fmt.Sprintf(" AND bill_id = $%d", argIdx)
		args = append(args, filter.BillNumber)
		argIdx++
	}

	query += " LIMIT 1"

	row := r.db.QueryRowContext(ctx, query, args...)

	var b domain.Bill
	err := row.Scan(
		&b.ID,
		&b.Name,
		&b.Type,
		&b.BillNumber,
		&b.Amount,
		&b.DueDate,
		&b.ImageID,
		&b.ApartmentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, bill.ErrNotFound
		}
		return nil, err
	}

	return &b, nil
}

func (r *billRepo) GetUserBillShares(ctx context.Context, userID common.ID) ([]domain.UserBillShare, error) {
	query := `
        SELECT
            b.id,
            b.name,
            b.amount,
            COUNT(ua2.user_id) AS member_count,
            (b.amount / COUNT(ua2.user_id)) AS user_share,
            COALESCE(SUM(p.amount), 0) AS user_paid,
            (b.amount / COUNT(ua2.user_id)) - COALESCE(SUM(p.amount), 0) AS balance_due
        FROM users_apartments ua
        JOIN bills b ON b.apartment_id = ua.apartment_id
        JOIN users_apartments ua2 ON ua2.apartment_id = b.apartment_id AND ua2.created_at <= b.created_at
        LEFT JOIN payments p ON p.bill_id = b.id AND p.payer_id = ua.user_id
        WHERE ua.user_id = $1
        GROUP BY b.id, b.name, b.amount;
    `

	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []types.UserBillShare
	for rows.Next() {
		var s types.UserBillShare
		err := rows.Scan(
			&s.BillID,
			&s.BillName,
			&s.TotalAmount,
			&s.MemberCount,
			&s.SharePerUser,
			&s.UserPaid,
			&s.BalanceDue)
		if err != nil {
			return nil, err
		}
		s.UserID = userID.String()
		shares = append(shares, s)
	}
	return fp.Mapper(shares, func(ubs types.UserBillShare) domain.UserBillShare {
		return *types.UserBillShareStorageToDomain(&ubs)
	}), nil
}

func (r *billRepo) GetUserTotalDebt(ctx context.Context, userID common.ID) (int, error) {
	query := `
        SELECT SUM(user_share - COALESCE(user_paid, 0)) AS total_debt FROM (
            SELECT
                b.id,
                (b.amount / COUNT(ua2.user_id)) AS user_share,
                COALESCE(SUM(p.amount), 0) AS user_paid
            FROM users_apartments ua
            JOIN bills b ON b.apartment_id = ua.apartment_id
            JOIN users_apartments ua2 ON ua2.apartment_id = b.apartment_id AND ua2.created_at <= b.created_at
            LEFT JOIN payments p ON p.bill_id = b.id AND p.payer_id = ua.user_id
            WHERE ua.user_id = $1
            GROUP BY b.id, b.amount
        ) AS bill_balances;
    `

	var totalDebt sql.NullInt64
	err := r.db.QueryRow(query, userID.String()).Scan(&totalDebt)
	if err != nil {
		return 0, err
	}

	if !totalDebt.Valid {
		return 0, nil
	}
	return int(totalDebt.Int64), nil
}
