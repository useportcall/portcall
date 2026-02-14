-- Add item-level billing intervals to plan_items, subscription_items, and billing_meters
-- Migration: 20260131_add_item_level_billing_intervals
-- Description: Enables different billing intervals for individual plan items (e.g., annual fixed price + monthly metered)

-- Add interval columns to plan_items
ALTER TABLE plan_items
ADD COLUMN IF NOT EXISTS interval VARCHAR(255) NOT NULL DEFAULT 'inherit',
ADD COLUMN IF NOT EXISTS interval_count INTEGER NOT NULL DEFAULT 1;

-- Add interval and reset columns to subscription_items
ALTER TABLE subscription_items
ADD COLUMN IF NOT EXISTS interval VARCHAR(255) NOT NULL DEFAULT 'inherit',
ADD COLUMN IF NOT EXISTS interval_count INTEGER NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS last_reset_at TIMESTAMP NULL,
ADD COLUMN IF NOT EXISTS next_reset_at TIMESTAMP NULL;

-- Add interval columns to billing_meters
ALTER TABLE billing_meters
ADD COLUMN IF NOT EXISTS interval VARCHAR(255) NOT NULL DEFAULT 'inherit',
ADD COLUMN IF NOT EXISTS interval_count INTEGER NOT NULL DEFAULT 1;

-- Create index for efficient querying of subscription items by next reset
CREATE INDEX IF NOT EXISTS idx_subscription_items_next_reset ON subscription_items(subscription_id, next_reset_at);

-- Create index for efficient querying of billing meters by next reset
CREATE INDEX IF NOT EXISTS idx_billing_meters_next_reset ON billing_meters(subscription_id, next_reset_at);

-- Update existing fixed plan items to use 'inherit' (default behavior)
-- Metered items can be updated manually to have different intervals
COMMENT ON COLUMN plan_items.interval IS 'Billing interval for this item: inherit (from plan), week, month, year';
COMMENT ON COLUMN plan_items.interval_count IS 'Number of intervals for the billing cycle (e.g., 1 month, 3 months)';
COMMENT ON COLUMN subscription_items.interval IS 'Billing interval for this item: inherit (from subscription), week, month, year';
COMMENT ON COLUMN subscription_items.next_reset_at IS 'When this subscription item will next be billed/reset';
