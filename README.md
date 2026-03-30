# ⚡ High-Load URL Shortener

> Сервис сокращения ссылок с фокусом на производительность — аналог bit.ly на Go

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat-square&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-UI-85EA2D?style=flat-square&logo=swagger&logoColor=black)
![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)

---

## 🏗️ Архитектура

```
┌─────────────┐     ┌──────────────────────────────────────────┐
│   Client    │────▶│              Rate Limiter                 │
└─────────────┘     └──────────────────┬───────────────────────┘
                                       │
                    ┌──────────────────▼───────────────────────┐
                    │              API Handler                  │
                    │   POST /shorten    GET /:id               │
                    └──────┬───────────────────┬───────────────┘
                           │                   │
                    ┌──────▼──────┐     ┌──────▼──────┐
                    │  LRU Cache  │     │  LRU Cache  │
                    │  (memory)   │     │  (memory)   │
                    └──────┬──────┘     └──────┬──────┘
                           │                   │ miss
                    ┌──────▼───────────────────▼──────┐
                    │           PostgreSQL             │
                    │     Connection Pool (25)         │
                    └─────────────────────────────────┘
```

### Жизненный цикл запроса

```
POST /shorten {"url": "https://google.com"}
    → Генерация base62 ID (8 символов)
    → Сохранить в LRU Cache + PostgreSQL
    → Вернуть {"id": "gcMsE4ec"}

GET /gcMsE4ec
    → Найти в LRU Cache → HIT  → 301 Redirect (< 5ms)
    → Найти в LRU Cache → MISS → PostgreSQL → 301 Redirect
    → Не найдено → 404
```

---

## 🚀 Запуск

### Требования
- Go 1.23+
- Docker + Docker Compose

### Установка

```bash
# Клонировать репозиторий
git clone https://github.com/your-username/url-shortener
cd url-shortener

# Поднять PostgreSQL
docker-compose up -d

# Установить зависимости
go mod tidy

# Запустить
go run main.go
```

### API

```bash
# Сократить ссылку
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
# → {"id": "gcMsE4ec"}

# Перейти по ссылке
curl -L http://localhost:8080/gcMsE4ec
```

### Swagger UI

```
http://localhost:8080/swagger/
```

---

## 📊 Результаты нагрузочного тестирования

Инструмент: **k6** | Нагрузка: **500 concurrent users** | Длительность: **2 минуты**

### Сценарий: POST /shorten (только запись)

| Метрика | Значение |
|---|---|
| RPS | **373 req/s** |
| Avg latency | 6.88ms |
| Median | 4.26ms |
| p(95) | 11.96ms |
| Error rate | **0%** |

### Сценарий: POST + GET (смешанная нагрузка, уникальные URL)

| Метрика | Значение |
|---|---|
| RPS | **1043 req/s** |
| Avg latency | 108ms |
| Median | 46ms |
| p(95) | 370ms |
| Error rate | **0%** |

### Прогресс оптимизаций

| Этап | RPS | p(95) |
|---|---|---|
| Baseline (50 VUs) | 36 | 244ms |
| + Connection Pool (500 VUs) | 373 | 22ms |
| + LRU Cache (500 VUs) | 373 | **12ms** |

---

## 🔧 Стек

| Технология | Для чего |
|---|---|
| `net/http` | HTTP сервер без сторонних фреймворков |
| `sync.RWMutex` + `atomic` | Потокобезопасность без лишних блокировок |
| `hashicorp/golang-lru` | LRU cache — автовытеснение старых записей |
| `PostgreSQL` + connection pool | Персистентное хранилище, 25 соединений |
| `Docker Compose` | Воспроизводимое окружение |
| `Swagger` | Автодокументация API |
| `k6` | Нагрузочное тестирование |

---

## 💡 Технические решения

**LRU Cache вместо обычного map** — ограничивает потребление памяти. При 1М+ ссылок обычный map съест всю память, LRU хранит только N последних популярных.

**Connection Pool** — переиспользование соединений с PostgreSQL дало 10x прирост RPS и 4x снижение latency.

**atomic.AddInt64 для счётчика кликов** — инкремент без мьютекса. Мьютекс держим только на чтение указателя из LRU, сам инкремент атомарный.

**Graceful Shutdown** — при Ctrl+C сервер дожидается текущих запросов, корректно закрывает соединение с БД.

---

