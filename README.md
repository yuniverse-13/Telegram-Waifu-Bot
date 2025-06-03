# Waifu Telegram Bot (Телеграм-бот "Вайфу")

Телеграм-бот на Go для поиска и отображения информации о персонажах (вайфу) из вашей базы данных PostgreSQL.

## ✨ Возможности

*   Поиск персонажей по ID.
*   Поиск персонажей по имени или альтернативному имени (регистронезависимый для основного имени).
*   Отображение случайного персонажа.
*   Отображение информации о персонаже: имя, ID, тайтл, описание, рейтинг и изображение.
*   Использование PostgreSQL для хранения данных.
*   Использование GORM как ORM для взаимодействия с базой данных.
*   Управление миграциями схемы БД с помощью `golang-migrate`.
*   Возможность запуска в Docker-контейнере с помощью Docker Compose.
*   Форматирование сообщений с помощью MarkdownV2 для лучшей читаемости.

## 📋 Требования

*   Go (версия 1.24 или выше рекомендуется)
*   PostgreSQL (версия 17 или выше рекомендуется)
*   Docker и Docker Compose (рекомендуется для простоты развертывания)
*   `golang-migrate` CLI (если вы планируете управлять миграциями вручную вне Docker)

## 🚀 Установка и Запуск

### С помощью Docker Compose (Рекомендуемый способ)

Это самый простой способ запустить бота вместе с базой данных.

1.  **Клонируйте репозиторий:**
    ```bash
    git clone https://github.com/YourUsername/telegram-waifu-bot.git
    cd telegram-waifu-bot
    ```
    *(Замени `YourUsername/telegram-waifu-bot` на URL твоего репозитория)*

