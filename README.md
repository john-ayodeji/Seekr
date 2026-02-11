# Seekr

Seekr is a **mini, privacy-first search engine** focused on **usefulness over optimization**.

It is **not as good as Google** — and it is **not trying to be**.

Seekr exists as a small, honest experiment in search:
no tracking, no ads, no manipulation — just results.

---

## What Seekr Is (and Isn’t)

### What Seekr Is
- A mini search engine
- Privacy-first by default
- Opinionated about relevance
- Built for learning, experimentation, and clarity
- Transparent and explainable

### What Seekr Is Not
- A Google replacement
- A web-scale crawler
- An ad-driven platform
- A personalized surveillance system

If you need the entire internet indexed in milliseconds — use Google.  
If you want to understand **how search works** without the noise — use Seekr.

---

## Motivation

Search engines today optimize for:
- ads
- SEO gaming
- engagement metrics
- user profiling

Seekr is motivated by a simple question:

> What if search just tried to be helpful?

This project explores:
- relevance over ranking hacks
- intent over keywords
- simplicity over scale
- users over advertisers

---

## Quick Start

> ⚠️ Seekr is early-stage and under active development.

### Prerequisites
- Go (backend)
- Node.js / Bun (frontend, optional)
- Docker (optional but recommended)

### Clone the repository

```bash
git clone https://github.com/your-username/seekr.git
cd seekr
```

### Run the backend

```bash
go run ./cmd/server
```

---

## Environment Variables

Seekr is configured using environment variables.

### Required Variables

| Variable | Description | Example |
|--------|------------|---------|
| PORT | Port the server runs on | 3000 |
| DATABASE_URL | PostgreSQL connection string | postgres://user:password@localhost:5432/seekr |
| RABBITMQ_URL | RabbitMQ connection string | amqp://guest:guest@localhost:5672/ |


### Example `.env`

```env
PORT=3000
DATABASE_URL=postgres://seekr:seekr@localhost:5432/seekr
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
LOG_LEVEL=debug
ENV=development
```

---

## API Endpoints

### Search

```http
GET /api/v1/search?q=query
```

**Description:**  
Search the index for relevant documents based on the query.

**Query Params**
- `q` — search query string

---

### Submit Sitemap

```http
POST /api/v1/sitemap
```

**Body**
```json
{
  "sitemap_url": "https://example.com/sitemap.xml"
}
```

**Description:**  
Submits a sitemap URL for indexing.

---

## Usage

1. Start the server
2. Ensure Postgres and RabbitMQ are running
3. Submit a sitemap
4. Query the search endpoint

No tracking. No ads. No personalization.

---

## Contributing

Contributions are welcome, especially around:
- ranking logic
- indexing strategies
- query parsing
- performance improvements
- documentation and tests

### How to Contribute

1. Fork the repository
2. Create a feature branch
3. Keep changes focused
4. Open a PR with context

---

## Project Philosophy

Seekr does not chase scale.

It chases:
- clarity
- honesty
- usefulness

Big search engines optimize for growth.  
Seekr optimizes for understanding.
