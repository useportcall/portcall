-- Add billing_exempt column to apps table
-- Migration: 20260120_add_billing_exempt_to_apps
-- Description: Add billing_exempt flag for dogfood/internal apps that should skip billing logic

ALTER TABLE apps
ADD COLUMN IF NOT EXISTS billing_exempt BOOLEAN NOT NULL DEFAULT false;

-- Add index for efficient querying of non-exempt apps
CREATE INDEX IF NOT EXISTS idx_apps_billing_exempt ON apps(billing_exempt);

-- Update any existing 'Portcall Live' or 'Portcall Test' apps to be billing exempt
UPDATE apps
SET billing_exempt = true
WHERE name IN ('Portcall Live', 'Portcall Test')
  AND billing_exempt = false;
