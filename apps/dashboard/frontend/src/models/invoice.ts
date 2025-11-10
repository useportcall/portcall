export type Invoice = {
  id: string;
  subscription_id?: string;
  invoice_number: string;
  currency: string;
  due_by?: string;
  status: string;
  total?: number;
  subtotal?: number;
  pdf_url?: string;
  email_url?: string;
  recipient_email: string;
  recipient_id: string;
  created_at: string;
  updated_at: string;
  items: any[];
};
