CREATE DATABASE IF NOT EXISTS electroshop;

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    thumbnail VARCHAR(255),
    category VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(15, 2) NOT NULL,
    image_link VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO products (thumbnail, category, name, price, image_link) VALUES
('thumbnail1.jpg', 'Iphone', 'Iphone 13 Pro', 12000000, 'iphone13pro.jpg'),
('thumbnail2.jpg', 'Samsung', 'Samsung X flip', 20000000, 'samsungxflip.jpg'),
('thumbnail3.jpg', 'Xiaomi', 'Xiaomi Redmi Note 11 Pro', 3200000, 'xiaomiredminote11pro.jpg');

-- Membuat indeks untuk pencarian
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_name ON products(name);