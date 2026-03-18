# Veritabanı bağlantı adresini bir değişkene atıyoruz ki her yerde tekrar yazmayalım
DB_URL=postgresql://postgres:123456@localhost:5432/eav_db?sslmode=disable

# Tabloları oluşturur (Up)
migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

# Tabloları siler (Down)
migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down


# Sistemin çökmesi durumunda migration versiyonunu zorla düzeltmek için
migrateforce:
	migrate -path migrations -database "$(DB_URL)" force 1


.PHONY: migrateup migratedown migrateforce