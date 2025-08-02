-- Создание базы данных
CREATE DATABASE IF NOT EXISTS game_world;
USE game_world;

-- Таблица ObjectList
CREATE TABLE object_list (
                             id INT PRIMARY KEY,
                             name VARCHAR(255),
                             image VARCHAR(255),
                             description TEXT
);

-- Таблица Object
CREATE TABLE object (
                        id INT PRIMARY KEY,
                        object_list_id INT,
                        FOREIGN KEY (object_list_id) REFERENCES object_list(id) ON DELETE CASCADE
);

-- Таблица ItemList
CREATE TABLE item_list (
                           id INT PRIMARY KEY,
                           object_id INT,
                           rarity INT,
                           is_stackable BOOLEAN,
                           FOREIGN KEY (object_id) REFERENCES object(id) ON DELETE CASCADE
);

-- Таблица Item
CREATE TABLE item (
                      id INT PRIMARY KEY,
                      object_id INT,
                      item_list_id INT,
                      FOREIGN KEY (object_id) REFERENCES object(id) ON DELETE CASCADE,
                      FOREIGN KEY (item_list_id) REFERENCES item_list(id) ON DELETE CASCADE
);

-- Таблица EntityList
CREATE TABLE entity_list (
                             id INT PRIMARY KEY,
                             object_list_id INT,
                             damage DOUBLE,
                             speed DOUBLE,
                             cooldown DOUBLE,
                             damage_radius DOUBLE,
                             is_angry BOOLEAN,
                             visual_radius DOUBLE,
                             max_health DOUBLE,
                             model VARCHAR(255),
                             spawn VARCHAR(255),
                             is_open BOOLEAN,
                             is_spawning BOOLEAN,
                             is_pick_up BOOLEAN,
                             FOREIGN KEY (object_list_id) REFERENCES object_list(id) ON DELETE CASCADE
);

-- Таблица Entity
CREATE TABLE entity (
                        id INT PRIMARY KEY,
                        object_id INT,
                        entity_list_id INT,
                        health DOUBLE,
                        x DOUBLE,
                        y DOUBLE,
                        z DOUBLE,
                        FOREIGN KEY (object_id) REFERENCES object(id) ON DELETE CASCADE,
                        FOREIGN KEY (entity_list_id) REFERENCES entity_list(id) ON DELETE CASCADE
);

-- Таблица Inventory
CREATE TABLE inventory (
                           id VARCHAR(255) PRIMARY KEY,
                           entity_id VARCHAR(255),
                           item_id VARCHAR(255),
                           FOREIGN KEY (item_id) REFERENCES item(id) ON DELETE CASCADE
);

-- Таблица Player
CREATE TABLE player (
                        id INT PRIMARY KEY,
                        user_id INT,
                        x DOUBLE,
                        y DOUBLE,
                        z DOUBLE,
                        name VARCHAR(255)
);

-- Таблица World
CREATE TABLE world (
                       id VARCHAR(255) PRIMARY KEY,
                       x INT,
                       y INT,
                       value DOUBLE
);