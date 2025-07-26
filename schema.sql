-- Таблица пользователей для регистрации и аутентификации
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Уникальный идентификатор пользователя
    username VARCHAR(255) UNIQUE NOT NULL,         -- Имя пользователя (должно быть уникальным)
    password_hash VARCHAR(255) NOT NULL,            -- Хеш пароля
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() -- Время создания записи
    );

-- Таблица для отслеживания местоположения игроков
CREATE TABLE IF NOT EXISTS player_locations (
                                                player_id UUID PRIMARY KEY,                     -- Идентификатор игрока (связан с user.id)
                                                x FLOAT NOT NULL,                               -- Координата X
                                                y FLOAT NOT NULL,                               -- Координата Y
                                                z FLOAT NOT NULL,                               -- Координата Z
                                                updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Время последнего обновления
    FOREIGN KEY (player_id) REFERENCES users(id) ON DELETE CASCADE -- Внешний ключ к таблице users
    );

-- Индекс для быстрого поиска по имени пользователя
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- Триггер для автоматического обновления updated_at в player_locations
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_player_locations_updated_at ON player_locations;
CREATE TRIGGER trg_player_locations_updated_at
    BEFORE UPDATE ON player_locations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
