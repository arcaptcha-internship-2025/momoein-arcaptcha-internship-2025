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
		status,
		paid_at,
		due_date,
		image_id,
		apartment_id
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id;
	`
	args := []any{
		b.Name, b.Type.String(), b.BillNumber,
		b.Amount, b.Status, b.PaidAt,
		b.DueDate, b.ImageID, b.ApartmentID,
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&b.ID)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *billRepo) Read(ctx context.Context, filter *domain.BillFilter) (*domain.Bill, error) {
	query := `SELECT 
		id, name, bill_type, bill_id, amount, status, 
		paid_at, due_date, image_id, apartment_id 
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
		&b.Status,
		&b.PaidAt,
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

type UserBillShare struct {
	BillID       string
	BillName     string
	TotalAmount  int
	MemberCount  int
	SharePerUser int
	UserPaid     int
	BalanceDue   int
}

func (r *billRepo) GetUserBillShares(userID string) ([]UserBillShare, error) {
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

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []UserBillShare
	for rows.Next() {
		var s UserBillShare
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
		shares = append(shares, s)
	}
	return shares, nil
}

func (r *billRepo) GetUserTotalDebt(userID string) (int, error) {
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
	err := r.db.QueryRow(query, userID).Scan(&totalDebt)
	if err != nil {
		return 0, err
	}

	if !totalDebt.Valid {
		return 0, nil
	}
	return int(totalDebt.Int64), nil
}
