"use client";

import "../i18n";
import { PropsWithChildren } from "react";

export default function I18nProvider({ children }: PropsWithChildren) {
    return <>{children}</>;
}
