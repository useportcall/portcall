# Deprecated Apps

These services are no longer part of active local/dev/prod runtime:

- `cron` logic moved to `libs/go/cronx`, started by `apps/dashboard`.
- `webhook` logic moved to `libs/go/webhookx`, routes served by `apps/dashboard`.

The folders remain for reference only and are intentionally excluded from the
active CLI deploy/run workflows.
