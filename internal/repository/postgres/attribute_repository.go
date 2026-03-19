package postgres

import (
	"context"
	"eav-intentory/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttributeRepository struct {
	db *pgxpool.Pool
}

func NewAttributeRepository(pool *pgxpool.Pool) *AttributeRepository {
	return &AttributeRepository{db: pool}
}

func (a *AttributeRepository) Create(ctx context.Context, attribute *domain.Attribute) error {

	query := `insert into attributes (code,name,data_type) values ($1,$2,$3) returning id`

	err := a.db.QueryRow(ctx, query, attribute.Code, attribute.Name, attribute.DataType).Scan(&attribute.ID)
	if err != nil {
		return fmt.Errorf("attribute eklenirken hata olsutu %w", err)
	}
	return nil
}

func (a *AttributeRepository) GetById(ctx context.Context, id int) (*domain.Attribute, error) {
	query := `select id, name, code, data_type from attributes where id=$1`

	var attribute domain.Attribute // bunu pointer olarak almakla almamak ne değiştirir ? ve ek olarak attrbute := domain.attribute{} seklinde olusturmak ile farkı nedir ?
	// pointer olarak tanımlarsan nesneyi yaratmaz ve sadece pointer olusrturur ve bu nesneye atama yapmaya calısırken panic verir olmayana atama yapamazsın !
	// attribute := domain.Attribute{} olarak tanımlamanında bi farkı olmazdı sadece idiomatic oluor böyle
	err := a.db.QueryRow(ctx, query, id).Scan(&attribute.ID, &attribute.Name, &attribute.Code, &attribute.DataType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Nitelik bulunamadı id:%d", id)
		}

		return nil, fmt.Errorf("nitelik sorgularnırken olususan hata %w", err)
	}

	return &attribute, nil
}

func (a *AttributeRepository) GetAll(ctx context.Context, limit int, offset int) ([]domain.Attribute, int, error) {

	var totalCount int
	totalCountQuery := `select count(*) from attributes`
	a.db.QueryRow(ctx, totalCountQuery).Scan(&totalCount)

	query := `select id, name, code, data_type from attributes order by name asc limit $1 offset $2`

	rows, err := a.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("nitelikler sorgularnırken olususan hata %w", err)
	}

	defer rows.Close()
	var attributes []domain.Attribute // neden := ile tanımlanmıyor burdada kafamı karıstırdı bu konuyuda netleştirmek istiyorum
	// := ile acarsak 0 elemanlı dizi olsuturo ama varla olusturursak nil slice olusturuyo ve ilk append ile allocate edilior
	for rows.Next() {
		// attr := domain.Attribute{} // var ile kullanmanin farkları ne ? bu yanlıs mı ?
		// bunuda idiomatic olması için var ile tanımlamak gerekiyormus cünkü direk init etmiyoruz
		var attr domain.Attribute
		err = rows.Scan(&attr.ID, &attr.Name, &attr.Code, &attr.DataType)
		if err != nil {
			return nil, 0, fmt.Errorf("veritabani objesi bind edilirken olusan hata %w", err)
		}

		attributes = append(attributes, attr)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("nitelik listesi dongusunde hata: %w", err)
	}

	return attributes, totalCount, nil
}
