// import logger from "@/lib/logger";
import { useTranslation } from "react-i18next";
import { useRouteError } from "react-router-dom";

const isProduction: boolean = import.meta.env.PROD;

const ErrorPage = () => {
  const { t } = useTranslation();
  const error = useRouteError();

  return (
    <div className="w-screen h-screen flex justify-center items-center">
      <div className="flex flex-col gap-6">
        <h1 className="text-2xl font-semibold">{t("views.error.title")}</h1>
        <p>{t("views.error.description")}</p>
        {!isProduction && (
          <pre className="overflow-auto p-4 bg-zinc-900 text-zinc-50 rounded font-mono">
            {error ? (error as Error).stack : null}
          </pre>
        )}
      </div>
    </div>
  );
};

export default ErrorPage;
