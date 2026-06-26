export type MeResponse = {
  name?: string | null;
  gender?: string | null;
  about?: string | null;
  languages?: string[] | null;
  addressCity?: string | null;
  lat?: number | null;
  lon?: number | null;
  preferredDistance?:number | null;
  avatarurl?: string | null;
};

export type City = {
  label: string;
  countryCode?: string;
  lat: number;
  lon: number;
};

export type BioResponse = {
  gender?: string | null;
  prefferedDistance?: number | null;
};

export type PhotoUrl = {
  photo_url?: string | null;
}

export type ChildResponse = {
  name?: string | null;
	birthday?: string | null;
	gender?: string | null;
	about_short?: string | null;
	interests?: string[] | null;
	activity_level?: string | null;
	limitations?: string[] | null;
	allergies? : string[] | null;
	play_styles? : string[] | null;
  interests_weight?: number | null;        
  activity_level_weight?: number | null;   
  limitations_weight?: number | null;      
  allergies_weight?: number | null;       
  play_styles_weight?: number | null;     
  max_age_difference?: number | null; 
}


export type CombinedMe = MeResponse & BioResponse & PhotoUrl;

export type ChildProfile = {
  name: string;
  ageYears: number;
  gender: string;
  aboutShort: string;
  topInterests: string[];
};

export type UserProfile = {
  name: string;           // parent name
  about: string;
  languages: string[];
  addressCity: string;
  child: ChildProfile;    // child.name  
};

export type UserPhoto = {
  avatarurl?: string | null;
} 

export type CombinedUser = UserProfile & UserPhoto;
