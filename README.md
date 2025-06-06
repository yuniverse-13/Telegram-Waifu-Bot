# Waifu Telegram Bot (Телеграм-бот "Вайфу")

Телеграм-бот на Go для поиска, отображения информации и **оценки** персонажей (вайфу) из вашей базы данных PostgreSQL.

## ✨ Возможности

*   Поиск персонажей по ID.
*   Поиск персонажей по имени или альтернативному имени (регистронезависимый для основного имени).
*   Отображение случайного персонажа.
*   **Система рейтинга персонажей:**
    *   Оценка персонажа командой `/rate <ID или Имя> <оценка 1-10>`.
    *   Интерактивная оценка через кнопки (1⭐ - 10⭐) на карточке персонажа.
    *   Отображение информации о персонаже: имя, ID, тайтл, описание, изображение, **средний рейтинг**, **количество голосов** и **ваша персональная оценка**.
*   Использование PostgreSQL для хранения данных.
*   Использование GORM как ORM для взаимодействия с базой данных.
*   Управление миграциями схемы БД с помощью `golang-migrate`.
*   Возможность запуска в Docker-контейнере с помощью Docker Compose.
*   Форматирование сообщений с помощью MarkdownV2 для лучшей читаемости.

## 📋 Требования

*   Go (версия 1.21 или выше рекомендуется; ваш проект использует 1.24)
*   PostgreSQL (версия 13 или выше рекомендуется; ваш проект упоминает 17)
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
    *(Замените `YourUsername/telegram-waifu-bot` на URL вашего репозитория)*

