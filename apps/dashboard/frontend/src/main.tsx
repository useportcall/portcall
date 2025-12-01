import "./index.css";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode, Suspense } from "react";
import { createRoot } from "react-dom/client";
import {
  createBrowserRouter,
  Navigate,
  Outlet,
  RouterProvider,
} from "react-router-dom";
import { Toaster } from "sonner";
import OnboardGuard from "./components/onboard-guard";
import Layout from "./components/layouts/layout";
import { TooltipProvider } from "./components/ui/tooltip";
import CompanyDetailsView from "./features/company-details/company-details-view";
import DeveloperView from "./features/developer/developer-view";
import ErrorPage from "./features/error/error-view";
import HomeView from "./features/home/home-view";
import InvoiceListView from "./features/invoices/invoice-list-view";
import PaymentIntegrationsView from "./features/payment-integrations/payment-integrations-view";
import PlanListView from "./features/plans/plan-list-view";
import { EditPlan } from "./features/plans/plan-view";
import SubscriptionListView from "./features/subscriptions/subscription-list-view";
import SubscriptionView from "./features/subscriptions/subscription-view";
import UserView from "./features/users/user-view";
import UsersView from "./features/users/users-view";
import "./index.css";
import {
  AuthProvider,
  RedirectToSignIn,
  SignedIn,
  SignedOut,
} from "./lib/keycloak/auth";
import { AppProvider } from "./components/app-provider";
import QuoteListView from "./features/quotes/quote-list-view";
import { QuoteView } from "./features/quotes/quote-view";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: "/",
    errorElement: <ErrorPage />,
    element: (
      <>
        <SignedIn>
          <Suspense
            fallback={
              <div className="w-full h-full bg-gray-100 animate-pulse rounded" />
            }
          >
            <AppProvider>
              <OnboardGuard>
                <Layout>
                  <Outlet />
                </Layout>
              </OnboardGuard>
            </AppProvider>
          </Suspense>
        </SignedIn>
        <SignedOut>
          <RedirectToSignIn />
        </SignedOut>
      </>
    ),
    children: [
      {
        index: true,
        element: <HomeView />,
      },
      {
        path: "/developer",
        element: <DeveloperView />,
      },

      {
        path: "/invoices",
        element: <InvoiceListView />,
      },
      {
        path: "/subscriptions",
        element: <SubscriptionListView />,
      },
      {
        path: "/subscriptions/:id",
        element: <SubscriptionView />,
      },
      {
        path: "/company",
        element: <CompanyDetailsView />,
      },
      {
        path: "/integrations",
        element: <PaymentIntegrationsView />,
      },
      {
        path: "/users",
        element: <UsersView />,
      },
      {
        path: "/users/:id",
        element: <UserView />,
      },
      {
        path: "/plans",
        element: <PlanListView />,
      },
      {
        path: "/plans/:id",
        element: <EditPlan />,
      },
      {
        path: "/quotes",
        element: <QuoteListView />,
      },
      {
        path: "/quotes/:id",
        element: <QuoteView />,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" />,
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AuthProvider>
      <TooltipProvider>
        <QueryClientProvider client={queryClient}>
          <Toaster />
          <RouterProvider router={router} />
        </QueryClientProvider>
      </TooltipProvider>
    </AuthProvider>
  </StrictMode>
);
