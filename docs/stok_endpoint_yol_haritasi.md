# Stok Yönetim Sistemi — Endpoint ve Test Yol Haritası

Bu doküman, ürünler (`products`) ve depolar (`warehouses`) hazır olduktan sonra `stocks` ve `stock_movements` yapısını adım adım kurmak için hazırlanmıştır.

Ana mantık:

```text
products           => Ürün tanımları
warehouses         => Depo tanımları
stocks             => Ürün + depo bazlı mevcut stok bakiyesi
stock_movements    => Stok hareket geçmişi
```

Önemli karar:

```text
stocks tablosu doğrudan CRUD yapılmaz.
Stok miktarı, stock_movements üzerinden oluşur.
```

Yani başlangıçta şu endpointleri yazma:

```http
POST   /api/stocks
PUT    /api/stocks/{id}
DELETE /api/stocks/{id}
```

Bunun yerine stok miktarını değiştiren işlemleri `stock-movements` endpointleri üzerinden yap.

---

# 0. Başlamadan Önce Kontrol Et

Şu anda sistemde en az 1 ürün ve 1 depo olmalı.

## Kontrol SQL

```sql
SELECT * FROM products;
SELECT * FROM warehouses;
```

## Beklenen Durum

```text
En az 1 product kaydı var.
En az 1 warehouse kaydı var.
```

Örnek:

```text
Product:
id = 1
name = "Laptop"
sku = "LPT-001"

Warehouse:
id = 1
name = "Ana Depo"
code = "ANA"
```

---

# 1. İlk Endpoint: Stok Girişi

İlk yazman gereken endpoint budur.

## Endpoint

```http
POST /api/stock-movements/purchase-in
```

## Amaç

Bir ürünü bir depoya stok olarak eklemek.

## Request Body

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 100
}
```

## İş Akışı

```text
1. warehouse_id geçerli mi?
2. product_id geçerli mi?
3. quantity 0'dan büyük mü?
4. stocks tablosunda bu warehouse_id + product_id var mı?
5. Yoksa stocks satırı oluştur.
6. stock_movements tablosuna PURCHASE_IN hareketi ekle.
7. stocks.available_quantity değerini artır.
```

## Usecase Metodu

```go
PurchaseIn(ctx, input)
```

## Kullanılacak Repository Metotları

```go
stockRepo.CreateIfNotExists(ctx, tx, warehouseID, productID)
stockMovementRepo.Create(ctx, tx, movement)
stockRepo.IncreaseAvailable(ctx, tx, warehouseID, productID, quantity)
```

## Test Etmen Gerekenler

### Test 1 — Başarılı Stok Girişi

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 100
}
```

Beklenen sonuç:

```text
HTTP 200 veya 201
stocks tablosunda warehouse_id = 1 ve product_id = 1 için satır oluşmalı.
available_quantity = 100 olmalı.
stock_movements tablosunda PURCHASE_IN kaydı oluşmalı.
```

Kontrol SQL:

```sql
SELECT * 
FROM stocks 
WHERE warehouse_id = 1 
  AND product_id = 1;

SELECT * 
FROM stock_movements 
WHERE warehouse_id = 1 
  AND product_id = 1
ORDER BY id DESC;
```

### Test 2 — Aynı Ürüne Tekrar Giriş Yap

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 50
}
```

Beklenen sonuç:

```text
Yeni stocks satırı oluşmamalı.
Mevcut stocks.available_quantity 150 olmalı.
stock_movements tablosuna ikinci PURCHASE_IN hareketi eklenmeli.
```

Kontrol SQL:

```sql
SELECT *
FROM stocks
WHERE warehouse_id = 1
  AND product_id = 1;

SELECT COUNT(*)
FROM stock_movements
WHERE warehouse_id = 1
  AND product_id = 1
  AND movement_type = 'PURCHASE_IN';
```

### Test 3 — Geçersiz Quantity

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 0
}
```

Beklenen sonuç:

```text
HTTP 400
Stok değişmemeli.
Hareket kaydı oluşmamalı.
```

---

# 2. İkinci Endpoint: Stok Listeleme

Stok girişi çalıştıktan sonra mevcut stokları listele.

## Endpoint

```http
GET /api/stocks
```

## Amaç

Depolardaki ürünlerin mevcut stok bakiyelerini görmek.

