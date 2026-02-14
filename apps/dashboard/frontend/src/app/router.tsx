import OnboardGuard from "@/components/onboard-guard";
import { AppProvider } from "@/components/app-provider";
import Layout from "@/components/layouts/layout";
import CompanyDetailsView from "@/features/company-details/company-details-view";
import DeveloperView from "@/features/developer/developer-view";
import ErrorPage from "@/features/error/error-view";
import HomeView from "@/features/home/home-view";
import InvoiceListView from "@/features/invoices/invoice-list-view";
import PaymentIntegrationsView from "@/features/payment-integrations/payment-integrations-view";
import PlanListView from "@/features/plans/plan-list-view";
import { EditPlan } from "@/features/plans/plan-view";
import QuoteListView from "@/features/quotes/quote-list-view";
import { QuoteView } from "@/features/quotes/quote-view";
import SubscriptionListView from "@/features/subscriptions/subscription-list-view";
import SubscriptionView from "@/features/subscriptions/subscription-view";
import { UsageView } from "@/features/usage";
import UserView from "@/features/users/user-view";
import UsersView from "@/features/users/users-view";
import { RedirectToSignIn, useAuth } from "@/lib/keycloak/auth";
import { Suspense } from "react";
import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";

const routes = [
  { index: true, element: <HomeView /> },
  { path: "/developer", element: <DeveloperView /> },
  { path: "/invoices", element: <InvoiceListView /> },
  { path: "/subscriptions", element: <SubscriptionListView /> },
  { path: "/subscriptions/:id", element: <SubscriptionView /> },
  { path: "/company", element: <CompanyDetailsView /> },
  { path: "/integrations", element: <PaymentIntegrationsView /> },
  { path: "/users", element: <UsersView /> },
  { path: "/users/:id", element: <UserView /> },
  { path: "/plans", element: <PlanListView /> },
  { path: "/plans/:id", element: <EditPlan /> },
  { path: "/quotes", element: <QuoteListView /> },
  { path: "/quotes/:id", element: <QuoteView /> },
  { path: "/usage", element: <UsageView /> },
];

function AuthenticatedLayout() {
  const { isLoaded, isAuthenticated } = useAuth();

  if (!isLoaded) {
    return (
      <div className="auth-transition-screen" role="status" aria-live="polite">
        <div className="auth-transition-card">Checking session...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <RedirectToSignIn />;
  }

  return (
    <Suspense fallback={<div className="w-full h-full bg-muted animate-pulse rounded" />}>
      <AppProvider>
        <OnboardGuard>
          <Layout>
            <Outlet />
          </Layout>
        </OnboardGuard>
      </AppProvider>
    </Suspense>
  );
}

export const router = createBrowserRouter([
  { path: "/", errorElement: <ErrorPage />, element: <AuthenticatedLayout />, children: routes },
  { path: "*", element: <Navigate to="/" /> },
]);
