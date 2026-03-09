# Business Domain Knowledge

## About This Project

We're building a unified commerce platform for small-to-medium e-commerce
merchants. Today these businesses juggle spreadsheets, Shopify, email threads,
and manual processes. The goal is one system that handles products, orders,
customers, and inventory — a single source of truth.

## What's In and What's Out

In scope for v1: product catalog management, order lifecycle, customer
accounts, inventory tracking across warehouses.

Out of scope: marketplace/multi-vendor support, subscription billing,
analytics dashboards. These are v2 considerations.

## Glossary

- **Customer** — A registered user who browses products and places orders.
- **Product** — An item for sale with a price, description, and stock level.
- **Order** — A transaction tying a customer to one or more products.
- **Inventory** — Stock counts per product per warehouse.
- **Line Item** — A single product entry in an order, with quantity and unit price.
- **Merchant** — The business operating the storefront.

## Who's Involved

- **Customers** browse, search, and purchase.
- **Merchant Admins** manage products and inventory, review orders.
- **Operations** handles picking, packing, and shipping.
- **Platform Engineering** builds and maintains the system.

## How the Business Works

A customer places zero or more orders. Each order has line items referencing
products. Every product tracks inventory per warehouse. Fulfilling an order
decrements inventory and triggers a shipment.

## The Typical Customer Journey

1. Customer browses the catalog and adds items to a cart.
2. Customer checks out; system validates stock and calculates the total.
3. Payment processes; order is confirmed.
4. Warehouse picks, packs, ships. Inventory decrements.
5. Customer gets a shipping confirmation with tracking.

## How Things Work Today

Most merchants use a combination of Shopify for their storefront, Google Sheets
for inventory, and email for order coordination. Data lives in multiple places
and frequently gets out of sync.

## What's Broken

- Inventory drifts between systems — merchants oversell regularly.
- Order processing is copy-paste between email and spreadsheets.
- No unified customer view; support checks 3+ tools per inquiry.
- Reporting requires manual aggregation across platforms every week.

## What Success Looks Like

- 95% of pilot merchants report faster order processing.
- Zero inventory discrepancies after 30 days of parallel operation.
- Customer-facing API uptime above 99.9% in the first quarter.