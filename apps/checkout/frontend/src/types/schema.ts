import z from "zod";

export const CheckoutFormSchema = z.object({
  line1: z.string().min(1, "validation.line1_required"),
  line2: z.string().optional(),
  city: z.string().min(1, "validation.city_required"),
  state: z.string().optional(),
  postal_code: z.string().min(1, "validation.postal_code_required"),
  country: z.string().min(1, "validation.country_required"),
  vat_number: z.string().optional(),
  is_business: z.boolean().default(false).optional(),
});

export type CheckoutFormSchema = z.infer<typeof CheckoutFormSchema>;
