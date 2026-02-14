-- Make billing_address_id nullable on subscriptions table for free plans
ALTER TABLE subscriptions ALTER COLUMN billing_address_id DROP NOT NULL;
ALTER TABLE subscriptions ALTER COLUMN billing_address_id SET DEFAULT NULL;
