# Admin API

This API is used for debugging and triggering various queue jobs.

## Endpoints

- `GET /ping`: Health check endpoint.
- `POST /queues/:queue_id/tasks/:task_id`: Enqueue jobs to a queue.

## Example usage

```bash
curl --location "http://localhost:9100/queues/email_queue/tasks/send_invoice_paid_email" \
-X POST \
--header "Accept: application/json" \
--header 'Content-Type: application/json' \
--data '{
    "invoice_number": "INV-000009",
    "amount_paid": "$56.00",
    "year": 2025,
    "date_paid": "15th November, 2025",
    "company_name": "The Prancing Pony",
    "recipient_name": "Frodo Baggins"
}'
```