## Örnek Response

```json
[
  {
    "stock_id": 1,
    "warehouse_id": 1,
    "warehouse_name": "Ana Depo",
    "product_id": 1,
    "product_name": "Laptop",
    "sku": "LPT-001",
    "available_quantity": 150,
    "reserved_quantity": 0
  }
]
```

## Kullanılacak SQL

```sql
SELECT 
    s.id AS stock_id,
    s.warehouse_id,
    w.name AS warehouse_name,
    s.product_id,
    p.name AS product_name,
    p.sku,
    s.available_quantity,
    s.reserved_quantity
FROM stocks s
JOIN warehouses w ON w.id = s.warehouse_id
JOIN products p ON p.id = s.product_id
ORDER BY s.id DESC;
```

## Test Etmen Gerekenler

### Test 1 — Liste Boş Değil Mi?

Beklenen sonuç:

```text
Daha önce purchase-in yaptığın ürün listede görünmeli.
available_quantity doğru olmalı.
```

### Test 2 — Warehouse Filtresi

Endpoint:

```http
GET /api/stocks?warehouse_id=1
```

Beklenen sonuç:

```text
Sadece warehouse_id = 1 olan stoklar dönmeli.
```

### Test 3 — Product Filtresi

Endpoint:

```http
GET /api/stocks?product_id=1
```

Beklenen sonuç:

```text
Sadece product_id = 1 olan stoklar dönmeli.
```

---

# 3. Üçüncü Endpoint: Stok Hareket Geçmişi Listeleme

## Endpoint

```http
GET /api/stock-movements
```

## Amaç

Stok giriş, çıkış, transfer ve düzeltme geçmişini listelemek.

## Örnek Response

```json
[
  {
    "id": 1,
    "warehouse_id": 1,
    "warehouse_name": "Ana Depo",
    "product_id": 1,
    "product_name": "Laptop",
    "sku": "LPT-001",
    "quantity": 100,
    "movement_type": "PURCHASE_IN"
  }
]
```

## Kullanılacak SQL

```sql
SELECT 
    sm.id,
    sm.warehouse_id,
    w.name AS warehouse_name,
    sm.product_id,
    p.name AS product_name,
    p.sku,
    sm.quantity,
    sm.movement_type
FROM stock_movements sm
JOIN warehouses w ON w.id = sm.warehouse_id
JOIN products p ON p.id = sm.product_id
ORDER BY sm.id DESC;
```

## Test Etmen Gerekenler

### Test 1 — Hareketler Listeleniyor Mu?

Beklenen sonuç:

```text
Daha önce yaptığın PURCHASE_IN kayıtları listelenmeli.
```

### Test 2 — Warehouse Filtresi

Endpoint:

```http
GET /api/stock-movements?warehouse_id=1
```

Beklenen sonuç:

```text
Sadece ilgili depoya ait hareketler gelmeli.
```

### Test 3 — Product Filtresi

Endpoint:

```http
GET /api/stock-movements?product_id=1
```

Beklenen sonuç:

```text
Sadece ilgili ürüne ait hareketler gelmeli.
```

### Test 4 — Movement Type Filtresi

Endpoint:

```http
GET /api/stock-movements?movement_type=PURCHASE_IN
```

Beklenen sonuç:

```text
Sadece PURCHASE_IN hareketleri gelmeli.
```

---

# 4. Dördüncü Endpoint: Stok Çıkışı

Bu aşamada hâlâ `orders` tablosu eklemek zorunda değilsin.

Buradaki işlem, basit manuel stok çıkışıdır.

## Endpoint

```http
POST /api/stock-movements/stock-out
```

veya DB hareket tipine daha yakın isim istersen:

```http
POST /api/stock-movements/order-out
```

## Amaç

Bir depodan ürün çıkışı yapmak.

## Request Body

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 20
}
```

## İş Akışı

```text
1. warehouse_id geçerli mi?
2. product_id geçerli mi?
3. quantity 0'dan büyük mü?
4. İlgili stock satırı var mı?
5. available_quantity yeterli mi?
6. stock_movements tablosuna ORDER_OUT hareketi ekle.
7. stocks.available_quantity değerini azalt.
```

## Usecase Metodu

```go
StockOut(ctx, input)
```

veya

```go
OrderOut(ctx, input)
```

## Kullanılacak Repository Metotları

```go
stockRepo.GetForUpdate(ctx, tx, warehouseID, productID)
stockMovementRepo.Create(ctx, tx, movement)
stockRepo.DecreaseAvailable(ctx, tx, warehouseID, productID, quantity)
```

## Test Etmen Gerekenler

### Test 1 — Başarılı Stok Çıkışı

Önce stokta 150 adet olduğunu varsayalım.

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 20
}
```

