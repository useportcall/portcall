declare module "braintree-web-drop-in" {
  interface DropinCreateOptions {
    authorization: string;
    container: HTMLElement | string;
    card?: { vault?: { vaultCard?: boolean } };
  }
  interface DropinInstance {
    requestPaymentMethod(): Promise<{ nonce: string; type: string }>;
    teardown(): Promise<void>;
  }
  export function create(options: DropinCreateOptions): Promise<DropinInstance>;
}
