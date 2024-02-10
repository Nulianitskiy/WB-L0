package repository

import (
	"WB-L0/internal/model"
	"github.com/jmoiron/sqlx"
	"log"
)

func GetAllOrders(db *sqlx.DB) []model.Order {
	var orders []model.Order

	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery,
			&order.Payment,
			&order.Item,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			log.Printf("failed to scan orders: %v", err)
			return nil
		}

		orders = append(orders, order)
	}

	log.Printf("Orders wrote to memory")
	return orders
}

func GetOrdersByCustomer(db *sqlx.DB, customer string) []model.Order {
	var orders []model.Order

	rows, err := db.Query("SELECT * FROM orders WHERE customer_id=$1", customer)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery,
			&order.Payment,
			&order.Item,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			log.Printf("failed to scan orders: %v", err)
			return nil
		}

		orders = append(orders, order)
	}

	log.Printf("Orders wrote to memory")
	return orders
}

func PostDelivery(db *sqlx.DB, delivery model.Delivery, orderID string) error {
	_, err := db.Exec("INSERT INTO delivery (order_id, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		orderID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err != nil {
		return err
	}
	return nil
}

func PostPayment(db *sqlx.DB, payment model.Payment, orderID string) error {
	_, err := db.Exec("INSERT INTO payment (transaction_id, order_id, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		payment.TransactionID, orderID, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDT, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err != nil {
		return err
	}
	return nil
}

func PostItem(db *sqlx.DB, item model.Item, orderID string) error {
	_, err := db.Exec("INSERT INTO items (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
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