Beklenen sonuç:

```text
HTTP 200 veya 201
available_quantity 150'den 130'a düşmeli.
stock_movements tablosuna ORDER_OUT kaydı eklenmeli.
```

Kontrol SQL:

```sql
SELECT *
FROM stocks
WHERE warehouse_id = 1
  AND product_id = 1;

SELECT *
FROM stock_movements
WHERE warehouse_id = 1
  AND product_id = 1
ORDER BY id DESC;
```

### Test 2 — Yetersiz Stok

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 999999
}
```

Beklenen sonuç:

```text
HTTP 400
"insufficient stock" benzeri hata dönmeli.
stocks.available_quantity değişmemeli.
stock_movements tablosuna ORDER_OUT kaydı eklenmemeli.
```

### Test 3 — Quantity 0 veya Negatif

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 0
}
```

Beklenen sonuç:

```text
HTTP 400
Stok değişmemeli.
Hareket oluşmamalı.
```

---

# 5. Beşinci Endpoint: Depolar Arası Transfer

Bu aşamaya geçmeden önce sistemde en az 2 depo olmalı.

## Ön Kontrol SQL

```sql
SELECT * FROM warehouses;
```

## Endpoint

```http
POST /api/stock-movements/transfer
```

## Amaç

Bir ürünü bir depodan başka bir depoya aktarmak.

## Request Body

```json
{
  "from_warehouse_id": 1,
  "to_warehouse_id": 2,
  "product_id": 1,
  "quantity": 30
}
```

## İş Akışı

```text
1. from_warehouse_id geçerli mi?
2. to_warehouse_id geçerli mi?
3. product_id geçerli mi?
4. from_warehouse_id ve to_warehouse_id aynı mı?
5. quantity 0'dan büyük mü?
6. Kaynak depoda stok var mı?
7. Kaynak depoda yeterli stok var mı?
8. Hedef depoda bu ürün için stock satırı var mı?
9. Yoksa hedef depo için stock satırı oluştur.
10. Kaynak depoya TRANSFER_OUT hareketi yaz.
11. Hedef depoya TRANSFER_IN hareketi yaz.
12. Kaynak depodaki available_quantity değerini azalt.
13. Hedef depodaki available_quantity değerini artır.
```

## Usecase Metodu

```go
Transfer(ctx, input)
```

## Kullanılacak Repository Metotları

```go
stockRepo.GetForUpdate(ctx, tx, fromWarehouseID, productID)
stockRepo.CreateIfNotExists(ctx, tx, toWarehouseID, productID)
stockMovementRepo.Create(ctx, tx, transferOutMovement)
stockMovementRepo.Create(ctx, tx, transferInMovement)
stockRepo.DecreaseAvailable(ctx, tx, fromWarehouseID, productID, quantity)
stockRepo.IncreaseAvailable(ctx, tx, toWarehouseID, productID, quantity)
```

## Test Etmen Gerekenler

### Test 1 — Başarılı Transfer

Request:

```json
{
  "from_warehouse_id": 1,
  "to_warehouse_id": 2,
  "product_id": 1,
  "quantity": 30
}
```

Beklenen sonuç:

```text
HTTP 200 veya 201
warehouse_id = 1 olan stok 30 azalmalı.
warehouse_id = 2 olan stok 30 artmalı.
stock_movements tablosuna 2 kayıt eklenmeli:
- TRANSFER_OUT
- TRANSFER_IN
```

Kontrol SQL:

```sql
SELECT *
FROM stocks
WHERE product_id = 1
ORDER BY warehouse_id;

SELECT *
FROM stock_movements
WHERE product_id = 1
ORDER BY id DESC;
```

### Test 2 — Aynı Depoya Transfer Engellenmeli

Request:

```json
{
  "from_warehouse_id": 1,
  "to_warehouse_id": 1,
  "product_id": 1,
  "quantity": 30
}
```

