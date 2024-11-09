package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func main() {
	exampleGoFakeItSQL()
}
func exampleGoFakeItSQL() {
	// так-то это должно храниться в переменных окружениях, но для вашего удобства это тут
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	gofakeit.Seed(11)
	res, err := gofakeit.SQL(&gofakeit.SQLOptions{
		Table: "orders",
		Count: 100000,
		Fields: []gofakeit.Field{
			{Name: "id", Function: "autoincrement"},
			{Name: "user_id", Function: "number", Params: gofakeit.MapParams{"min": {"1"}, "max": {"10"}}},
			{Name: "created_at", Function: "date"},
			{Name: "expiry_date", Function: "date", Params: gofakeit.MapParams{"min": {"now"}, "max": {"30d"}}},
			{Name: "deadline_return_order", Function: "date", Params: gofakeit.MapParams{"min": {"now"}, "max": {"30d"}}},
			{Name: "is_issued", Function: "bool"},
			{Name: "is_returned_from_user", Function: "bool"},
			{Name: "packaging", Function: "randomstring", Params: gofakeit.MapParams{"strs": {"'film'", "'box'", "'bag'"}}},
			{Name: "weigh", Function: "float32", Params: gofakeit.MapParams{"min": {"1.0"}, "max": {"50.0"}}},
			{Name: "cost", Function: "float32", Params: gofakeit.MapParams{"min": {"10.0"}, "max": {"500.0"}}},
		},
	})
	// это сделано для нас, поэтому можно фаталить с ошибкой, т.к. пользователь не увидит)))
	if err != nil {
		log.Fatal(err)
	}
	results, err := pool.Exec(ctx, res)
	if err != nil || results.RowsAffected() == 0 {
		log.Fatal(err)
	}
	fmt.Println("fake orders have been added")
}
