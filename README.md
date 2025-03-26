# Server and Clients Simulation Project

## Overview
A Go-based simulation of a rate-limited HTTP server with multiple clients sending requests and collecting statistics.

## Features

### Server Implementation
- **Configurable port** from ENV file
- Supports **GET/POST** requests
- **Rate limiting**: 5 requests/second
- **Randomized responses** (70% positive / 30% negative)
- **Statistics endpoint** (`GET /stats`)

### Clients
- **Posting Clients (x2)**:
  - Each sends **100 POST requests** via goroutines
  - **2 workers** per client (5 reqs/sec limit each)
  - Collects detailed response statistics
- **Monitoring Client**:
  - Checks server status **every 5 seconds**
  - Simple availability monitoring

### Statistics Tracking
- **Response tracking**:
  - Per-client breakdown
  - Server-wide totals
- **JSON endpoint** for aggregated data
- **Client-side stats** displayed after completion

## Technical Details
- **Response codes**:
  - Positive: `200 OK`, `202 Accepted`
  - Negative: `400 Bad Request`, `500 Internal Error`
- **Rate limiting** with overflow handling
- **Concurrent processing** with goroutines
