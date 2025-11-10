// import logger from "@/lib/logger";
import { useRouteError } from "react-router-dom";

const isProduction: boolean = import.meta.env.PROD;

const ErrorPage = () => {
  const error = useRouteError();

  return (
    <div className="w-screen h-screen flex justify-center items-center">
      <div className="flex flex-col gap-6">
        <h1 className="text-2xl font-semibold">Oops! Something went wrong.</h1>
        <p>Sorry, an unexpected error has occurred.</p>
        {!isProduction && (
          <pre className="overflow-auto p-4 bg-slate-900 text-slate-50 rounded font-mono">
            {error ? (error as Error).stack : null}
          </pre>
        )}
      </div>
    </div>
  );
};

export default ErrorPage;
