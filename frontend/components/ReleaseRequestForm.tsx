"use client";

// This file creates release requests with track, date, and target store choices.

import { FormEvent, useMemo, useState } from "react";
import type {
  CreateReleaseRequestPayload,
  Store,
  Track,
} from "@/types/domain";

type ReleaseRequestFormProps = {
  tracks: Track[];
  stores: Store[];
  onCreate: (payload: CreateReleaseRequestPayload) => Promise<void>;
};

export function ReleaseRequestForm({
  tracks,
  stores,
  onCreate,
}: ReleaseRequestFormProps) {
  const defaultDate = useMemo(() => {
    const date = new Date();
    date.setDate(date.getDate() + 21);
    return date.toISOString().slice(0, 10);
  }, []);

  const [trackId, setTrackId] = useState(tracks[0]?.id ?? "");
  const [title, setTitle] = useState("");
  const [releaseDate, setReleaseDate] = useState(defaultDate);
  const [storeIds, setStoreIds] = useState(stores.map((store) => store.id));
  const [saving, setSaving] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSaving(true);
    await onCreate({ trackId, title, releaseDate, storeIds });
    setTitle("");
    setReleaseDate(defaultDate);
    setSaving(false);
  }

  function toggleStore(storeId: string) {
    setStoreIds((current) =>
      current.includes(storeId)
        ? current.filter((id) => id !== storeId)
        : [...current, storeId],
    );
  }

  return (
    <section className="panel">
      <h2>배포 신청 생성</h2>
      <form className="form-grid" onSubmit={handleSubmit}>
        <label className="full-span">
          Track
          <select
            disabled={tracks.length === 0}
            required
            value={trackId}
            onChange={(event) => setTrackId(event.target.value)}
          >
            {tracks.map((track) => (
              <option key={track.id} value={track.id}>
                {track.title}
              </option>
            ))}
          </select>
        </label>

        <label className="full-span">
          Release title
          <input
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Blank uses track title + Single"
          />
        </label>

        <label className="full-span">
          Release date
          <input
            required
            type="date"
            value={releaseDate}
            onChange={(event) => setReleaseDate(event.target.value)}
          />
        </label>

        <div className="full-span stack">
          {stores.map((store) => (
            <label className="checkbox-row" key={store.id}>
              <input
                checked={storeIds.includes(store.id)}
                type="checkbox"
                onChange={() => toggleStore(store.id)}
              />
              {store.name} ({store.territory})
            </label>
          ))}
        </div>

        <button
          className="full-span"
          disabled={saving || tracks.length === 0 || storeIds.length === 0}
          type="submit"
        >
          {saving ? "Submitting..." : "Submit release"}
        </button>
      </form>
    </section>
  );
}
