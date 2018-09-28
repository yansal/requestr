package app

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yansal/requestr/internal/broker"
	"github.com/yansal/requestr/internal/mux"
	"github.com/yansal/requestr/internal/requestr"
	"github.com/yansal/requestr/internal/server"
	"github.com/yansal/requestr/internal/store"
	"github.com/yansal/requestr/internal/tasks"
	"github.com/yansal/requestr/internal/worker"
)

type App struct {
	port string
	mux  *http.ServeMux

	receiver broker.Receiver
	tasks    tasks.Tasks
}

func New() *App {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	app := new(App)

	pgURL := "host=/var/run/postgresql sslmode=disable"
	db := sqlx.MustConnect("postgres", pgURL)
	store := store.New(db)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://:6379"
	}
	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	poolsize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	redisOpts.PoolSize = poolsize
	redis := redis.NewClient(redisOpts)
	if err := redis.Ping().Err(); err != nil {
		panic(err)
	}
	publisher := broker.NewRedisPublisher(redis)
	template := template.Must(template.New("").ParseGlob("templates/*.html"))

	app.port = os.Getenv("PORT")
	if app.port == "" {
		app.port = "8080"
	}
	app.mux = mux.New(store, publisher, template)

	app.receiver = broker.NewRedisReceiver(redis)
	pool, err := requestr.NewPool(ctx, store, app.receiver, tasks.DebugHandler, publisher)
	if err != nil {
		panic(err)
	}
	app.tasks = tasks.New(pool)

	return app
}

func (app *App) StartServer(ctx context.Context) func() error {
	return server.Start(ctx, app.port, app.mux)
}

func (app *App) StartWorker(ctx context.Context) func() error {
	return worker.Start(ctx, app.receiver, app.tasks)
}
