import { CreditCard, TestTube2, Landmark } from "lucide-react";

function StripeIcon({ className = "w-6 h-6" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="24" height="24" rx="4" fill="#635BFF" />
      <path d="M11.5 8.5c0-.83.68-1.5 1.51-1.5.83 0 1.49.67 1.49 1.5h2c0-1.93-1.57-3.5-3.49-3.5-1.93 0-3.51 1.57-3.51 3.5 0 3.5 5 2.91 5 5 0 .83-.67 1.5-1.5 1.5s-1.5-.67-1.5-1.5h-2c0 1.93 1.57 3.5 3.5 3.5s3.5-1.57 3.5-3.5c0-3.5-5-2.91-5-5z" fill="white" />
    </svg>
  );
}

function MockIcon({ className = "w-6 h-6" }: { className?: string }) {
  return <div className={`${className} rounded-lg bg-gradient-to-br from-amber-400 to-orange-500 flex items-center justify-center`}><TestTube2 className="w-4 h-4 text-white" /></div>;
}

function BraintreeIcon({ className = "w-6 h-6" }: { className?: string }) {
  return (
    <div className={`${className} rounded-lg bg-gradient-to-br from-blue-500 to-blue-700 flex items-center justify-center`}>
      <Landmark className="w-4 h-4 text-white" />
    </div>
  );
}

function DefaultIcon({ className = "w-6 h-6" }: { className?: string }) {
  return <div className={`${className} rounded-lg bg-muted flex items-center justify-center`}><CreditCard className="w-4 h-4 text-muted-foreground" /></div>;
}

export function getProviderIcon(source: string, className?: string) {
  switch (source.toLowerCase()) {
    case "stripe":
      return <StripeIcon className={className} />;
    case "braintree":
      return <BraintreeIcon className={className} />;
    case "local":
    case "mock":
      return <MockIcon className={className} />;
    default:
      return <DefaultIcon className={className} />;
  }
}

export function getProviderDisplayName(source: string) {
  if (source.toLowerCase() === "stripe") return "Stripe";
  if (source.toLowerCase() === "braintree") return "Braintree";
  if (source.toLowerCase() === "local") return "Mock Provider";
  return source.charAt(0).toUpperCase() + source.slice(1);
}

export function getProviderDescription(source: string) {
  if (source.toLowerCase() === "stripe") return "Production-ready payment processing";
  if (source.toLowerCase() === "braintree") return "PayPal / Braintree payment processing";
  if (source.toLowerCase() === "local") return "For testing and development";
  return "Payment provider";
}