Beklenen sonuç:

```text
HTTP 400
"source and target warehouse cannot be same" benzeri hata dönmeli.
Stok değişmemeli.
Hareket oluşmamalı.
```

### Test 3 — Yetersiz Stok

Request:

```json
{
  "from_warehouse_id": 1,
  "to_warehouse_id": 2,
  "product_id": 1,
  "quantity": 999999
}
```

Beklenen sonuç:

```text
HTTP 400
Stok değişmemeli.
TRANSFER_OUT ve TRANSFER_IN hareketleri oluşmamalı.
```

---

# 6. Altıncı Endpoint: Stok Düzeltme

Bu endpoint sayım farkı veya manuel düzeltme için kullanılır.

## Endpoint

```http
POST /api/stock-movements/adjustment
```

## Amaç

Mevcut stok miktarını sayım sonucuna göre düzeltmek.

## Request Body

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "new_quantity": 120,
  "reason": "Sayım düzeltmesi"
}
```

## İş Akışı

```text
1. warehouse_id geçerli mi?
2. product_id geçerli mi?
3. new_quantity negatif mi?
4. Mevcut stock satırını bul.
5. Mevcut miktar ile yeni miktar arasındaki farkı hesapla.
6. Fark 0 ise işlem yapma veya özel response dön.
7. Fark pozitifse stok artır.
8. Fark negatifse stok azalt.
9. stock_movements tablosuna ADJUSTMENT hareketi ekle.
```

## Örnek

```text
Mevcut stok: 100
Yeni sayım sonucu: 120
Fark: +20

Stok 20 artırılır.
ADJUSTMENT hareketi quantity = 20 olarak yazılır.
```

```text
Mevcut stok: 100
Yeni sayım sonucu: 85
Fark: -15

Stok 15 azaltılır.
ADJUSTMENT hareketi quantity = 15 olarak yazılır.
```

## Dikkat

Mevcut şemada `ADJUSTMENT` hareketinin artış mı azalış mı olduğunu net gösterecek alan yok.

Bu yüzden ileride şu iki seçenekten birini düşün:

```text
1. ADJUSTMENT_IN ve ADJUSTMENT_OUT hareket tipleri eklemek
2. stock_movements tablosuna direction alanı eklemek
```

Minimum sürümde şimdilik `ADJUSTMENT` ile devam edebilirsin.

## Test Etmen Gerekenler

### Test 1 — Stok Artıran Düzeltme

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "new_quantity": 200,
  "reason": "Sayımda fazla çıktı"
}
```

Beklenen sonuç:

```text
available_quantity 200 olmalı.
stock_movements içine ADJUSTMENT kaydı eklenmeli.
```

### Test 2 — Stok Azaltan Düzeltme

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "new_quantity": 80,
  "reason": "Sayımda eksik çıktı"
}
```

Beklenen sonuç:

```text
available_quantity 80 olmalı.
stock_movements içine ADJUSTMENT kaydı eklenmeli.
```

### Test 3 — Negatif New Quantity

Request:

```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "new_quantity": -10,
  "reason": "Hatalı test"
}
```

Beklenen sonuç:

```text
HTTP 400
Stok değişmemeli.
Hareket oluşmamalı.
```

---

# 7. Transaction Testleri

Stok sistemi transaction olmadan güvenilir olmaz.

Her stok değiştiren endpoint şu işlemleri tek transaction içinde yapmalı:

```text
1. stock_movements kaydı oluşturma
2. stocks bakiyesini güncelleme
```

Eğer biri başarısız olursa ikisi de geri alınmalı.

## Transaction Kullanılacak Endpointler

```http
POST /api/stock-movements/purchase-in
POST /api/stock-movements/stock-out
POST /api/stock-movements/transfer
POST /api/stock-movements/adjustment
```

## Test Fikri

Bilinçli olarak hareket kaydı oluştuktan sonra hata fırlat.

Beklenen sonuç:

```text
stock_movements kaydı kalmamalı.
stocks tablosu değişmemeli.
```

---

# 8. Concurrency Testleri

Aynı anda iki stok çıkışı geldiğinde stok eksiye düşmemeli.

## Riskli Senaryo

```text
Mevcut stok: 10

