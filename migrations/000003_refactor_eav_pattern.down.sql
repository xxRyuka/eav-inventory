--  YENİ KURDUĞUMUZ KURUMSAL TABLOLARI SİLİYORUZ (Önce çocuklar, sonra baba)
DROP TABLE IF EXISTS product_attribute_values;
DROP TABLE IF EXISTS category_attributes;
DROP TABLE IF EXISTS attributes;

-- ESKİ BASİT EAV SİSTEMİNİ GERİ GETİRİYORUZ (Rolback)
CREATE TABLE IF NOT EXISTS category_attributes
(
    id          SERIAL PRIMARY KEY,
    category_id INTEGER      NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    data_type   VARCHAR(255) NOT NULL,
    is_required BOOLEAN      NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS product_attribute_values
(
    id                    SERIAL PRIMARY KEY,
    product_id            INTEGER NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    category_attribute_id INTEGER NOT NULL REFERENCES category_attributes (id) ON DELETE CASCADE,
    value                 TEXT    NOT NULL
);