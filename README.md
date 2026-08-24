# Quokka

Quokka is a web-based music creation and sharing platform. It lets users play notes in the browser, record melodies, upload compositions, share public works, and render offline audio effects through a C++ DSP engine.

The project is designed as a multi-language reference application:

- Vue 3 and TypeScript power the browser studio.
- Ruby on Rails provides the primary API and data model.
- C++ handles offline audio processing.
- PHP provides a lightweight compatibility API.
- Docker Compose wires the development stack together.

Quokka also reserves an internal integration point for Liora through `POST /internal/liora_trigger`.

## Features

- Two-octave browser keyboard from C4 to B5.
- Computer keyboard and mouse input support.
- Web Audio synthesis with `sine`, `square`, `sawtooth`, and `triangle` waveforms.
- MediaRecorder-based recording with playback and download.
- Recording upload to the Rails API.
- Public composition list from PostgreSQL.
- Likes and comments for public compositions.
- Optional sign-up and login through Devise JWT.
- Offline C++ audio processing for reverb, delay, and basic pitch shifting.
- Rails health check endpoint.
- Action Cable collaboration channel scaffold.
- Sidekiq background job scaffold.
- PHP compatibility API scaffold.

## Architecture

```text
Quokka/
├── frontend/        Vue 3 + TypeScript studio
├── backend_ruby/    Rails 8 API, PostgreSQL, Active Storage, Devise JWT
├── backend_php/     Lightweight PHP API compatibility layer
├── audio_engine/    C++20 command-line audio DSP engine
├── infra/nginx/     Local reverse proxy configuration
├── docs/            Development and troubleshooting notes
├── quokka_reference/ Learning path and architecture reference
└── docker-compose.yml
```

| Layer             | Stack                                           | Responsibility                                              |
| ----------------- | ----------------------------------------------- | ----------------------------------------------------------- |
| Frontend          | Vue 3, TypeScript, Web Audio API, MediaRecorder | Music input, synthesis, recording, visualization, upload UI |
| Primary API       | Ruby on Rails 8, PostgreSQL                     | Users, compositions, uploads, comments, likes, playlists    |
| Audio engine      | C++20, CMake, libsndfile                        | Offline WAV rendering with DSP effects                      |
| Queue             | Sidekiq, Redis                                  | Long-running audio and AI jobs                              |
| Realtime          | Action Cable                                    | Collaboration room events                                   |
| Proxy             | Nginx                                           | Unified local entry point                                   |
| Compatibility API | PHP 8                                           | Alternate API surface for experiments                       |

## Requirements

The easiest path is Docker Desktop with the Linux/WSL2 engine enabled.

For manual development, install:

- Node.js 22+
- Ruby 3.3+
- PostgreSQL 17+
- Redis 7+
- PHP 8.3+
- CMake 3.20+
- A C++20 compiler
- libsndfile
- ffmpeg

## Quick Start

From the project root:

```bash
docker compose up --build
```

Open:

```text
http://localhost:8088
```

Useful service URLs:

```text
Frontend dev server: http://localhost:5173
Rails API:           http://localhost:3000
Unified entry:       http://localhost:8088
PHP API:             http://localhost:8080
PostgreSQL:          localhost:5432
Redis:               localhost:6379
```

The first startup may take a while because Docker installs dependencies, builds the C++ audio engine, prepares the database, and seeds demo data.

Demo account:

```text
Email:    demo@quokka.local
Password: password123
```

## Running The App

1. Open `http://localhost:8088`.
2. Click or press mapped keys on the virtual keyboard.
3. Choose waveform, volume, visualizer mode, and effect parameters.
4. Click `Record`.
5. Play a short melody.
6. Click `Stop`.
7. Click `Upload` to save the original recording.
8. Click `C++ Reverb` to process the recording through the audio engine.
9. Use `Refresh` in the composition list to reload public works.

## Keyboard Mapping

White keys:

```text
A S D F G H J K L Z X C V B
```

Black keys:

```text
Q W E R T Y U I O P
```

## API Endpoints

```http
GET    /api/v1/health
GET    /api/v1/compositions
POST   /api/v1/compositions
GET    /api/v1/compositions/:id
PATCH  /api/v1/compositions/:id
DELETE /api/v1/compositions/:id
GET    /api/v1/compositions/:composition_id/comments
POST   /api/v1/compositions/:composition_id/comments
POST   /api/v1/compositions/:composition_id/like
DELETE /api/v1/compositions/:composition_id/like
POST   /api/v1/audio/process
POST   /api/v1/signup
POST   /api/v1/login
DELETE /api/v1/logout
POST   /api/v1/ai/chords
POST   /api/v1/ai/mix
POST   /api/v1/ai/style_transfer
POST   /api/v1/collaboration/:room_id/events
POST   /internal/liora_trigger
```

Health check:

```bash
curl http://localhost:8088/api/v1/health
```

Upload a composition:

```bash
curl -F "composition[title]=First melody" \
  -F "composition[bpm]=120" \
  -F "composition[key_signature]=C" \
  -F "composition[midi_data]={\"events\":[]}" \
  -F "composition[audio]=@recording.webm" \
  http://localhost:8088/api/v1/compositions
```

Process audio through the C++ engine:

```bash
curl -F "title=Processed melody" \
  -F "reverb=0.65" \
  -F "delay_ms=250" \
  -F "pitch_semitones=0" \
  -F "audio=@recording.webm" \
  http://localhost:8088/api/v1/audio/process
```

## Manual Development

Frontend:

```bash
cd frontend
npm install
npm run dev
```

Frontend verification:

```bash
npm run build
npm test
npm run lint
```

Rails API:

```bash
cd backend_ruby
bundle install
bin/rails db:prepare
bin/rails db:seed
bin/rails server
```

Sidekiq:

```bash
cd backend_ruby
bundle exec sidekiq
```

C++ audio engine:

```bash
cd audio_engine
cmake -S . -B build
cmake --build build
./build/quokka_audio input.wav output.wav --reverb 0.8 --delay-ms 250 --pitch-semitones 0
```

PHP compatibility API:

```bash
cd backend_php
php -S 0.0.0.0:8080 -t public
```

## Docker Commands

Start all services:

```bash
docker compose up --build
```

Start in the background:

```bash
docker compose up -d --build
```

View service status:

```bash
docker compose ps
```

View logs:

```bash
docker compose logs -f rails frontend nginx
```

Stop services:

```bash
docker compose down
```

Reset local database volumes:

```bash
docker compose down -v
docker compose up --build
```

## Testing Status

The frontend project has been verified with:

```bash
npm run build
npm test
npm run lint
```

The Rails and PHP code has been checked for syntax. Docker Compose has been validated and the main local endpoints have been tested through `http://localhost:8088`.

## Notes

- The C++ pitch-shift effect is an MVP implementation based on linear resampling. Replace it with a phase vocoder for production-quality pitch shifting.
- The Rails API allows guest creation, likes, and comments during MVP development.
- Active Storage uses local disk storage by default. S3 or MinIO can be added later.
- The AI endpoints currently return deterministic scaffold responses and are ready for future model integration.
- The Liora integration endpoint is reserved for internal automation and future sound-module triggers.
