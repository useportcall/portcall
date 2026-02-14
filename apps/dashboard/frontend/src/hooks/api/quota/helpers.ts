export function redirectWithQuery(param: string) {
  window.location.href = `${window.location.origin}${window.location.pathname}?${param}`;
}

export function redirectToCheckout(checkoutUrl: string) {
  window.location.href = checkoutUrl;
}
