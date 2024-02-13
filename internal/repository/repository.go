package repository

import (
	"WB-L0/internal/model"
	"github.com/jmoiron/sqlx"
)

func GetAllOrders(db *sqlx.DB) ([]model.Order, error) {
	var orders []model.Order

	if err := db.Select(&orders, "SELECT * FROM orders"); err != nil {
		return nil, err
	}

	// Для каждого заказа получаем связанные с ним товары, доставку и оплату
	for i := range orders {
		var items []model.Item
		if err := db.Select(&items, "SELECT * FROM items WHERE order_uid = $1", orders[i].OrderUID); err != nil {
			return nil, err
		}
		orders[i].Item = items

		var delivery model.Delivery
		if err := db.Get(&delivery, "SELECT * FROM delivery WHERE order_uid = $1", orders[i].OrderUID); err != nil {
			return nil, err
		}
		orders[i].Delivery = delivery

		var payment model.Payment
		if err := db.Get(&payment, "SELECT * FROM payment WHERE order_uid = $1", orders[i].OrderUID); err != nil {
			return nil, err
		}
		orders[i].Payment = payment
	}

	return orders, nil
}

func GetOrdersById(db *sqlx.DB, id string) ([]model.Order, error) {
	var orders []model.Order

	// Сначала получаем все заказы
	if err := db.Select(&orders, "SELECT * FROM orders WHERE order_uid=$1", id); err != nil {
		return nil, err
	}

	// Для каждого заказа получаем связанные с ним товары, доставку и оплату
	for i := range orders {
		var items []model.Item
		if err := db.Select(&items, "SELECT * FROM items WHERE order_uid = $1", orders[i].OrderUID); err != nil {
			return nil, err
		}
		orders[i].Item = items

		var delivery model.Delivery
		if err := db.Get(&delivery, "SELECT * FROM delivery WHERE order_uid = $1", orders[i].OrderUID); err != nil {
			return nil, err
		}
		orders[i].Delivery = delivery

		var payment model.Payment
		if err := db.Get(&payment, "SELECT * FROM payment WHERE order_uid = $1", orders[i].OrderUID); err != nil {
			return nil, err
		}
		orders[i].Payment = payment
	}

	return orders, nil
}

func PostDelivery(db *sqlx.DB, delivery model.Delivery, orderID string) error {
	_, err := db.Exec("INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		orderID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err != nil {
		return err
	}
	return nil
}

func PostPayment(db *sqlx.DB, payment model.Payment, orderID string) error {
	_, err := db.Exec("INSERT INTO payment (transaction_id, order_uid, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		payment.TransactionID, orderID, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDT, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err != nil {
		return err
	}
	return nil
}

func PostItem(db *sqlx.DB, item model.Item, orderID string) error {
	_, err := db.Exec("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		orderID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
	if err != nil {
		return err
	}
	return nil
}

func PostOrder(db *sqlx.DB, order model.Order) error {
	_, err := db.Exec("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}
	return nil
}
