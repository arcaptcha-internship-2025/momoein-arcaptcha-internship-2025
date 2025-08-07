package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/internal/common"
	paymentd "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/domain"
	paymentp "github.com/arcaptcha-internship-2025/momoein-apartment/internal/payment/port"
)

type paymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) paymentp.Repo {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) CreatePayment(
	ctx context.Context, p *paymentd.Payment,
) (
	*paymentd.Payment, error,
) {
	query := `
		INSERT INTO payments (
			bill_id,
			payer_id,
			amount,
			paid_at,
			status,
			gateway,
			transaction_id,
			callback_data
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id
	`

	args := []any{
		p.BillID,
		p.PayerID,
		p.Amount,
		p.PaidAt,
		p.Status,
		p.Gateway,
		p.TransactionID,
		p.CallbackData,
	}

	var idStr string
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&idStr)
	if err != nil {
		return nil, err
	}

	p.ID = common.IDFromText(idStr)
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = p.CreatedAt
	return p, nil
}

func (r *paymentRepo) BatchCreatePayment(
	ctx context.Context, ps []*paymentd.Payment,
) (
	[]*paymentd.Payment, error,
) {
	if len(ps) == 0 {
		return []*paymentd.Payment{}, nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Rollback if not already committed
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO payments (
			bill_id,
			payer_id,
			amount,
			paid_at,
			status,
			gateway,
			transaction_id,
			callback_data
		) VALUES 
	`

	args := []any{}
	valueStrings := []string{}

	for i, p := range ps {
		start := i*8 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			start, start+1, start+2, start+3, start+4, start+5, start+6, start+7,
		))

		args = append(args,
			p.BillID,
			p.PayerID,
			p.Amount,
			p.PaidAt,
			p.Status,
			p.Gateway,
			p.TransactionID,
			p.CallbackData,
		)
	}

	query += strings.Join(valueStrings, ", ") + " RETURNING id"

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	now := time.Now().UTC()
	for i := 0; rows.Next(); i++ {
		var idStr string
		if err := rows.Scan(&idStr); err != nil {
			return nil, err
		}
		ps[i].ID = common.IDFromText(idStr)
		ps[i].CreatedAt = now
		ps[i].UpdatedAt = now
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ps, nil
}

func (r *paymentRepo) UpdateStatus(
	ctx context.Context,
	paymentIDs []common.ID,
	status paymentd.PaymentStatus,
) error {
	if len(paymentIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(paymentIDs))
	args := make([]any, 0, len(paymentIDs)+1)

	args = append(args, status)
	for i, id := range paymentIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2) // start from $2
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		UPDATE payments
		SET status = $1,
		    updated_at = NOW()
		WHERE id IN (%s)
	`, strings.Join(placeholders, ", "))

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *paymentRepo) UserBillBalanceDue(
	ctx context.Context, userID, billID common.ID,
) (
	int64, error,
) {
	query := `
        SELECT
            (b.amount / COUNT(ua2.user_id)) - COALESCE(SUM(p.amount), 0) AS balance_due
        FROM users_apartments ua
        JOIN bills b ON b.apartment_id = ua.apartment_id
        JOIN users_apartments ua2 ON ua2.apartment_id = b.apartment_id AND ua2.created_at <= b.created_at
        LEFT JOIN payments p ON p.bill_id = b.id AND p.payer_id = ua.user_id
        WHERE ua.user_id = $1 AND b.id = $2
        GROUP BY b.id, b.amount;
    `

	var balanceDue int64
	err := r.db.QueryRowContext(ctx, query, userID.String(), billID.String()).Scan(&balanceDue)
	if err != nil {
		return 0, err
	}

	return balanceDue, nil
}

func (r *paymentRepo) UserBillsBalanceDue(
	ctx context.Context, userID common.ID,
) (
	[]paymentd.BillWithAmount, error,
) {
	query := `
        SELECT
			b.id AS bill_id,
			ROUND((b.amount::numeric / COUNT(DISTINCT ua2.user_id)) - COALESCE(SUM(p.amount), 0), 2) AS balance_due
		FROM users_apartments ua
		JOIN bills b ON b.apartment_id = ua.apartment_id
		JOIN users_apartments ua2 
			ON ua2.apartment_id = b.apartment_id AND ua2.created_at <= b.created_at
		LEFT JOIN payments p 
			ON p.bill_id = b.id AND p.payer_id = ua.user_id AND p.deleted_at IS NULL
		WHERE ua.user_id = $1
		GROUP BY b.id, b.amount
		HAVING ROUND((b.amount::numeric / COUNT(DISTINCT ua2.user_id)) - COALESCE(SUM(p.amount), 0), 2) > 0
		ORDER BY b.created_at DESC;
    `

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := []paymentd.BillWithAmount{}
	for rows.Next() {
		var bwa paymentd.BillWithAmount
		err := rows.Scan(&bwa.BillID, &bwa.Amount)
		if err != nil {
			return nil, err
		}
		bills = append(bills, bwa)
	}

	return bills, nil
}