Aynı anda iki request gelir:
Request 1: 8 adet çıkış
Request 2: 8 adet çıkış
```

Yanlış sistemde sonuç:

```text
available_quantity = -6
```

Doğru sistemde:

```text
Bir request başarılı olur.
Diğeri yetersiz stok hatası alır.
available_quantity = 2
```

## Repository Tarafında Gerekli SQL

```sql
SELECT *
FROM stocks
WHERE warehouse_id = $1
  AND product_id = $2
FOR UPDATE;
```

veya güvenli update:

```sql
UPDATE stocks
SET available_quantity = available_quantity - $3
WHERE warehouse_id = $1
  AND product_id = $2
  AND available_quantity >= $3;
```

---

# 9. Önerilen Endpoint Sırası

Bu sırayı takip et.

## Aşama 1

```http
POST /api/stock-movements/purchase-in
```

Amaç:

```text
Stok girişi yapabilmek.
```

Bunu bitirmeden diğerlerine geçme.

---

## Aşama 2

```http
GET /api/stocks
```

Amaç:

```text
Mevcut stok bakiyesini görebilmek.
```

---

## Aşama 3

```http
GET /api/stock-movements
```

Amaç:

```text
Stok hareket geçmişini görebilmek.
```

---

## Aşama 4

```http
POST /api/stock-movements/stock-out
```

Amaç:

```text
Depodan manuel stok çıkışı yapabilmek.
```

---

## Aşama 5

```http
POST /api/stock-movements/transfer
```

Amaç:

```text
Depolar arası ürün transferi yapabilmek.
```

---

## Aşama 6

```http
POST /api/stock-movements/adjustment
```

Amaç:

```text
Sayım veya manuel düzeltme yapabilmek.
```

---

# 10. Bu Aşamada Eklememen Gerekenler

Şimdilik aşağıdakileri ekleme:

```text
orders
customers
invoices
shipments
reservations
reserved stock workflow
```

Bunlar sonraki modüller.

Şu an sadece stok motorunu kur.

Minimum stok motoru:

```text
Giriş yap.
Çıkış yap.
Transfer yap.
Düzeltme yap.
Listele.
Geçmişi görüntüle.
```

---

# 11. Minimum Başarı Kriteri

Bu aşamanın sonunda şunları yapabiliyor olmalısın:

```text
1. Bir depoya ürün girişi yapabiliyorum.
2. Stok bakiyesini görebiliyorum.
3. Stok hareket geçmişini görebiliyorum.
4. Depodan ürün çıkışı yapabiliyorum.
5. Yetersiz stokta hata alıyorum.
6. Depolar arası transfer yapabiliyorum.
7. Sayım düzeltmesi yapabiliyorum.
8. Hiçbir işlem stok eksiye düşürmüyor.
9. Hareket kaydı oluşmadan stok değişmiyor.
10. Stok değişmeden hareket kaydı kalmıyor.
```

---

# 12. Nihai Endpoint Listesi

```http
POST /api/stock-movements/purchase-in
POST /api/stock-movements/stock-out
POST /api/stock-movements/transfer
POST /api/stock-movements/adjustment

GET  /api/stocks
GET  /api/stock-movements
```

---

# 13. Sonraki Aşama

Bu yapı çalıştıktan sonra sisteme şunlar eklenebilir:

```text
1. created_at alanı
2. reason veya description alanı
3. reference_id alanı
4. ADJUSTMENT_IN / ADJUSTMENT_OUT ayrımı
5. reserved stock sistemi
6. order modülü
7. shipment/sevkiyat modülü
8. audit log
```

Özellikle `stock_movements` tablosuna ileride şu alanlar eklenmeli:

```sql
ALTER TABLE stock_movements
ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();

ALTER TABLE stock_movements
ADD COLUMN reason text NULL;

ALTER TABLE stock_movements
ADD COLUMN reference_id uuid NULL;
```

---

# Kısa Özet

Şu an yapacağın şey order sistemi değil.

Şu an yapacağın şey:

```text
stock movement oluşturmak
stock balance güncellemek
```

En önce bunu yaz:

```http
POST /api/stock-movements/purchase-in
```

Sonra bunu yaz:

```http
GET /api/stocks
```

Sonra bunu yaz:

```http
GET /api/stock-movements
```

Bunlar çalışınca stok çıkışı, transfer ve adjustment ekle.
