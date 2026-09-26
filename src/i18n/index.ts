import { en } from "./en";
import { id } from "./id";
import type { Locale, SiteCopy } from "../types/site";

export const translations: Record<Locale, SiteCopy> = { en, id };

export function getCopy(locale: Locale): SiteCopy {
  return translations[locale];
}
