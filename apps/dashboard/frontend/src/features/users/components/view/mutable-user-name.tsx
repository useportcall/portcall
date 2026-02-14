import { useUpdateUser } from "@/hooks";
import { User } from "@/models/user";
import { useTranslation } from "react-i18next";

export function MutableUserName({ user }: { user: User }) {
  const { t } = useTranslation();
  const { mutate } = useUpdateUser(user.id);

  return (
    <input
      className="outline-none"
      defaultValue={user.name}
      placeholder={t("views.user.fields.name_placeholder")}
      onBlur={(e) => {
        if (!e.target.value || e.target.value === user.name) return;
        mutate({ name: e.target.value });
      }}
    />
  );
}
