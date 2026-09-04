<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue';

import Keyboard from '@/components/Keyboard.vue';
import Mixer from '@/components/Mixer.vue';
import Playlist from '@/components/Playlist.vue';
import Recorder from '@/components/Recorder.vue';
import Visualizer from '@/components/Visualizer.vue';
import { useAudio } from '@/composables/useAudio';
import { useRecording } from '@/composables/useRecording';
import { useWebSocket, type CollaborationSocket } from '@/composables/useWebSocket';
import { useCompositionsStore } from '@/stores/compositions';
import type { NoteDefinition, Waveform } from '@/types/audio';

const audio = useAudio();
const recorder = useRecording();
const compositionsStore = useCompositionsStore();

const waveform = ref<Waveform>('sine');
const volume = ref(0.6);
const visualizerMode = ref<'waveform' | 'spectrum'>('waveform');
const title = ref('Untitled melody');
const uploadStatus = ref('');
let collaborationSocket: CollaborationSocket | null = null;

watch(volume, (nextVolume) => audio.setVolume(nextVolume), { immediate: true });

onMounted(() => {
  void compositionsStore.fetchCompositions();
  collaborationSocket = useWebSocket('default-room', (event) => {
    if (event.type === 'note_on') {
      audio.startNote(`remote:${event.note}`, event.frequency, waveform.value);
    } else {
      audio.stopNote(`remote:${event.note}`);
    }
  });
});

onUnmounted(() => {
  collaborationSocket?.close();
});

async function noteDown(note: NoteDefinition): Promise<void> {
  await audio.resume();
  audio.startNote(note.id, note.frequency, waveform.value);
  collaborationSocket?.sendNoteEvent({ type: 'note_on', note: note.id, frequency: note.frequency, at: performance.now() });
}

function noteUp(note: NoteDefinition): void {
  audio.stopNote(note.id);
  collaborationSocket?.sendNoteEvent({ type: 'note_off', note: note.id, frequency: note.frequency, at: performance.now() });
}

async function startRecording(): Promise<void> {
  await audio.resume();
  recorder.startRecording(audio.destination.stream);
}

async function stopRecording(): Promise<void> {
  await recorder.stopRecording();
}

async function uploadRecording(): Promise<void> {
  if (!recorder.recordedBlob.value) return;

  uploadStatus.value = 'Uploading...';
  try {
    await compositionsStore.createComposition({
      title: title.value.trim() || 'Untitled melody',
      audio: recorder.recordedBlob.value,
      midi_data: { events: audio.events.slice(-1000) },
    });
    uploadStatus.value = 'Uploaded.';
  } catch (error) {
    uploadStatus.value = error instanceof Error ? error.message : 'Upload failed.';
  }
}
</script>

<template>
  <main class="studio">
    <header class="studio-header">
      <div>
        <p>Quokka Studio</p>
        <h1>Create, record, and share a melody.</h1>
      </div>
      <input v-model="title" aria-label="Composition title" />
    </header>

    <Mixer v-model:waveform="waveform" v-model:volume="volume" v-model:visualizer-mode="visualizerMode" />
    <Keyboard :waveform="waveform" @note-down="noteDown" @note-up="noteUp" />
    <Visualizer :analyser="audio.analyser" :mode="visualizerMode" />

    <Recorder
      :audio-url="recorder.audioUrl.value"
      :is-paused="recorder.isPaused.value"
      :is-recording="recorder.isRecording.value"
      @download="recorder.downloadRecording"
      @pause="recorder.pauseRecording"
      @resume="recorder.resumeRecording"
      @start="startRecording"
      @stop="stopRecording"
      @upload="uploadRecording"
    />

    <p v-if="uploadStatus" class="status">{{ uploadStatus }}</p>
    <Playlist
      :compositions="compositionsStore.compositions"
      :error="compositionsStore.error"
      :loading="compositionsStore.loading"
    />
  </main>
</template>

<style scoped>
.studio {
  display: grid;
  gap: 22px;
  inline-size: min(1120px, calc(100vw - 32px));
  margin: 0 auto;
  padding: 32px 0 52px;
}

.studio-header {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(220px, 360px);
  gap: 20px;
  align-items: end;
}

.studio-header p,
.studio-header h1 {
  margin: 0;
}

.studio-header p {
  color: #16736a;
  font-weight: 800;
  text-transform: uppercase;
}

.studio-header h1 {
  max-inline-size: 720px;
  font-size: clamp(2.1rem, 6vw, 4.4rem);
  line-height: 0.96;
}

.studio-header input {
  min-block-size: 42px;
  border: 1px solid #cfd8e3;
  border-radius: 6px;
  padding: 0 12px;
  font: inherit;
}

.status {
  margin: 0;
  color: #16736a;
  font-weight: 800;
}

@media (max-width: 760px) {
  .studio-header {
    grid-template-columns: 1fr;
  }
}
</style>
