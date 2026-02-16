import { getPlan } from "@repo/ui/api/get-plan.ts";
import { SubscribeButton } from "@repo/ui/components/subscribe-button";
import { Check } from "@repo/ui/icons";
import { cn } from "@repo/ui/lib/utils.ts";
import { PlanSubmitButton } from "@repo/ui/submit-button";

const plans = [
  {
    name: "Free",
    description:
      "For individuals who want to try out the most advanced AI audio",
    price: 0,
    credits: "10k credits/month",
    features: [
      "Text to Speech",
      "Speech to Text",
      "Music",
      "Agents",
      "3 Projects in Studio",
      "Automated Dubbing",
      "API Access",
    ],
    cta: "GET STARTED",
    planId: process.env.NEXT_PUBLIC_ELEVENLABS_FREE_PLAN_ID || "",
    highlighted: false,
  },
  {
    name: "Starter",
    description: "For hobbyists creating projects with AI audio",
    price: 5,
    credits: "30k credits/month",
    features: [
      { content: "Everything in Free, plus", feature: false },
      "Commercial License",
      "Instant Voice Cloning",
      "20 Projects in Studio",
      "Dubbing Studio",
      "Music commercial use",
    ],
    cta: "GET STARTED",
    planId: process.env.NEXT_PUBLIC_ELEVENLABS_STARTER_PLAN_ID || "",
    highlighted: false,
  },
  {
    name: "Creator",
    popular: true,
    discount: "First month 50% off",
    originalPrice: 22,
    price: 11,
    credits: "100k credits/month",
    description: "For creators making premium content for global audiences",
    features: [
      { content: "Everything in Starter, plus", feature: false },
      "Professional Voice Cloning",
      "Additional Credits",
      "192kbps quality audio",
    ],
    cta: "GET STARTED",
    planId: process.env.NEXT_PUBLIC_ELEVENLABS_CREATOR_PLAN_ID || "",
    highlighted: true,
  },
  {
    name: "Pro",
    description: "For creators ramping up their content production",
    price: 99,
    credits: "500k credits/month",
    features: [
      { content: "Everything in Creator, plus", feature: false },
      "44.1kHz PCM audio output via API",
    ],
    cta: "GET STARTED",
    planId: process.env.NEXT_PUBLIC_ELEVENLABS_PRO_PLAN_ID || "",
    highlighted: false,
  },
];

export default async function PricingPage() {
  return (
    <>
      <div className="min-h-screen py-12 px-4 -z-10">
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="mb-12">
            <h1 className="font-light text-4xl text-gray-900 mb-4">Pricing</h1>
            <p className="text-xl text-gray-600">
              Plans built for creators and business of all sizes
            </p>
          </div>

          {/* Pricing Cards */}
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-4 gap-8">
            {plans.map((plan) => (
              <div
                key={plan.name}
                className={cn(
                  "relative border-0 shadow-xl rounded-3xl overflow-hidden py-5",
                  { "outline-4 outline-black": plan.highlighted }
                )}
              >
                <div className="px-4 flex flex-col  xl:min-h-40">
                  <h3 className="text-xl font-semibold">{plan.name}</h3>
                  <p className="text-sm text-gray-600 mt-2 min-h-12">
                    {plan.description}
                  </p>
                  <h4 className="text-sm text-black font-bold">
                    {plan.credits}
                  </h4>
                </div>

                <div className="mt-10 px-4">
                  <div className="mb-0 text-black">
                    <div className="mt-4 flex items-baseline">
                      <span className="text-4xl mr-1">${plan.price}</span>
                      <span className="text-sm text-gray-500 mt-2">
                        per month
                      </span>
                    </div>
                  </div>

                  <PlanSubmitButton
                    planId={plan.planId}
                    className={cn(
                      "w-full rounded-full mt-6 bg-gray-200 hover:bg-gray-300 text-sm text-gray-800 font-semibold py-3",
                      { "bg-black text-white": plan.highlighted }
                    )}
                  />

                  <ul className="space-y-3 mt-6">
                    {plan.features.map((feature, idx) => {
                      if (typeof feature === "object") {
                        return (
                          <li key={idx} className="text-sm mt-4">
                            {feature.content}
                          </li>
                        );
                      } else {
                        return (
                          <li key={idx} className="flex items-start">
                            <Check className="size-3 mr-3 mt-0.5 inline-block" />
                            <span className="text-sm text-black font-semibold">
                              {feature}
                            </span>
                          </li>
                        );
                      }
                    })}
                  </ul>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </>
  );
}
