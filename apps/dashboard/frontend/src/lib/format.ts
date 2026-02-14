import { Plan } from "@/models/plan";
import { Quote } from "@/models/quote";

/**
 * Formatting utilities for consistent data display across the application.
 */

/**
 * Format a plan interval into a short string (mo, qtr, yr, wk).
 */
export function formatInterval(interval: string): string {
    const intervalMap: Record<string, string> = {
        month: "mo",
        quarter: "qtr",
        year: "yr",
        week: "wk",
    };
    return intervalMap[interval] || interval;
}

/**
 * Format currency amount (in cents) with currency code and interval.
 */
export function formatCurrency(
    amount: number,
    currency: string,
    interval?: string
): string {
    const currencyUpper = currency.toUpperCase();
    const fee = Number(amount / 100).toFixed(2);

    if (interval) {
        return `${currencyUpper} ${fee} / ${formatInterval(interval)}`;
    }

    return `${currencyUpper} ${fee}`;
}

/**
 * Format fixed charge from a plan.
 * Returns a formatted string like "USD 49.00 / mo" or "Free" for zero amount.
 */
export function formatFixedCharge(plan: Plan): string {
    const fixedPlanItem = plan.items.find(
        (item) => item.pricing_model === "fixed"
    );

    if (!fixedPlanItem) return "";

    if (fixedPlanItem.unit_amount === 0) return "Free";

    return formatCurrency(
        fixedPlanItem.unit_amount,
        plan.currency,
        plan.interval
    );
}

/**
 * Format fixed charge from a quote.
 * Returns a formatted string like "USD 49.00 / mo" or "Free" for zero amount.
 */
export function formatQuoteFixedCharge(quote: Quote): string {
    if (!quote.plan) return "—";

    const fixedPlanItem = quote.plan.items?.find(
        (item) => item.pricing_model === "fixed"
    );

    if (!fixedPlanItem) return "—";

    if (fixedPlanItem.unit_amount === 0) return "Free";

    return formatCurrency(
        fixedPlanItem.unit_amount,
        quote.plan.currency,
        quote.plan.interval
    );
}

/**
 * Format a date string for display.
 */
export function formatDate(
    dateStr: string,
    options?: Intl.DateTimeFormatOptions
): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString(
        "en-GB",
        options || {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
        }
    );
}
