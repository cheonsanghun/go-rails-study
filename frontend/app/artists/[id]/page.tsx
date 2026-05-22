import Link from "next/link";
import { ArtistWorkspace } from "@/components/ArtistWorkspace";
import {
  getArtist,
  getReleaseRequests,
  getSalesReports,
  getStores,
  getTracks,
} from "@/lib/api";
import type {
  Artist,
  ReleaseRequest,
  SalesReport,
  Store,
  Track,
} from "@/types/domain";

// This file loads all data needed for one artist's distribution workspace.

type ArtistPageProps = {
  params: Promise<{
    id: string;
  }>;
};

export default async function ArtistPage({ params }: ArtistPageProps) {
  const { id } = await params;
  let artist: Artist | null = null;
  let tracks: Track[] = [];
  let releases: ReleaseRequest[] = [];
  let stores: Store[] = [];
  let salesReports: SalesReport[] = [];
  let errorMessage = "";

  try {
    [artist, tracks, releases, stores, salesReports] = await Promise.all([
      getArtist(id),
      getTracks(id),
      getReleaseRequests(id),
      getStores(),
      getSalesReports(id),
    ]);
  } catch (error) {
    errorMessage =
      error instanceof Error
        ? error.message
        : "Failed to load artist workspace";
  }

  if (errorMessage || !artist) {
    return (
      <main className="page-shell">
        <Link className="back-link" href="/">
          Back to artists
        </Link>
        <div className="notice">
          Artist workspace를 불러올 수 없습니다: {errorMessage}
        </div>
      </main>
    );
  }

  return (
    <ArtistWorkspace
      artist={artist}
      initialTracks={tracks}
      initialReleaseRequests={releases}
      stores={stores}
      initialSalesReports={salesReports}
    />
  );
}
