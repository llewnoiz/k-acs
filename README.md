# gtr069 — TR-069 Inform ACS + 장치 목록

CPE(장치)가 보내는 TR-069(CWMP) **Inform** 메시지를 수신·저장하는 **Go ACS** 와, 그 장치를
목록으로 보여주는 **Next.js** 화면으로 구성된 데모 프로젝트.

```
gtr069/
├── acs/   # Go ACS 백엔드 (Inform 수신 + SQLite 저장 + JSON API)
└── web/   # Next.js 장치 목록 UI
```

## 아키텍처

```
 CPE ──(1) POST /acs  cwmp:Inform SOAP──▶  Go ACS (:7547) ──▶ SQLite (acs.db)
     ◀─(2) cwmp:InformResponse───────────       │
 CPE ──(3) POST /acs  (empty)──────────▶        │ GET /api/devices (JSON)
     ◀─(4) 204 No Content──────────────         ▼
                                          Next.js (:3000)  ── 5초 폴링 ──▶ 테이블
```

- **범위**: Inform 수신 전용 POC (InformResponse 반환 + 빈 POST 204 세션 종료).
- **저장소**: SQLite (`acs/acs.db`, `modernc.org/sqlite` — CGO 불필요).
- **보안**: 평문 HTTP, 인증 없음 (로컬 데모용).

## 1. ACS 백엔드 실행

```bash
cd acs
go run .          # :7547 에서 리슨. DB=acs.db 자동 생성
```

환경변수(선택): `ACS_ADDR`(기본 `:7547`), `ACS_DB`(기본 `acs.db`).

엔드포인트:

| 메서드 | 경로 | 설명 |
|--------|------|------|
| POST | `/acs` | CWMP Inform 수신 → InformResponse. 빈 POST → 204 |
| GET | `/api/devices` | 장치 목록 (JSON 배열) |
| GET | `/api/devices/{serial}` | 장치 상세 (전체 파라미터 포함) |

### 단위 테스트

```bash
cd acs && go test ./...
```

## 2. Inform 시뮬레이션 (장치 등록)

```bash
# 샘플 장치 SN-0001 등록
curl -X POST http://localhost:7547/acs \
  -H 'Content-Type: text/xml' \
  --data-binary @acs/testdata/inform_sample.xml

# 다른 장치로 등록 (serial/제조사/product class 치환)
sed -e 's/SN-0001/SN-0002/' -e 's/Acme Networks/Globex/' -e 's/Router-X100/Gateway-G9/' \
  acs/testdata/inform_sample.xml \
  | curl -X POST http://localhost:7547/acs -H 'Content-Type: text/xml' --data-binary @-

# 등록 결과 확인
curl http://localhost:7547/api/devices
```

## 3. 프런트엔드 실행

```bash
cd web
npm install
npm run dev       # http://localhost:3000
```

`web/app` 의 장치 목록 화면이 5초마다 `/api/devices` 를 폴링하여 새 Inform 을 자동 반영한다.
ACS 주소를 바꾸려면 `web/.env.local` 에 `NEXT_PUBLIC_ACS_API` 를 설정한다
(`.env.local.example` 참고).

## 확장 여지 (범위 밖)

- GetParameterValues / SetParameterValues 등 큐잉 RPC 처리로 세션 완성형 확장
- ACS → CPE Connection Request (역방향 세션 개시)
- HTTP Basic 인증 + TLS
