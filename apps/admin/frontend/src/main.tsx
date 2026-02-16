import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider, useAuth } from "@/lib/keycloak/auth";
import Layout from "@/components/layout";
import DashboardView from "@/views/dashboard-view";
import AppsListView from "@/views/apps-list-view";
import AppDetailView from "@/views/app-detail-view";
import ConnectionsManagementView from "@/views/connections-management-view";
import DogfoodView from "@/views/dogfood-view";
import DogfoodUserDetailView from "@/views/dogfood-user-detail-view";
import QueuesView from "@/views/queues-view";
import JobsView from "@/views/jobs-view";
import "./index.css";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 1000, // 30 seconds
      retry: 1,
    },
  },
});

function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-900 via-purple-800 to-slate-900">
      <div className="text-center text-white space-y-6">
        <div className="flex items-center justify-center gap-3 mb-8">
          <svg
            className="h-16 w-16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
          >
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
          </svg>
        </div>
        <h1 className="text-4xl font-bold">Portcall Admin</h1>
        <p className="text-purple-200">Initializing authentication...</p>
        <div className="animate-pulse">
          <div className="h-2 w-48 bg-purple-600 rounded mx-auto"></div>
        </div>
      </div>
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </QueryClientProvider>
  );
}

function AppContent() {
  const { isLoaded, isAuthenticated, login } = useAuth();

  // Show loading screen while initializing
  if (!isLoaded) {
    return <LoginPage />;
  }

  // If not authenticated, show login button
  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-900 via-purple-800 to-slate-900">
        <div className="text-center text-white space-y-6">
          <div className="flex items-center justify-center gap-3 mb-8">
            <svg
              className="h-16 w-16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
            </svg>
          </div>
          <h1 className="text-4xl font-bold">Portcall Admin</h1>
          <button
            onClick={login}
            className="px-6 py-3 bg-purple-600 hover:bg-purple-700 rounded-lg font-medium transition-colors"
          >
            Sign In
          </button>
        </div>
      </div>
    );
  }

  // Authenticated - show the app
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<DashboardView />} />
          <Route path="/apps" element={<AppsListView />} />
          <Route path="/apps/:appId" element={<AppDetailView />} />
          <Route path="/connections" element={<ConnectionsManagementView />} />
          <Route path="/dogfood" element={<DogfoodView />} />
          <Route
            path="/dogfood/users/:userId"
            element={<DogfoodUserDetailView />}
          />
          <Route path="/queues" element={<QueuesView />} />
          <Route path="/jobs" element={<JobsView />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