2.  **Настройте переменные окружения:**
    Создайте файл `.env` в корне проекта, скопировав `.env.example` (если он у вас есть) или создав новый:
    ```env
    # Токен вашего Telegram бота
    TELEGRAM_BOT_TOKEN=ВАШ_ТЕЛЕГРАМ_БОТ_ТОКЕН

    # Настройки для подключения к PostgreSQL из docker-compose
    # Эти переменные используются сервисом postgres в docker-compose.yml для инициализации БД
    POSTGRES_USER=waifu_bot_user
    POSTGRES_PASSWORD=supersecretpassword
    POSTGRES_DB=waifu_bot_db
    # Порт, который PostgreSQL внутри Docker будет слушать и который может быть проброшен на хост
    # DB_PORT_HOST=5432 # Используется, если вы хотите пробросить порт на хост-машину (см. docker-compose.yml)

    # DATABASE_URL для самого бота (он будет обращаться к БД внутри Docker-сети)
    # Имя хоста 'postgres_db' соответствует имени сервиса БД в docker-compose.yml
    DATABASE_URL=postgres://waifu_bot_user:supersecretpassword@postgres_db:5432/waifu_bot_db?sslmode=disable

    # Следующие переменные передаются в docker-compose.yml для сервиса waifu_bot,
    # чтобы он мог сформировать DATABASE_URL или для других нужд.
    # В текущей конфигурации docker-compose.yml бота используется напрямую DATABASE_URL.
    # Если ваш docker-compose.yml собирает DATABASE_URL из этих частей, оставьте их.
    # DB_USER=${POSTGRES_USER}
    # DB_PASSWORD=${POSTGRES_PASSWORD}
    # DB_NAME=${POSTGRES_DB}
    # DB_HOST=postgres_db # Имя сервиса БД в Docker-сети
    # DB_PORT_CONTAINER=5432 # Порт БД внутри контейнера
    ```
    **Важно:** Замените `ВАШ_ТЕЛЕГРАМ_БОТ_ТОКЕН` на реальный токен вашего бота, полученный от [@BotFather](https://t.me/BotFather). Пароль и другие данные также рекомендуется сменить на более надежные. Убедитесь, что ваш `docker-compose.yml` правильно использует эти переменные (особенно `DATABASE_URL` для сервиса бота).

3.  **Запустите сервисы:**
    ```bash
    docker-compose up --build -d
    ```
    Опция `-d` запускает контейнеры в фоновом режиме. При первом запуске миграции базы данных будут применены автоматически (согласно `entrypoint.sh`).

4.  **Просмотр логов (если нужно):**
    ```bash
    docker-compose logs -f waifu_bot
    docker-compose logs -f postgres_db
    ```

5.  **Остановка:**
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
    TELEGRAM_BOT_TOKEN=ВАШ_ТЕЛЕГРАМ_БОТ_ТОКЕН
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
    # Убедитесь, что переменная DATABASE_URL установлена в вашем окружении или передайте ее явно:
    # export DATABASE_URL="postgres://ВАШ_ПОЛЬЗОВАТЕЛЬ_БД:ВАШ_ПАРОЛЬ_БД@localhost:5432/ВАША_БАЗА_ДАННЫХ?sslmode=disable"
    migrate -path ./migrations -database "${DATABASE_URL}" up
    ```

6.  **Запустите бота:**
    ```bash
    go run cmd/bot/main.go
    ```

## ⚙️ Конфигурация

Бот конфигурируется с помощью переменных окружения (обычно через файл `.env`):

*   `TELEGRAM_BOT_TOKEN`: **Обязательно.** Токен вашего Telegram бота.
*   `DATABASE_URL`: **Обязательно.** Строка подключения к PostgreSQL.
    *   Формат: `postgres://[user]:[password]@[host]:[port]/[dbname]?sslmode=[disable|require|etc.]`
*   (Для Docker Compose также используются `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` для инициализации сервиса БД, как указано в разделе установки).

## 🤖 Использование (Команды бота)

*   `/start` - Приветственное сообщение и список команд.
*   `/info` - Показывает список доступных команд.
*   `/character <ID или Имя>` - Поиск персонажа по его ID или имени/альтернативному имени.
    *   Пример: `/character 1`
    *   Пример: `/character Фрирен`
    При отображении карточки персонажа будет показан его средний рейтинг, количество голосов и ваша личная оценка (если есть). Также будет доступна кнопка "Оценить ✨".
*   `/randomcharacter` - Показывает случайного персонажа из базы данных. Карточка также будет содержать информацию о рейтинге и кнопку для оценки.
*   `/rate <ID или Имя> <оценка>` - Позволяет оценить персонажа. Оценка должна быть числом от 1 до 10.
    *   Пример: `/rate 1 10`
    *   Пример: `/rate Фрирен 10`
*   **Интерактивная оценка:**
    *   На карточке персонажа нажмите кнопку **"Оценить ✨"**.
    *   Вам будут предложены кнопки с оценками от 1⭐ до 10⭐.
    *   Выберите одну из них, чтобы поставить или обновить свою оценку. Карточка персонажа обновится, отобразив актуальную информацию о рейтинге.

## 🛠️ Миграции Базы Данных

Миграции находятся в директории `/migrations` и управляются с помощью `golang-migrate`. Они определяют схему вашей базы данных и позволяют её версионировать.

*   **Создание новой миграции:**
    ```bash
    migrate create -ext sql -dir migrations -seq имя_миграции
    ```
    Это создаст два файла: `..._имя_миграции.up.sql` (для изменений "вперед") и `..._имя_миграции.down.sql` (для отката изменений).

*   **Применение миграций (пример для локальной разработки):**
    ```bash
    migrate -path ./migrations -database "postgres://user:pass@host:port/db?sslmode=disable" up
    ```

*   **Откат миграций:**
    ```bash
    migrate -path ./migrations -database "postgres://user:pass@host:port/db?sslmode=disable" down 1 # Откатить последнюю примененную миграцию
    migrate -path ./migrations -database "postgres://user:pass@host:port/db?sslmode=disable" down # Откатить все миграции
    ```

## 🤝 Внесение Вклада (Contributing)

Буду рад вашим Pull Request'ам! Если у вас есть идеи по улучшению или вы нашли ошибку, пожалуйста, создайте Issue или Pull Request.

1.  Сделайте форк проекта.
2.  Создайте новую ветку (`git checkout -b feature/новая-фича`).
3.  Внесите свои изменения.
4.  Сделайте коммит (`git commit -am 'Добавлена новая фича'`).
5.  Отправьте изменения в ваш форк (`git push origin feature/новая-фича`).
6.  Создайте Pull Request.