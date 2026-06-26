
export type ChildFields = {
  name: string;
  birthday: string;              
  gender: string;
  about_short: string;
  interests: string[];
  activity_level: string;
  limitations: string[];
  allergies: string[];        
  play_styles: string[];
  interests_weight?: number | null;        
  activity_level_weight?: number | null;   
  limitations_weight?: number | null;      
  allergies_weight?: number | null;       
  play_styles_weight?: number | null;     
  max_age_difference?: number | null; 
};


const trimEq = (a: string, b: string) => a.trim() === b.trim();

const arrNormalize = (arr: string[]) =>
  arr.map(s => s.trim()).filter(Boolean);

const arrEqual = (a: string[], b: string[]) =>
  a.length === b.length && a.every((v, i) => v === b[i]);


export function buildChildPayload(
  current: ChildFields,
  original: ChildFields
): Record<string, unknown> {
  const payload: Record<string, unknown> = {};

  // Scalars
  if (!trimEq(current.name, original.name)) payload.name = current.name.trim();
  if (current.birthday !== original.birthday) payload.birthday = current.birthday; // assume valid YYYY-MM-DD
  if (current.gender !== original.gender) payload.gender = current.gender;
  if (!trimEq(current.about_short, original.about_short)) payload.about_short = current.about_short.trim();
  if (current.activity_level !== original.activity_level) payload.activity_level = current.activity_level;

  // Arrays (normalized)
  const interestsNow = arrNormalize(current.interests);
  const interestsOld = arrNormalize(original.interests);
  if (!arrEqual(interestsNow, interestsOld)) payload.interests = interestsNow;

  const limitationsNow = arrNormalize(current.limitations);
  const limitationsOld = arrNormalize(original.limitations);
  if (!arrEqual(limitationsNow, limitationsOld)) payload.limitations = limitationsNow;

  const allergiesNow = arrNormalize(current.allergies);
  const allergiesOld = arrNormalize(original.allergies);
  if (!arrEqual(allergiesNow, allergiesOld)) payload.allergies = allergiesNow;

  const playNow = arrNormalize(current.play_styles);
  const playOld = arrNormalize(original.play_styles);
  if (!arrEqual(playNow, playOld)) payload.play_styles = playNow;

  if (current.interests_weight !== original.interests_weight)
  payload.interests_weight = current.interests_weight;

if (current.activity_level_weight !== original.activity_level_weight)
  payload.activity_level_weight = current.activity_level_weight;

if (current.limitations_weight !== original.limitations_weight)
  payload.limitations_weight = current.limitations_weight;

if (current.allergies_weight !== original.allergies_weight)
  payload.allergies_weight = current.allergies_weight;

if (current.play_styles_weight !== original.play_styles_weight)
  payload.play_styles_weight = current.play_styles_weight;

// NEW: max age difference (>= 0)
if (current.max_age_difference !== original.max_age_difference)
  payload.max_age_difference = current.max_age_difference;

  return payload;
}
