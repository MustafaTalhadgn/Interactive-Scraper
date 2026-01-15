# Database Schema Documentation

## Overview
This document describes the database schema used for the Threat Intelligence Platform.

---

## Table: `sources`

### Description
Stores monitored sources used for intelligence collection such as forums, markets, and paste sites.

### Columns

| Column Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | SERIAL | PK | Unique source identifier |
| `name` | VARCHAR(255) | NOT NULL | Source name |
| `url` | TEXT | NOT NULL, UNIQUE | Source URL |
| `category` | VARCHAR(50) | - | Source category (forum, market, paste) |
| `criticality` | VARCHAR(20) | CHECK | Importance level |
| `enabled` | BOOLEAN | DEFAULT true | Source status |
| `scrape_interval` | INTERVAL | DEFAULT 1 hour | Scraping frequency |
| `last_scraped_at` | TIMESTAMP | - | Last scrape timestamp |
| `created_at` | TIMESTAMP | DEFAULT now() | Creation time |

---

## Table: `intelligence_data`

### Description
Stores analyzed intelligence data collected from monitored sources.

### Columns

| Column Name | Type | Constraints | Description |
|------------|------|-------------|-------------|
| `id` | SERIAL | PK | Unique record ID |
| `source_id` | INTEGER | FK → sources(id) | Related source |
| `title` | VARCHAR(500) | - | Content title |
| `analyzed_content` | TEXT | NOT NULL | Processed content |
| `source_url` | TEXT | NOT NULL | Original content URL |
| `criticality_score` | INT | CHECK (0–100) | Severity score |
| `published_at` | TIMESTAMP | - | Publish time |
| `created_at` | TIMESTAMP | DEFAULT now() | Record creation time |

---

## Relationships

- `sources (1) → (N) intelligence_data`
- Cascading delete enabled (`ON DELETE CASCADE`)

---

## Notes
- Criticality values are standardized across the system
- Scraping intervals are configurable per source
