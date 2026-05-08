package postgres_tx_manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager transaction yönetimini temsil eden interface'tir.
//
// Neden interface?
// Çünkü application/usecase katmanı somut pgx implementasyonuna değil,
// bu davranışa bağımlı olsun istiyoruz.
//
// Usecase şunu bilsin:
//
//   - Bir işi transaction içinde çalıştırabiliyorum.
//
// Usecase şunu bilmesin:
//
//   - pgxpool nedir?
//   - pgx.Tx nedir?
//   - Begin / Commit / Rollback nasıl çalışır?
type TxManager interface {
	WithInTx(ctx context.Context, fn func(ctx context.Context, tx DbTx) error) error
}

// PgxTxManager , TxManager interface'inin pgx ile yazılmış somut implementasyonudur.
//
// Bu struct infrastructure/shared tarafındadır.
// Çünkü pgx teknik detaydır.
// Domain veya usecase bunu bilmemelidir.
type PgxTxManager struct {
	pool *pgxpool.Pool
}

// NewPgxTxManager dışarıdan pgxpool.Pool alır.
//
// pgxpool.Pool uygulama başlarken main/wiring tarafında oluşturulur.
// TxManager ise bu pool üzerinden transaction başlatır.
func NewPgxTxManager(pool *pgxpool.Pool) TxManager {
	return &PgxTxManager{pool: pool}
}

// WithInTx WithinTx verilen fonksiyonu transaction içinde çalıştırır.
//
// Akış:
//
//  1. Transaction başlatılır.
//  2. fn çalıştırılır.
//  3. fn hata dönerse rollback yapılır.
//  4. fn hata dönmezse commit yapılır.
//  5. Commit hata verirse hata yukarı döndürülür.
//
// Buradaki en önemli mimari nokta:
//
//   - Transaction sınırını usecase belirler.
//   - Transaction'ın teknik yönetimini TxManager yapa
func (p PgxTxManager) WithInTx(ctx context.Context, fn func(ctx context.Context, tx DbExecutor) error) error {

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	// Burada fn'ye tx veriyoruz.
	//
	// fn içinde çağrılan repository'ler bu tx üzerinden çalışırsa
	// hepsi aynı transaction içinde çalışmış olur.

	//BURDA TAM OLARAK NE OLDUGUNU ANLAMADIM PARAMETRE OLARAK FONKSİYON ALDIK OKEY ? AMA FONKSİYONU PARAMETRELİ Mİ GONDERCEKLER NASIL GONDERCEKLER BURDA UYGULADIGIMIZ SYNTAXI ANLAMADIM BUNUDA COK İYİ OGRENMEK İSTİYORUM
	if err = fn(ctx, tx); err != nil {
		rollBackErr := tx.Rollback(ctx)
		if rollBackErr != nil {
			return errors.Join(
				fmt.Errorf("tx err : %w", err),
				fmt.Errorf("rollback err %w", rollBackErr),
			)
		}
		return fmt.Errorf("tx err : %w", err)

	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction : %w", err)
	}
	return nil
}
