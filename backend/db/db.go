package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

func NewDB(logger *zap.SugaredLogger, host, user, password, dbname string, port int) (*sql.DB, error) {
	// psqlInfo := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
	// 	host, user, password, host, port, dbname)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	logger.Info(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	ctx, cnl_fnx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cnl_fnx()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
