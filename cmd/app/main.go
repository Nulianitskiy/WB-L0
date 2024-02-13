package main

import (
	"WB-L0/internal/cache"
	"WB-L0/internal/model"
	"WB-L0/internal/repository"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	libcache "github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"time"
)

func main() {

	// Подключение к базе данных

	connectionString := "user=dbuser password=tameimpala dbname=wbl0 host=localhost port=5436 sslmode=disable"
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	// Ping базы данных для проверки подключения
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки подключения к базе данных:", err)
	}
	log.Println("Подключение к базе данных PostgreSQL успешно")

	// Работа с кэшами

	// Создаем кэш с жизненным циклом объектов в 5 минут и очисткой устаревших объектов каждые 10 минут
	c := libcache.New(5*time.Minute, 10*time.Minute)

	orders, err := repository.GetAllOrders(db)
	if err != nil {
		log.Fatal("Ошибка получения данных из базы данных: ", err)
	}
	// Сохранение данных в кэше
	for _, order := range orders {
		c.Set(order.OrderUID, order, libcache.DefaultExpiration)
	}

	// Настройка веб-сервер

	router := gin.Default()

	router.LoadHTMLGlob("static/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})

	router.GET("/orders", func(ctx *gin.Context) {
		orders := mycache.GetCacheOrders(c)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, orders)
		}
	})

	router.GET("/orders/:customer", func(ctx *gin.Context) {
		customer := ctx.Param("customer")
		orders := mycache.GetCacheOrdersByCustomer(c, customer)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, orders)
		}
	})

	// Подключение к NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	// Осуществление подписки
	_, err = nc.Subscribe("foo", func(m *nats.Msg) {
		var order model.Order
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			log.Println("Ошибка при разборе JSON:", err)
		} else {

			// Добавление сообщение в кэш и бд
			c.Set(order.OrderUID, order, libcache.DefaultExpiration)

			err = repository.PostOrder(db, order)
			if err != nil {
				log.Printf("Ошибка при записи order в базу данных: %s\n", err)
			}
			err = repository.PostDelivery(db, order.Delivery, order.OrderUID)
			if err != nil {
				log.Printf("Ошибка при записи delivery в базу данных: %s\n", err)
			}
			err = repository.PostPayment(db, order.Payment, order.OrderUID)
			if err != nil {
				log.Printf("Ошибка при записи payment в базу данных: %s\n", err)
			}
			for _, item := range order.Item {
				err = repository.PostItem(db, item, order.OrderUID)
				if err != nil {
					log.Printf("Ошибка при записи item в базу данных: %s\n", err)
				}
			}

			log.Printf("Получено сообщение, uid: %s\n", order.OrderUID)
		}
	})
	if err != nil {
		log.Fatal("Ошибка подписки в NATS", err)
	}

	router.Run(":8080")
}
