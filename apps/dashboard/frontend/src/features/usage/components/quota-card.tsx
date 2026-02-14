import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { QuotaStatus } from "@/hooks/api/quota";
import { AlertTriangle } from "lucide-react";
import { useTranslation } from "react-i18next";

export function QuotaCard({ title, description, icon: Icon, quota }: { title: string; description: string; icon: React.ComponentType<{ className?: string }>; quota?: QuotaStatus }) {
  const { t } = useTranslation();
  if (!quota) return <Card><CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2"><CardTitle className="text-sm font-medium">{title}</CardTitle><Icon className="h-4 w-4 text-muted-foreground" /></CardHeader><CardContent><p className="text-muted-foreground text-sm">{t("common.loading")}</p></CardContent></Card>;

  const percentage = quota.is_unlimited ? 0 : Math.min(100, Math.round((quota.usage / quota.quota) * 100));
  const isNearLimit = !quota.is_unlimited && percentage >= 80;
  const isExceeded = quota.is_exceeded;

  return (
    <Card className={isExceeded ? "border-red-200 bg-red-50/50" : ""}><CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2"><div><CardTitle className="text-sm font-medium">{title}</CardTitle><CardDescription className="text-xs mt-1">{description}</CardDescription></div><Icon className={`h-4 w-4 ${isExceeded ? "text-red-500" : "text-muted-foreground"}`} /></CardHeader><CardContent className="space-y-3"><div className="flex items-end justify-between"><div className="text-2xl font-bold">{quota.usage?.toLocaleString()}</div><div className="text-sm text-muted-foreground">{quota.is_unlimited ? <Badge variant="secondary">{t("views.usage.quotas.unlimited")}</Badge> : <span>{t("views.usage.quotas.of", { count: quota.quota ?? 0 })}</span>}</div></div>{!quota.is_unlimited && <Progress value={percentage} className={isExceeded ? "bg-red-100" : isNearLimit ? "bg-amber-100" : ""} />}{isExceeded && <AlertMessage text={t("views.usage.quotas.exceeded")} className="text-red-600" />}{isNearLimit && !isExceeded && <AlertMessage text={t("views.usage.quotas.approaching", { count: quota.remaining })} className="text-amber-600" />}</CardContent></Card>
  );
}

function AlertMessage({ text, className }: { text: string; className: string }) {
  return <div className={`flex items-center gap-2 text-sm ${className}`}><AlertTriangle className="h-4 w-4" /><span>{text}</span></div>;
}
