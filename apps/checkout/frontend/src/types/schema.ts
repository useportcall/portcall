import z from "zod";

export const CheckoutFormSchema = z.object({
  line1: z.string().min(1, "Line 1 is required"),
  line2: z.string().optional(),
  city: z.string().min(1, "City is required"),
  state: z.string().optional(),
  postal_code: z.string().min(1, "Postal code is required"),
  country: z.string().min(1, "Country is required"),
  vat_number: z.string().optional(),
  is_business: z.boolean().default(false).optional(),
});

export type CheckoutFormSchema = z.infer<typeof CheckoutFormSchema>;
