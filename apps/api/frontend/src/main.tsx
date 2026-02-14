import "./index.css";
import { StrictMode, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import HomePage from "./pages/home";
import ErrorPage from "./pages/error";

// Dark mode handler
function DarkModeWrapper({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // Check system preference
    const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    if (isDark) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }

    // Listen for changes
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const handleChange = (e: MediaQueryListEvent) => {
      if (e.matches) {
        document.documentElement.classList.add("dark");
      } else {
        document.documentElement.classList.remove("dark");
      }
    };

    mediaQuery.addEventListener("change", handleChange);
    return () => mediaQuery.removeEventListener("change", handleChange);
  }, []);

  return <>{children}</>;
}

const router = createBrowserRouter(
  [
    {
      path: "/",
      element: <HomePage />,
      errorElement: <ErrorPage />,
    },
  ],
  {
    basename: "/docs",
  },
);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <DarkModeWrapper>
      <RouterProvider router={router} />
    </DarkModeWrapper>
  </StrictMode>,
);
