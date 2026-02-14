"use client";

import { CheckoutSession } from "@/types/api";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import axios from "axios";
import {
  CheckoutSessionCredentials,
  readPaymentLinkCredentials,
  readCheckoutSessionCredentials,
  stripCheckoutTokenFromURL,
} from "./checkout-session-params";
import { redeemPaymentLink } from "./redeem-payment-link";
import { getSessionErrorMessage } from "./checkout-session-error";

export function useGetCheckoutSession() {
  const didInit = useRef(false);
  const mounted = useRef(false);
  const [credentials, setCredentials] =
    useState<CheckoutSessionCredentials | null>(null);
  const [isInvalidLink, setIsInvalidLink] = useState(false);
  const [isRedeemingLink, setIsRedeemingLink] = useState(false);
  const [linkError, setLinkError] = useState<string | null>(null);

  useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);

  useEffect(() => {
    if (didInit.current) {
      return;
    }
    didInit.current = true;

    const parsed = readCheckoutSessionCredentials(window.location.search);
    if (parsed) {
      setCredentials(parsed);
      stripCheckoutTokenFromURL();
      return;
    }

    const paymentLink = readPaymentLinkCredentials(window.location.search);
    if (!paymentLink) {
      setIsInvalidLink(true);
      return;
    }

    stripCheckoutTokenFromURL();
    redeemPaymentLink(
      paymentLink,
      () => !mounted.current,
      setIsRedeemingLink,
      setIsInvalidLink,
      setLinkError,
    );
  }, []);

  const query = useQuery({
    queryKey: ["checkout-sessions", credentials?.id],
    queryFn: async () => {
      const path = "/api/checkout-sessions/" + credentials!.id;
      const { data } = await axios.get<{ data: CheckoutSession }>(path, {
        headers: {
          "X-Checkout-Session-Token": credentials!.token,
          "Cache-Control": "no-store",
        },
      });

      return data.data;
    },
    enabled: !!credentials,
  });

  const isWaitingForCredentials = !credentials && !isInvalidLink;
  const isSessionError = !!credentials && query.isError;
  const sessionError = getSessionErrorMessage(query.error);

  return {
    ...query,
    credentials,
    linkError,
    sessionError,
    isInvalidLink,
    isSessionError,
    isLoading:
      isRedeemingLink ||
      isWaitingForCredentials ||
      (!!credentials && query.isLoading),
  };
}
