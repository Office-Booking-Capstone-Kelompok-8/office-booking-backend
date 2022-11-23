CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    role INT NOT NULL,
    is_verified bool NOT NULL
);

CREATE TABLE profile_pictures (
    id UUID PRIMARY KEY,
    url VARCHAR(255) NOT NULL,
    thumbnail_url VARCHAR(255) NOT NULL
);

CREATE TABLE user_details (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    profile_picture_id INT NOT NULL REFERENCES profile_pictures(id) ON DELETE CASCADE
);

CREATE TABLE payment_pictures (
    id UUID PRIMARY KEY,
    url INT,
    alt VARCHAR(255)
);

CREATE TABLE payments (
    id INT PRIMARY KEY,
    name INT NOT NULL,
    account_number VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    logo_id VARCHAR(255) NOT NULL REFERENCES payment_pictures(id) ON DELETE CASCADE
);

CREATE TABLE cities (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE districts (
    id INT PRIMARY KEY,
    city_id INT NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE buildings (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    capacity INT NOT NULL,
    annual_price INT NOT NULL,
    monthly_price INT NOT NULL,
    owner VARCHAR(255) NOT NULL,
    size INT NOT NULL,
    city_id INT NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    district_id INT NOT NULL REFERENCES districts(id) ON DELETE CASCADE,
    address VARCHAR(255) NOT NULL,
    logitude FLOAT NOT NULL,
    latitude FLOAT NOT NULL
);

CREATE TABLE pictures (
    id UUID PRIMARY KEY,
    building_id UUID NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    thumbnail_url VARCHAR(255) NOT NULL,
    alt VARCHAR(255) NOT NULL
);

CREATE TABLE facility_categories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL
);

CREATE TABLE facilities (
    id INT PRIMARY KEY,
    building_id VARCHAR(255) NOT NULL,
    category_id INT NOT NULL REFERENCES facility_categories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE status (
    id INT PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE reservations (
    id VARCHAR(255) PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    building_id VARCHAR(255) NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status_id INT NOT NULL REFERENCES status(id) ON DELETE CASCADE,
    message VARCHAR(255) NOT NULL
);