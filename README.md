# Artist Distribution Study

Wano / TuneCore Japan 같은 음악 유통 서비스를 의식해서 만든 학습용 미니 프로젝트입니다. 인디 아티스트가 곡을 등록하고, 여러 스토어에 배포 신청을 만들고, 스토어별 배포 상태와 간단한 수익 리포트를 확인하는 흐름을 다룹니다.

## 목표

- TypeScript로 도메인 타입을 명확히 설계한다.
- Next.js App Router로 화면과 클라이언트 액션을 구현한다.
- Go 표준 라이브러리로 REST API를 만든다.
- PostgreSQL 교체를 염두에 두고 repository 인터페이스를 분리한다.
- 프론트엔드와 백엔드가 REST API로 연결되는 구조를 이해한다.

## 폴더 구조

```text
artist-distribution-study/
  frontend/
    app/          Next.js App Router 페이지
    components/   화면 단위 React 컴포넌트
    lib/          API 호출 함수
    types/        TypeScript 도메인 타입
  backend/
    cmd/          실행 진입점
    internal/
      domain/     Go 도메인 모델
      handler/    HTTP 핸들러
      service/    비즈니스 로직
      repository/ 저장소 인터페이스와 in-memory 구현
```

## 실행 방법

백엔드 파일 구조와 계층별 책임은 [`backend/ARCHITECTURE.md`](backend/ARCHITECTURE.md)에 정리했습니다.

### 1. Backend 실행

```bash
cd artist-distribution-study/backend
go run ./cmd/api
```

Backend는 `http://localhost:8080`에서 실행됩니다.

확인:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/artists
```

### 2. Frontend 실행

다른 터미널에서 실행합니다.

```bash
cd artist-distribution-study/frontend
npm install
npm run dev
```

브라우저에서 `http://localhost:3000`을 엽니다.

API 주소를 바꾸고 싶다면 다음처럼 실행할 수 있습니다.

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run dev
```

## REST API

| Method | Path | 설명 |
| --- | --- | --- |
| `GET` | `/api/artists` | 아티스트 목록 조회 |
| `GET` | `/api/artists/{artistId}` | 아티스트 상세 조회 |
| `GET` | `/api/artists/{artistId}/tracks` | 곡 목록 조회 |
| `POST` | `/api/artists/{artistId}/tracks` | 곡 등록 |
| `GET` | `/api/stores` | 배포 스토어 목록 조회 |
| `GET` | `/api/artists/{artistId}/release-requests` | 배포 신청 목록 조회 |
| `POST` | `/api/artists/{artistId}/release-requests` | 배포 신청 생성 |
| `PATCH` | `/api/release-requests/{releaseId}/status` | 전체 배포 상태 또는 스토어별 상태 변경 |
| `GET` | `/api/artists/{artistId}/sales-reports` | 수익 리포트 조회 |

## 예시 요청

곡 등록:

```bash
curl -X POST http://localhost:8080/api/artists/art_1001/tracks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "New City Lights",
    "isrc": "JPA1B2600999",
    "durationSeconds": 205,
    "genre": "Electronic",
    "language": "ja",
    "explicit": false,
    "audioFileName": "new-city-lights.wav"
  }'
```

배포 신청:

```bash
curl -X POST http://localhost:8080/api/artists/art_1001/release-requests \
  -H "Content-Type: application/json" \
  -d '{
    "trackId": "trk_1001",
    "title": "Late Night Conveyor - Single",
    "releaseDate": "2026-06-12",
    "storeIds": ["spotify", "apple_music", "line_music"]
  }'
```

스토어별 상태 변경:

```bash
curl -X PATCH http://localhost:8080/api/release-requests/rel_1001/status \
  -H "Content-Type: application/json" \
  -d '{
    "storeId": "spotify",
    "storeStatus": "live"
  }'
```

전체 배포 상태 변경:

```bash
curl -X PATCH http://localhost:8080/api/release-requests/rel_1001/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "live"
  }'
```

## PostgreSQL로 바꿀 때의 DB 설계 연습

처음 구현은 in-memory repository이지만, 아래 테이블로 바꾸기 쉽게 인터페이스를 나눠두었습니다.

```sql
CREATE TABLE artists (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  label_name TEXT NOT NULL,
  country TEXT NOT NULL,
  primary_genre TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE tracks (
  id TEXT PRIMARY KEY,
  artist_id TEXT NOT NULL REFERENCES artists(id),
  title TEXT NOT NULL,
  isrc TEXT,
  duration_seconds INTEGER NOT NULL,
  genre TEXT NOT NULL,
  language TEXT NOT NULL,
  explicit BOOLEAN NOT NULL,
  audio_file_name TEXT,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE stores (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  territory TEXT NOT NULL,
  supports_preorder BOOLEAN NOT NULL
);

CREATE TABLE release_requests (
  id TEXT PRIMARY KEY,
  artist_id TEXT NOT NULL REFERENCES artists(id),
  track_id TEXT NOT NULL REFERENCES tracks(id),
  title TEXT NOT NULL,
  release_date DATE NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE store_deliveries (
  release_request_id TEXT NOT NULL REFERENCES release_requests(id),
  store_id TEXT NOT NULL REFERENCES stores(id),
  status TEXT NOT NULL,
  last_updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (release_request_id, store_id)
);

CREATE TABLE sales_reports (
  id TEXT PRIMARY KEY,
  artist_id TEXT NOT NULL REFERENCES artists(id),
  track_id TEXT NOT NULL REFERENCES tracks(id),
  store_id TEXT NOT NULL REFERENCES stores(id),
  period TEXT NOT NULL,
  country TEXT NOT NULL,
  currency TEXT NOT NULL,
  units INTEGER NOT NULL,
  net_revenue_minor INTEGER NOT NULL,
  reported_at TIMESTAMPTZ NOT NULL
);
```

## 다음으로 공부할 개선 과제

1. PostgreSQL repository를 추가하고 `MemoryRepository`와 교체해보기.
2. `ReleaseStatus` 상태 전이를 제한하기. 예: `submitted`에서 바로 `live`로 못 가게 만들기.
3. 곡 등록 폼에 UPC, 작사가, 작곡가, 음원 파일 검증 같은 메타데이터 추가하기.
4. 수익 리포트를 월별, 스토어별, 국가별로 집계하는 API 추가하기.
5. Next.js에서 로딩 UI, 에러 UI, optimistic update를 적용해보기.
6. Go handler 테스트와 service 테스트를 추가하기.
7. 인증을 붙여서 아티스트별 접근 권한을 분리해보기.
8. OpenAPI 문서를 작성하고 프론트엔드 타입을 자동 생성해보기.
