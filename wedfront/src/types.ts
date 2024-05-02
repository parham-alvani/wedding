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
  name: string;
  id: string;
  spouse: string;
  answer?: Answer;
}
