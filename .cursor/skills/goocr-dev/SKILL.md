---
name: goocr-dev
description: Development guide for the PDF Intelligence Extraction Platform. Use when implementing Phase A/B/C tasks, hybrid OCR, platform adapters, Python layout/table services, or observability/K8s setup.
disable-model-invocation: true
---

# go-ocr Development Guide

Module: `github.com/lbkkk/OCR_go`. Full architecture: [docs/ARCHITECTURE.md](../../docs/ARCHITECTURE.md).

## Implementation phases (all mandatory, sequential)

### Phase A — Document core + CLI
- A1 detector, A2 extractor (pdf/docx), A3 hybrid OCR (Tesseract + Qwen + HybridEngine)
- A4 layout-service (Python gRPC), A5 table-service (Python gRPC), A6 Go gRPC clients
- A7 enricher, A8 markdown + converter, A9 CLI, A10 tests

### Phase B — Platform
- B1-B6 adapters (postgres, minio, rabbitmq, redis, usecase)
- B7-B9 REST + JWT + WebSocket + reprocess/search
- B10 worker, B11 docker-compose, B12 OpenAPI

### Phase C — Production
- C1 retry/DLQ/idempotency, C2 Kafka, C3-C6 observability (Prometheus/Grafana/Loki/Jaeger)
- C7-C9 Docker/K8s/CI-CD, C10 load test

## Hybrid OCR

Tesseract pass 1 → Qwen vision pass 2 (image + draft) → fallback Tesseract on error.

## Pipeline

```
detector -> extractor | hybrid OCR -> layout gRPC -> table gRPC -> enricher -> markdown
```

## Conventions

- English for code/comments/CLI. `gofmt`, wrap errors with `%w`, propagate `context.Context`.
- Secrets via env/config. Verify with `go build ./...` and `go test ./...`.
