import type {
  Artist,
  CreateReleaseRequestPayload,
  CreateTrackPayload,
  ReleaseRequest,
  SalesReport,
  Store,
  Track,
  UpdateReleaseStatusPayload,
} from "@/types/domain";

// This file centralizes REST calls so components do not know endpoint details.

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type ArtistsResponse = { artists: Artist[] };
type ArtistResponse = { artist: Artist };
type TracksResponse = { tracks: Track[] };
type TrackResponse = { track: Track };
type StoresResponse = { stores: Store[] };
type ReleaseRequestsResponse = { releaseRequests: ReleaseRequest[] };
type ReleaseRequestResponse = { releaseRequest: ReleaseRequest };
type SalesReportsResponse = { salesReports: SalesReport[] };
type ApiErrorResponse = { error?: string };

async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    cache: "no-store",
    headers: {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string> | undefined),
    },
  });

  if (!response.ok) {
    let message = `API request failed: ${response.status}`;
    try {
      const body = (await response.json()) as ApiErrorResponse;
      if (body.error) {
        message = body.error;
      }
    } catch {
      // Keep the HTTP status fallback when the response is not JSON.
    }
    throw new Error(message);
  }

  return response.json() as Promise<T>;
}

export async function getArtists(): Promise<Artist[]> {
  const data = await apiFetch<ArtistsResponse>("/api/artists");
  return data.artists;
}

export async function getArtist(artistId: string): Promise<Artist> {
  const data = await apiFetch<ArtistResponse>(`/api/artists/${artistId}`);
  return data.artist;
}

export async function getTracks(artistId: string): Promise<Track[]> {
  const data = await apiFetch<TracksResponse>(`/api/artists/${artistId}/tracks`);
  return data.tracks;
}

export async function createTrack(
  artistId: string,
  payload: CreateTrackPayload,
): Promise<Track> {
  const data = await apiFetch<TrackResponse>(`/api/artists/${artistId}/tracks`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
  return data.track;
}

export async function getStores(): Promise<Store[]> {
  const data = await apiFetch<StoresResponse>("/api/stores");
  return data.stores;
}

export async function getReleaseRequests(
  artistId: string,
): Promise<ReleaseRequest[]> {
  const data = await apiFetch<ReleaseRequestsResponse>(
    `/api/artists/${artistId}/release-requests`,
  );
  return data.releaseRequests;
}

export async function createReleaseRequest(
  artistId: string,
  payload: CreateReleaseRequestPayload,
): Promise<ReleaseRequest> {
  const data = await apiFetch<ReleaseRequestResponse>(
    `/api/artists/${artistId}/release-requests`,
    {
      method: "POST",
      body: JSON.stringify(payload),
    },
  );
  return data.releaseRequest;
}

export async function updateReleaseStatus(
  releaseId: string,
  payload: UpdateReleaseStatusPayload,
): Promise<ReleaseRequest> {
  const data = await apiFetch<ReleaseRequestResponse>(
    `/api/release-requests/${releaseId}/status`,
    {
      method: "PATCH",
      body: JSON.stringify(payload),
    },
  );
  return data.releaseRequest;
}

export async function getSalesReports(
  artistId: string,
): Promise<SalesReport[]> {
  const data = await apiFetch<SalesReportsResponse>(
    `/api/artists/${artistId}/sales-reports`,
  );
  return data.salesReports;
}
