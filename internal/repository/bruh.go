package repository

import (
	"context"
	"github.com/GameXost/wbTestCase/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryDelivery = `
		SELECT
		d.id, d.name,
		d.phone, d.zip,
		d.city, d.address,
		d.region, d.email
		FROM delivery AS d
		WHERE d.order_uid = $1;
		`

	queryBaseOrder = `
		SELECT
		o.track_number, o.entry,
		o.locale, o.internal_signature,
		o.customer_id, o.delivery_service,
		o.shardkey, o.sm_id,
		o.date_created, o.oof_shard
		FROM orders AS o
		WHERE o.order_uid = $1;
		`

	queryPayment = `
			SELECT
			p.transaction, p.request_id,
			p.currency, p.provider,
			p.amount, p.payment_dt,
			p.bank, p.delivery_cost,
			p.goods_total, p.custom_fee
			FROM payment AS p
			WHERE p.order_id = $1
			`

	queryItems = `
			SELECT
			i.chrt_id, i.track_number,
			i.price, i.rid,
			i.name, i.sale,
			i.size, i.total_price,
			i.nm_id, i.brand,
			i.status
			FROM items AS i
			WHERE i.order_uid = $1;
			`
)
const (
	queryInsertDelivery = `
						INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
						`
	queryInsertPayment = `
						INSERT INTO payment (order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
						`
	queryInsertItem = `
						INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
						`
	queryInsertOrder = `
						INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`
)

type dbExecutor interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type Repo struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool, tx: nil}
}
func (r *Repo) RepoWithTX(tx pgx.Tx) *Repo {
	return &Repo{pool: r.pool, tx: tx}
}

func (r *Repo) executor() dbExecutor {
	if r.tx == nil {
		return r.pool
	}
	return r.tx
}

func (r *Repo) GetBaseOrderOnId(ctx context.Context, OrderUId string) (*models.Order, error) {
	var order models.Order
	order.OrderUId = OrderUId
	err := r.executor().QueryRow(ctx, queryBaseOrder, OrderUId).Scan(
		&order.TrackNumber, &order.Entry,
		&order.Locale, &order.InternalSignature,
		&order.CustomerId, &order.DeliveryService,
		&order.Shardkey, &order.SmId,
		&order.DateCreated, &order.OofShard,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repo) GetDeliveryOnID(ctx context.Context, OrderUId string) (*models.Delivery, error) {
	var delivery models.Delivery
	delivery.OrderUId = OrderUId

	err := r.executor().QueryRow(ctx, queryDelivery, OrderUId).Scan(
		&delivery.Id, &delivery.Name,
		&delivery.Phone, &delivery.Zip,
		&delivery.City, &delivery.Address,
		&delivery.Region, &delivery.Email,
	)
	if err != nil {
		return nil, err
	}
	return &delivery, nil

}

func (r *Repo) GetPaymentOnID(ctx context.Context, OrderUId string) (*models.Payment, error) {
	var payment models.Payment

	err := r.executor().QueryRow(ctx, queryPayment, OrderUId).Scan(
		&payment.Transaction, &payment.RequestId,
		&payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentDt,
		&payment.Bank, &payment.DeliveryCost,
		&payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *Repo) GetItemsOnID(ctx context.Context, OrderUId string) ([]models.Item, error) {
	var items []models.Item

	rows, err := r.executor().Query(ctx, queryItems, OrderUId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Item
		err = rows.Scan(
			&item.ChrtId, &item.TrackNumber,
			&item.Price, &item.RID,
			&item.Name, &item.Sale,
			&item.Size, &item.TotalPrice,
			&item.NmId, &item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return items, nil
}

func (r *Repo) GetFullOrderOnId(ctx context.Context, OrderUId string) (*models.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	txRepo := r.RepoWithTX(tx)

	order, err := txRepo.GetBaseOrderOnId(ctx, OrderUId)
	if err != nil {
		return nil, err
	}

	items, err := txRepo.GetItemsOnID(ctx, OrderUId)
	if err != nil {
		return nil, err
	}

	delivery, err := txRepo.GetDeliveryOnID(ctx, OrderUId)
	if err != nil {
		return nil, err
	}

	payment, err := txRepo.GetPaymentOnID(ctx, OrderUId)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment
	order.Delivery = *delivery
	order.Items = items

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return order, nil

}

func (r *Repo) CreateDelivery(ctx context.Context, delivery *models.Delivery) error {
	_, err := r.executor().Exec(ctx, queryInsertDelivery, delivery.OrderUId, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CreatePayment(ctx context.Context, payment *models.Payment) error {
	_, err := r.executor().Exec(ctx, queryInsertPayment, payment.OrderId, payment.Transaction, payment.RequestId, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CreateItem(ctx context.Context, item *models.Item) error {
	_, err := r.executor().Exec(ctx, queryInsertItem, item.OrderUId, item.ChrtId, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CreateBaseOrder(ctx context.Context, order *models.Order) error {
	_, err := r.executor().Exec(ctx, queryInsertOrder, order.OrderUId, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CreateFullOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txRepo := r.RepoWithTX(tx)

	err = txRepo.CreateBaseOrder(ctx, order)
	if err != nil {
		return err
	}

	err = txRepo.CreatePayment(ctx, &order.Payment)
	if err != nil {
		return err
	}

	err = txRepo.CreateDelivery(ctx, &order.Delivery)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		err = txRepo.CreateItem(ctx, &item)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
