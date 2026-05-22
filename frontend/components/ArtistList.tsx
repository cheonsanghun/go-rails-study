import Link from "next/link";
import type { Artist } from "@/types/domain";

// This file renders artists as the first navigation step into the workspace.

type ArtistListProps = {
  artists: Artist[];
};

export function ArtistList({ artists }: ArtistListProps) {
  return (
    <section className="section">
      <h2>아티스트 목록</h2>
      <div className="grid artist-grid">
        {artists.map((artist) => (
          <Link
            className="artist-card"
            href={`/artists/${artist.id}`}
            key={artist.id}
          >
            <div>
              <h3>{artist.name}</h3>
              <p className="muted">{artist.labelName}</p>
            </div>
            <div className="meta-row">
              <span className="pill">{artist.country}</span>
              <span className="pill">{artist.primaryGenre}</span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
