import { ArtistList } from "@/components/ArtistList";
import { getArtists } from "@/lib/api";
import type { Artist } from "@/types/domain";

// This file renders the home page with the artist list entry point.

export default async function HomePage() {
  let artists: Artist[] = [];
  let errorMessage = "";

  try {
    artists = await getArtists();
  } catch (error) {
    errorMessage =
      error instanceof Error
        ? error.message
        : "Failed to load artists from API";
  }

  return (
    <main className="page-shell">
      <section className="hero">
        <div>
          <p className="eyebrow">Artist Distribution Study</p>
          <h1>인디 아티스트 음악 유통 관리</h1>
          <p className="lede">
            곡 등록, 배포 신청, 스토어별 진행 상태, 수익 리포트를 작은
            REST API와 Next.js 화면으로 연결해보는 학습용 예제입니다.
          </p>
        </div>
      </section>

      {errorMessage ? (
        <div className="notice">
          Backend API에 연결할 수 없습니다: {errorMessage}
        </div>
      ) : (
        <ArtistList artists={artists} />
      )}
    </main>
  );
}
