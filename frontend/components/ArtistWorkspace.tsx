"use client";

// This file coordinates client-side actions for one artist's distribution workflow.

import Link from "next/link";
import { useState } from "react";
import {
  createReleaseRequest,
  createTrack,
  updateReleaseStatus,
} from "@/lib/api";
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
import { ReleaseRequestForm } from "./ReleaseRequestForm";
import { ReleaseStatusPanel } from "./ReleaseStatusPanel";
import { SalesReportTable } from "./SalesReportTable";
import { TrackForm } from "./TrackForm";
import { TrackList } from "./TrackList";

type ArtistWorkspaceProps = {
  artist: Artist;
  initialTracks: Track[];
  initialReleaseRequests: ReleaseRequest[];
  stores: Store[];
  initialSalesReports: SalesReport[];
};

export function ArtistWorkspace({
  artist,
  initialTracks,
  initialReleaseRequests,
  stores,
  initialSalesReports,
}: ArtistWorkspaceProps) {
  const [tracks, setTracks] = useState(initialTracks);
  const [releaseRequests, setReleaseRequests] = useState(
    initialReleaseRequests,
  );
  const [feedback, setFeedback] = useState("");
  const [error, setError] = useState("");

  async function handleCreateTrack(payload: CreateTrackPayload) {
    setFeedback("");
    setError("");
    try {
      const track = await createTrack(artist.id, payload);
      setTracks((current) => [track, ...current]);
      setFeedback("Track registered.");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Track failed.");
    }
  }

  async function handleCreateRelease(payload: CreateReleaseRequestPayload) {
    setFeedback("");
    setError("");
    try {
      const releaseRequest = await createReleaseRequest(artist.id, payload);
      setReleaseRequests((current) => [releaseRequest, ...current]);
      setFeedback("Release request submitted.");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Release failed.");
    }
  }

  async function handleUpdateRelease(
    releaseId: string,
    payload: UpdateReleaseStatusPayload,
  ) {
    setFeedback("");
    setError("");
    try {
      const updated = await updateReleaseStatus(releaseId, payload);
      setReleaseRequests((current) =>
        current.map((request) => (request.id === updated.id ? updated : request)),
      );
      setFeedback("Distribution status updated.");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Update failed.");
    }
  }

  return (
    <main className="page-shell">
      <Link className="back-link" href="/">
        Back to artists
      </Link>

      <header className="workspace-header">
        <div>
          <p className="eyebrow">Artist Workspace</p>
          <h1>{artist.name}</h1>
          <p className="lede">
            {artist.labelName} / {artist.country} / {artist.primaryGenre}
          </p>
        </div>
        <div className="meta-row">
          <span className="pill">Artist ID: {artist.id}</span>
          <span className="pill">REST API connected</span>
        </div>
      </header>

      {feedback ? <p className="feedback">{feedback}</p> : null}
      {error ? <p className="feedback error">{error}</p> : null}

      <div className="workspace-grid">
        <div className="stack">
          <TrackList tracks={tracks} />
          <ReleaseStatusPanel
            releaseRequests={releaseRequests}
            onUpdate={handleUpdateRelease}
          />
          <SalesReportTable reports={initialSalesReports} />
        </div>

        <aside className="stack">
          <TrackForm artist={artist} onCreate={handleCreateTrack} />
          <ReleaseRequestForm
            tracks={tracks}
            stores={stores}
            onCreate={handleCreateRelease}
          />
        </aside>
      </div>
    </main>
  );
}
