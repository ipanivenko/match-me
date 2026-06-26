
import type { City } from "../../types/profile";

export type ProfileFields = {
  name: string;
  gender: string;
  about: string;
  preferredDistance: number;
  languages: string[];
  city: City | null; 
};

const arrEqual = (a: unknown[], b: unknown[]) =>
  a.length === b.length && a.every((v, i) => v === b[i]);

const cityEqual = (a: City | null, b: City | null) =>
  (!a && !b) ||
  (!!a &&
    !!b &&
    a.label === b.label &&
    (a.lat ?? null) === (b.lat ?? null) &&
    (a.lon ?? null) === (b.lon ?? null));

export function buildPayload(
  current: ProfileFields,
  original: ProfileFields
): Record<string, unknown> {
  const payload: Record<string, unknown> = {};

  // simple scalars
  if (current.name.trim() !== original.name.trim()) payload.name = current.name.trim();
  if (current.about.trim() !== original.about.trim()) payload.about = current.about.trim();
  if (current.gender !== original.gender) payload.gender = current.gender;
  if (current.preferredDistance !== original.preferredDistance)
    payload.preferredDistance = current.preferredDistance;


  // languages (array)
  if (!arrEqual(current.languages, original.languages)) {
    payload.languages = current.languages;
  }

  // city → send the fields backend expects.
  // API expects addressCity + lat/lon:
  if (!cityEqual(current.city, original.city)) {
    payload.addressCity = current.city?.label ?? "";
    if (current.city?.lat != null) payload.lat = current.city.lat;
    if (current.city?.lon != null) payload.lon = current.city.lon;
    
  }
  

  return payload;
}
