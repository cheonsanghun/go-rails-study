// This file mirrors backend domain models so API usage stays type-safe in React.

export type Artist = {
  id: string;
  name: string;
  labelName: string;
  country: string;
  primaryGenre: string;
  createdAt: string;
};

export type Track = {
  id: string;
  artistId: string;
  title: string;
  isrc: string;
  durationSeconds: number;
  genre: string;
  language: string;
  explicit: boolean;
  audioFileName: string;
  createdAt: string;
};

export type Store = {
  id: string;
  name: string;
  territory: string;
  supportsPreorder: boolean;
};

export type ReleaseStatus =
  | "draft"
  | "submitted"
  | "in_review"
  | "approved"
  | "delivered"
  | "live"
  | "rejected";

export type StoreDeliveryStatus =
  | "pending"
  | "processing"
  | "delivered"
  | "live"
  | "failed";

export type StoreDelivery = {
  storeId: string;
  storeName: string;
  status: StoreDeliveryStatus;
  lastUpdatedAt: string;
};

export type ReleaseRequest = {
  id: string;
  artistId: string;
  trackId: string;
  title: string;
  releaseDate: string;
  status: ReleaseStatus;
  storeDeliveries: StoreDelivery[];
  createdAt: string;
  updatedAt: string;
};

export type SalesReport = {
  id: string;
  artistId: string;
  trackId: string;
  trackTitle: string;
  storeId: string;
  storeName: string;
  period: string;
  country: string;
  currency: string;
  units: number;
  netRevenueMinor: number;
  reportedAt: string;
};

export type CreateTrackPayload = {
  title: string;
  isrc: string;
  durationSeconds: number;
  genre: string;
  language: string;
  explicit: boolean;
  audioFileName: string;
};

export type CreateReleaseRequestPayload = {
  trackId: string;
  title: string;
  releaseDate: string;
  storeIds: string[];
};

export type UpdateReleaseStatusPayload = {
  status?: ReleaseStatus;
  storeId?: string;
  storeStatus?: StoreDeliveryStatus;
};
