# Product Requirements

## What Users Need

- Customers: browse/search catalog, add to cart, place orders, view history, track shipments.
- Admins: add/update/remove products, view and manage all orders.

## Hard Constraints

- Every API endpoint requires JWT authentication.
- Order totals computed server-side only; client totals are ignored.
- Inventory can never go negative — reject orders that would understock.

## Acceptance Criteria

- Customer can view paginated product list (20/page default).
- Customer places an order and gets confirmation within 3 seconds.
- Admin can CRUD products via the admin API.
- Search returns results within 200ms for catalogs up to 100k.

## Privacy and Security

- PCI-DSS compliance for all payment handling.
- Sensitive user data encrypted at rest (AES-256) and in transit (TLS 1.3).
- Rate limiting: 100 requests/min per IP on public endpoints.

## Data Integrity Rules

- Prices are positive decimals, exactly 2 decimal places.
- Order status follows a strict state machine: pending → confirmed → shipped → delivered.
- No hard deletes — customers, orders, and products are soft-deleted only.

## Open Questions

- Shopify webhook payload format for inventory sync — needs a spike.
- Do existing merchants have consistent CSV formats for product import?
- No historical traffic baseline for new merchants; peak patterns unknown.

## What Can Go Wrong

- Stripe outage: orders queue with retry; customer sees "processing."
- Database failover: 10–30 second read-only window during promotion.
- S3 outage: product pages render without images; no hard failure.

## How We'll Validate

- Load test with 500 simulated concurrent users before beta.
- 2-week parallel run: compare old system vs. new system outputs.
- Automated regression suite covering every acceptance criterion.

## Delivery Phases

- Weeks 1–3: Core product CRUD + customer registration.
- Weeks 4–6: Order placement + inventory management.
- Weeks 7–9: Payment integration + admin dashboard.

## Go-to-Market

- Week 10: Beta with 10 pilot merchants.
- Week 12: Expand to 50 merchants.
- Week 16: General availability.
- GA + 4 weeks: Old system sunset.

## Stop Conditions

- Order error rate > 1% during beta → halt rollout.
- p95 latency > 5s on any endpoint → performance review.
- Critical failure → revert to legacy workflow within 2 hours.