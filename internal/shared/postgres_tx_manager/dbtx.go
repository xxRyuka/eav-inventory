package postgres_tx_manager

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DbExecutor  DBTX interface'i repository'lerin ihtiyaç duyduğu minimum veritabanı davranışını temsil eder.
//
// Neden var?
// Repository normalde SQL çalıştırmak ister.
// Ama SQL'i bazen doğrudan db pool ile, bazen de transaction ile çalıştırmamız gerekir.
//
// pgxpool.Pool da bu methodlara sahiptir.
// pgx.Tx de bu methodlara sahiptir.
//
// Bu sayede repository şu detayı bilmez:
//
//   - Ben transaction içinde miyim?
//   - Normal db bağlantısıyla mı çalışıyorum?
//
// Repository sadece şunu bilir:
//
//   - Bana verilen şey Exec / Query / QueryRow çalıştırabiliyor.
type DbExecutor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	// scan yok ?
}
