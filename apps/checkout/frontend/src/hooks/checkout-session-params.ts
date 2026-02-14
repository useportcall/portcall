const CHECKOUT_SESSION_ID_PATTERN = /^cs_[a-f0-9]{32}$/;
const PAYMENT_LINK_ID_PATTERN = /^pl_[a-f0-9]{32}$/;

export type CheckoutSessionCredentials = {
  id: string;
  token: string;
};

export type PaymentLinkCredentials = {
  id: string;
  token: string;
};

export function readCheckoutSessionCredentials(
  search: string,
): CheckoutSessionCredentials | null {
  const params = new URLSearchParams(search);
  const id = params.get("id")?.trim() ?? "";
  const token = params.get("st")?.trim() ?? "";

  if (!CHECKOUT_SESSION_ID_PATTERN.test(id) || token.length === 0) {
    return null;
  }

  return { id, token };
}

export function readPaymentLinkCredentials(
  search: string,
): PaymentLinkCredentials | null {
  const params = new URLSearchParams(search);
  const id = params.get("pl")?.trim() ?? "";
  const token = params.get("pt")?.trim() ?? "";

  if (!PAYMENT_LINK_ID_PATTERN.test(id)) {
    return null;
  }

  return { id, token };
}

export function stripCheckoutTokenFromURL(): void {
  const url = new URL(window.location.href);
  if (!url.searchParams.has("st") && !url.searchParams.has("pt")) {
    return;
  }

  url.searchParams.delete("st");
  url.searchParams.delete("pt");
  const query = url.searchParams.toString();
  const next = query ? `${url.pathname}?${query}` : url.pathname;
  window.history.replaceState(null, "", next);
}
