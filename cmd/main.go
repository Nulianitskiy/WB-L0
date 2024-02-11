package main

import (
	"WB-L0/internal/model"
	"WB-L0/internal/repository"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

func main() {
	// Replace with your actual database URL
	connectionString := "user=dbuser password=tameimpala dbname=wbl0 host=localhost port=5436 sslmode=disable"

	// Connect to the database
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return
	}

	fmt.Println("Connected to PostgreSQL database")

	router := gin.Default()

	router.LoadHTMLGlob("static/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	router.GET("/orders", func(c *gin.Context) {
		orders, err := repository.GetAllOrders(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, orders)
		}
	})
	router.GET("/orders/:customer", func(c *gin.Context) {
		customer := c.Param("customer")
		orders, err := repository.GetOrdersByCustomer(db, customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, orders)
		}
	})

	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	// Simple Async Subscriber
	_, err = nc.Subscribe("foo", func(m *nats.Msg) {
		var order model.Order
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			fmt.Println("Ошибка при разборе JSON:", err)
		}

		err = repository.PostOrder(db, order)
		if err != nil {
			fmt.Printf("Ошибка при записи order в бд: %s\n", err)
		}
		err = repository.PostDelivery(db, order.Delivery, order.OrderUID)
		if err != nil {
			fmt.Printf("Ошибка при записи delivery в бд: %s\n", err)
		}
		err = repository.PostPayment(db, order.Payment, order.OrderUID)
		if err != nil {
			fmt.Printf("Ошибка при записи payment в бд: %s\n", err)
		}
		for _, item := range order.Item {
			err = repository.PostItem(db, item, order.OrderUID)
			if err != nil {
				fmt.Printf("Ошибка при записи item в бд: %s\n", err)
			}
		}

		fmt.Printf("Received a message uid: %s\n", order.OrderUID)

	})
	if err != nil {
		log.Fatal(err)
	}

	router.Run(":8080")
}
