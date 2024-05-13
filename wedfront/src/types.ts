export interface CompatibilityItem {
  icon: string;
  title: string;
  url: string;
}

export interface FeatureItem {
  description: string;
  icon: string;
  title: string;
}

export interface FooterLink {
  description: string;
  icon: string;
  url: string;
}

export interface NavItem {
  title: string;
  url: string;
}

export interface Answer {
  coming: boolean;
  plus_one: boolean;
}
export interface Guest {
  first_name: string;
  last_name: string;
  id: string;
  is_family?: boolean;
  spouse_first_name?: string;
  spouse_last_name?: string;
  answer?: Answer;
}