2.  **Настройте переменные окружения:**
    Создайте файл `.env` в корне проекта, скопировав `.env.example` (если он у вас есть) или создав новый:
    ```env
    # Токен вашего Telegram бота
    TELEGRAM_BOT_TOKEN=ВАШ_ТЕЛЕГРАМ_БОТ_ТОКЕН

    # Настройки для подключения к PostgreSQL из docker-compose
    # Эти переменные используются сервисом postgres в docker-compose.yml
    POSTGRES_USER=waifu_bot_user
    POSTGRES_PASSWORD=supersecretpassword
    POSTGRES_DB=waifu_bot_db
    DB_PORT_HOST=5432 # Порт, который будет проброшен на хост-машину для PostgreSQL

    # DATABASE_URL для самого бота (он будет обращаться к БД внутри Docker-сети)
    # Имя хоста 'postgres_db' соответствует имени сервиса БД в docker-compose.yml
    DATABASE_URL=postgres://waifu_bot_user:supersecretpassword@postgres_db:5432/waifu_bot_db?sslmode=disable

    # Переменные для сервиса бота в docker-compose.yml (дублируют некоторые выше)
    DB_USER=${POSTGRES_USER}
    DB_PASSWORD=${POSTGRES_PASSWORD}
    DB_NAME=${POSTGRES_DB}
    DB_HOST=postgres_db # Имя сервиса БД в Docker-сети
    DB_PORT_CONTAINER=5432 # Порт БД внутри контейнера
    ```
    **Важно:** Замените `ВАШ_ТЕЛЕГРАМ_БОТ_ТОКЕН` на реальный токен вашего бота, полученный от [@BotFather](https://t.me/BotFather). Пароль и другие данные также рекомендуется сменить на более надежные.

3.  **Запустите сервисы:**
    ```bash
    docker-compose up --build
    ```
    При первом запуске миграции базы данных будут применены автоматически (согласно `entrypoint.sh`).

4.  **Остановка:**
    ```bash
    docker-compose down
    ```

### Локальный запуск (без Docker)

1.  **Клонируйте репозиторий:** (см. выше)

2.  **Настройте PostgreSQL:**
    Убедитесь, что у вас запущен и настроен сервер PostgreSQL. Создайте базу данных и пользователя для бота.

3.  **Настройте переменные окружения:**
    Создайте файл `.env` в корне проекта:
    ```env
    TELEGRAM_BOT_TOKEN=ВАШ_ТЕЛЕГRAM_БОТ_ТОКЕН
    DATABASE_URL=postgres://ВАШ_ПОЛЬЗОВАТЕЛЬ_БД:ВАШ_ПАРОЛЬ_БД@localhost:5432/ВАША_БАЗА_ДАННЫХ?sslmode=disable
    ```
    Замените значения на ваши.

4.  **Установите зависимости:**
    ```bash
    go mod tidy
    ```

5.  **Примените миграции базы данных:**
    Убедитесь, что у вас установлен `golang-migrate`. [Инструкции по установке](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).
    ```bash
    migrate -path ./migrations -database "${DATABASE_URL}" up
    ```
    (Убедитесь, что переменная `DATABASE_URL` доступна в вашем шелле или укажите URL напрямую)

6.  **Запустите бота:**
    ```bash
    go run cmd/bot/main.go
    ```

## ⚙️ Конфигурация

Бот конфигурируется с помощью переменных окружения:

*   `TELEGRAM_BOT_TOKEN`: **Обязательно.** Токен вашего Telegram бота.
*   `DATABASE_URL`: **Обязательно.** Строка подключения к PostgreSQL.
    *   Формат: `postgres://[user]:[password]@[host]:[port]/[dbname]?sslmode=[disable|require|etc.]`
*   (Для Docker Compose также используются `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` для инициализации сервиса БД).

## 🤖 Использование (Команды бота)

*   `/start` - Приветственное сообщение.
*   `/info` - Показывает список доступных команд.
*   `/character <ID или Имя>` - Поиск персонажа по его ID или имени/альтернативному имени.
    *   Пример: `/character 1`
    *   Пример: `/character Мегумин`
*   `/randomcharacter` - Показывает случайного персонажа из базы данных.

## 📁 Структура Проекта
.
├── .env # Файл с переменными окружения (локальный, не коммитится)
├── .env.example # Файл с примером твоего файла .env
├── .git # Директория Git
├── .gitignore # Файл для Git, указывающий игнорируемые файлы
├── Dockerfile # Dockerfile для сборки образа бота
├── README.md # Этот файл
├── cmd/
│ └── bot/
│ └── main.go # Точка входа приложения, инициализация
├── docker-compose.yml # Docker Compose файл для оркестрации сервисов (бот и БД)
├── entrypoint.sh # Скрипт точки входа для Docker-контейнера бота (применяет миграции)
├── go.mod # Файл управления зависимостями Go
├── go.sum # Контрольные суммы зависимостей Go
├── internal/
│ ├── bot/
│ │ ├── bot.go # Основная логика Telegram бота, инициализация API, цикл обновлений
│ │ └── handlers/ # Обработчики конкретных команд Telegram
│ │ ├── character_handler.go
│ │ ├── character_response.go # Формирование ответов для команды character
│ │ ├── random_character_handler.go
│ │ └── start_info_handler.go
│ └── characters/
│ └── characters.go # Модель персонажа (GORM) и репозиторий для работы с данными персонажей
└── migrations/ # SQL-миграции для базы данных (golang-migrate)
├── 000001_create_characters_table.down.sql
├── 000001_create_characters_table.up.sql
├── ... # Остальные файлы миграций

## 🛠️ Миграции Базы Данных

Миграции находятся в директории `/migrations` и управляются с помощью `golang-migrate`.

*   **Создание новой миграции:**
    ```bash
    migrate create -ext sql -dir migrations -seq имя_миграции
    ```
*   **Применение миграций (пример для локальной разработки):**
    ```bash
    migrate -path ./migrations -database "postgres://user:pass@host:port/db?sslmode=disable" up
    ```
*   **Откат миграций:**
    ```bash
    migrate -path ./migrations -database "postgres://user:pass@host:port/db?sslmode=disable" down 1 # Откатить последнюю
    ```

## 🤝 Внесение Вклада (Contributing)

Буду рад вашим Pull Request'ам! Если у вас есть идеи по улучшению или вы нашли ошибку, пожалуйста, создайте Issue или Pull Request.

1.  Сделайте форк проекта.
2.  Создайте новую ветку (`git checkout -b feature/новая-фича`).
3.  Внесите свои изменения.
4.  Сделайте коммит (`git commit -am 'Добавлена новая фича'`).
5.  Отправьте изменения в ваш форк (`git push origin feature/новая-фича`).
6.  Создайте Pull Request.