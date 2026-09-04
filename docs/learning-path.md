# Learning And Build Path

This path follows official documentation first, then maps each language concept to one Quokka-Lab feature. Build in small increments: every phase should leave the application runnable.

## Official Documentation Sources

- Go documentation: https://go.dev/doc/
- Effective Go: https://go.dev/doc/effective_go
- Go standard library: https://pkg.go.dev/std
- Vue guide: https://vuejs.org/guide/introduction.html
- Vue TypeScript guide: https://vuejs.org/guide/typescript/overview.html
- Vite guide: https://vite.dev/guide/
- Rails guides: https://guides.rubyonrails.org/
- Rails API applications: https://guides.rubyonrails.org/api_app.html
- Active Record basics: https://guides.rubyonrails.org/active_record_basics.html
- C++ reference: https://en.cppreference.com/w/
- CMake tutorial: https://cmake.org/cmake/help/latest/guide/tutorial/
- libsndfile documentation: https://libsndfile.github.io/libsndfile/
- gRPC Go documentation: https://grpc.io/docs/languages/go/
- Cap'n Proto documentation: https://capnproto.org/
- Docker Compose documentation: https://docs.docker.com/compose/

## Priority Order

| Priority | Build target | Why it comes here |
| --- | --- | --- |
| P0 | Docker Compose baseline | Everyone can run the same stack locally |
| P1 | Vue keyboard + Web Audio | Fastest visible feedback and core product identity |
| P2 | Go health/API gateway | Establishes the new public backend boundary |
| P3 | Rails composition storage | Makes recordings durable and shareable |
| P4 | WebSocket collaboration in Go | Uses Go where it is strongest: concurrent connections |
| P5 | C++ CLI audio processing | Adds CPU-heavy DSP without blocking the web tier |
| P6 | Go-to-C++ UDS protocol | Replaces process startup with a resident worker |
| P7 | Go-to-Rails gRPC | Replaces internal HTTP proxying with typed contracts |
| P8 | Auth, limits, observability | Hardens the product after the main loop works |

## Stage 1: Web Audio Basics

Study:

- Vue single-file components and Composition API.
- TypeScript interfaces and event typing.
- Web Audio API concepts: `AudioContext`, `OscillatorNode`, `GainNode`, and `AnalyserNode`.

Build:

- Render the two-octave keyboard.
- Map computer keys to notes.
- Play a note with one oscillator and one gain envelope.
- Draw waveform or spectrum data from the analyser.

Done when:

- `npm run build` succeeds.
- A note can be played with both mouse and keyboard.

## Stage 2: Browser Recording

Study:

- `MediaRecorder`.
- Blob URLs and downloads.
- Multipart form uploads.

Build:

- Record the Web Audio destination stream.
- Pause, resume, stop, preview, and download.
- Upload a `.webm` recording to `/api/v1/compositions`.

Done when:

- A browser recording can be uploaded and then appears in the public playlist.

## Stage 3: Go Core Service

Study:

- Go modules.
- `net/http`, request routing, JSON encoding, contexts, and graceful shutdown.
- gorilla/websocket connection lifecycle.

Build:

- `GET /api/v1/health`.
- AI placeholder endpoints.
- Rails reverse proxy for persistence APIs.
- `/ws/collaboration/:room_id` room broadcast.

Done when:

- `go test ./...` succeeds.
- Frontend API calls pass through Go rather than Rails directly.

## Stage 4: Rails Data Authority

Study:

- Rails API mode.
- Active Record models, validations, migrations, and associations.
- Active Storage for uploaded audio.
- Devise JWT for user accounts.

Build:

- Composition, track, effect, playlist, comment, and like models.
- CRUD endpoints for compositions.
- Upload storage through Active Storage.
- Demo seed data.

Done when:

- `bin/rails db:prepare` succeeds.
- `GET /api/v1/compositions` returns public records through the Go gateway.

## Stage 5: C++ Audio Engine

Study:

- Modern C++ value types and standard containers.
- CMake targets and compiler flags.
- libsndfile WAV read/write workflow.
- Basic DSP: feedback delay, Schroeder reverb, resampling tradeoffs.

Build:

- `quokka_audio input.wav output.wav --reverb 0.8`.
- Delay and pitch-shift flags.
- Clear error output and non-zero exit codes on failure.

Done when:

- Docker builds the `quokka_audio` executable.
- A WAV file can be processed from the command line.

## Stage 6: Internal Protocols

Study:

- gRPC service definitions and generated Go/Ruby integration options.
- Unix domain sockets.
- Cap'n Proto schemas and zero-copy message layout.

Build:

- Generate Go gRPC stubs from `rails_service.proto`.
- Add an internal Rails endpoint or gRPC server for quota and composition lookups.
- Convert the C++ CLI into a resident worker that listens on a UDS path.

Done when:

- Go can call Rails through a typed contract.
- Go can submit audio jobs to the resident C++ worker without starting a new process per job.

## Stage 7: Production Hardening

Study:

- Docker Compose production profiles or Kubernetes manifests.
- Redis Streams consumer groups.
- Structured logs, metrics, tracing, and health checks.
- JWT validation, rate limiting, and upload validation.

Build:

- Redis-backed task queue.
- Gateway rate limiter.
- Authenticated write APIs.
- Observability dashboard hooks.

Done when:

- The platform can be scaled by adding Go/C++ worker replicas while keeping Rails and PostgreSQL isolated behind internal networking.
