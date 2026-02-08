# Technical Specification

## System Overview

The platform is a Node.js/Express backend backed by PostgreSQL, with Redis for
caching and rate limiting. Product images go to S3. Payments through Stripe.
Email via SendGrid. The frontend is a React SPA consuming REST APIs.

## What We're Replacing

- Per-merchant spreadsheets → centralized PostgreSQL database.
- Email-based order flow → REST API with automated state transitions.
- Nightly batch inventory reconciliation → real-time sync.

## What We're Keeping

- Shopify stays as the storefront for merchants who use it; we integrate via webhooks.
- Stripe remains the payment processor.

## External Dependencies

- PostgreSQL 15+ (primary data store).
- Redis (sessions, rate limiting, caching).
- Stripe API (payments).
- AWS S3 (product images).
- SendGrid (transactional email).

## Migration Strategy

- Phase 1: Deploy new API alongside existing tools; merchants opt in.
- Phase 2: Import historical data (last 2 years) per merchant via CSV.
- Phase 3: Sunset old workflows once each merchant confirms parity.

## Backward Compatibility

- Shopify product catalogs importable via CSV.
- Historical order data migrated without loss.
- Old and new systems run in parallel during transition.

## Tricky Scenarios

- Two customers order the last unit simultaneously (inventory race condition).
- Products with variants (size, color) sharing a base SKU.
- Split shipments when items come from different warehouses.

## Design Decisions and Trade-offs

- PostgreSQL over DynamoDB: consistency wins, manual scaling is acceptable.
- REST over GraphQL: simpler for v1 clients, revisit later.
- Synchronous inventory check at checkout: adds ~50ms latency but prevents overselling.

## Infrastructure Requirements

- 500 concurrent users with sub-second response times.
- Point-in-time database recovery, 5-minute RPO.
- Docker-based deployments, zero-downtime rolling updates.

## Latency and Throughput Targets

- Product list: p95 < 200ms (up to 100k products).
- Order placement: end-to-end < 3 seconds.
- Search: < 200ms for 100k catalog.
- Inventory lookup: < 50ms per SKU.