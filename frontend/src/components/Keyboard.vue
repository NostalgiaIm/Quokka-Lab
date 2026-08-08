<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';

import type { NoteDefinition, Waveform } from '@/types/audio';

const props = defineProps<{
  waveform: Waveform;
}>();

const emit = defineEmits<{
  noteDown: [note: NoteDefinition];
  noteUp: [note: NoteDefinition];
}>();

const WHITE_KEYS = ['A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'Z', 'X', 'C', 'V', 'B'];
const BLACK_KEYS = ['Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'];
const BLACK_OFFSETS = [0.68, 1.68, 3.68, 4.68, 5.68, 7.68, 8.68, 10.68, 11.68, 12.68];

const activeNoteIds = ref(new Set<string>());

const notes = computed<NoteDefinition[]>(() => {
  const whiteNames = ['C4', 'D4', 'E4', 'F4', 'G4', 'A4', 'B4', 'C5', 'D5', 'E5', 'F5', 'G5', 'A5', 'B5'];
  const blackNames = ['C#4', 'D#4', 'F#4', 'G#4', 'A#4', 'C#5', 'D#5', 'F#5', 'G#5', 'A#5'];

  const whites = whiteNames.map((name, index) => ({
    id: name,
    name,
    frequency: frequencyForNote(name),
    keyboard: WHITE_KEYS[index],
    color: 'white' as const,
    octave: Number(name.at(-1)),
    whiteIndex: index,
  }));

  const blacks = blackNames.map((name, index) => ({
    id: name,
    name,
    frequency: frequencyForNote(name),
    keyboard: BLACK_KEYS[index],
    color: 'black' as const,
    octave: Number(name.at(-1)),
  }));

  return [...whites, ...blacks];
});

const whiteNotes = computed(() => notes.value.filter((note) => note.color === 'white'));
const blackNotes = computed(() => notes.value.filter((note) => note.color === 'black'));
const noteByKeyboard = computed(() => new Map(notes.value.map((note) => [note.keyboard.toLowerCase(), note])));

function frequencyForNote(note: string): number {
  const match = note.match(/^([A-G]#?)(\d)$/);
  if (!match) throw new Error(`Invalid note: ${note}`);

  const [, pitch, octave] = match;
  const semitoneFromC: Record<string, number> = {
    C: 0,
    'C#': 1,
    D: 2,
    'D#': 3,
    E: 4,
    F: 5,
    'F#': 6,
    G: 7,
    'G#': 8,
    A: 9,
    'A#': 10,
    B: 11,
  };

  const midiNumber = (Number(octave) + 1) * 12 + semitoneFromC[pitch];
  return 440 * 2 ** ((midiNumber - 69) / 12);
}

function pressNote(note: NoteDefinition): void {
  if (activeNoteIds.value.has(note.id)) return;

  activeNoteIds.value = new Set(activeNoteIds.value).add(note.id);
  emit('noteDown', note);
}

function releaseNote(note: NoteDefinition): void {
  if (!activeNoteIds.value.has(note.id)) return;

  const next = new Set(activeNoteIds.value);
  next.delete(note.id);
  activeNoteIds.value = next;
  emit('noteUp', note);
}

function handleKeyDown(event: KeyboardEvent): void {
  if (event.repeat) return;

  const note = noteByKeyboard.value.get(event.key.toLowerCase());
  if (!note) return;

  event.preventDefault();
  pressNote(note);
}

function handleKeyUp(event: KeyboardEvent): void {
  const note = noteByKeyboard.value.get(event.key.toLowerCase());
  if (!note) return;

  event.preventDefault();
  releaseNote(note);
}

function blackKeyLeft(index: number): string {
  return `calc(${BLACK_OFFSETS[index]} * var(--white-key-width))`;
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown);
  window.addEventListener('keyup', handleKeyUp);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown);
  window.removeEventListener('keyup', handleKeyUp);
});
</script>

<template>
  <section class="keyboard-shell" :data-waveform="props.waveform">
    <div class="keyboard" aria-label="Two octave piano keyboard">
      <button
        v-for="note in whiteNotes"
        :key="note.id"
        class="piano-key piano-key-white"
        :class="{ active: activeNoteIds.has(note.id) }"
        type="button"
        @mousedown="pressNote(note)"
        @mouseup="releaseNote(note)"
        @mouseleave="releaseNote(note)"
        @touchstart.prevent="pressNote(note)"
        @touchend.prevent="releaseNote(note)"
      >
        <span>{{ note.name }}</span>
        <kbd>{{ note.keyboard }}</kbd>
      </button>

      <button
        v-for="(note, index) in blackNotes"
        :key="note.id"
        class="piano-key piano-key-black"
        :class="{ active: activeNoteIds.has(note.id) }"
        :style="{ left: blackKeyLeft(index) }"
        type="button"
        @mousedown="pressNote(note)"
        @mouseup="releaseNote(note)"
        @mouseleave="releaseNote(note)"
        @touchstart.prevent="pressNote(note)"
        @touchend.prevent="releaseNote(note)"
      >
        <span>{{ note.name }}</span>
        <kbd>{{ note.keyboard }}</kbd>
      </button>
    </div>
  </section>
</template>

<style scoped>
.keyboard-shell {
  --white-key-width: min(6.3vw, 64px);
  --white-key-height: min(38vw, 240px);
  inline-size: fit-content;
  max-inline-size: 100%;
  margin: 0 auto;
  overflow-x: auto;
  padding: 20px 4px 8px;
}

.keyboard {
  position: relative;
  display: grid;
  grid-template-columns: repeat(14, var(--white-key-width));
  inline-size: calc(14 * var(--white-key-width));
  min-inline-size: 600px;
  block-size: var(--white-key-height);
}

.piano-key {
  border: 0;
  cursor: pointer;
  font: inherit;
  user-select: none;
  touch-action: manipulation;
}

.piano-key-white {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
  block-size: var(--white-key-height);
  padding: 0 4px 14px;
  border: 1px solid #cfd5df;
  border-radius: 0 0 6px 6px;
  background: linear-gradient(180deg, #fff 0%, #edf1f5 100%);
  color: #1b2636;
  box-shadow: inset 0 -8px 16px rgb(27 38 54 / 8%);
}

.piano-key-black {
  position: absolute;
  z-index: 2;
  inset-block-start: 0;
  inline-size: calc(var(--white-key-width) * 0.64);
  block-size: calc(var(--white-key-height) * 0.62);
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
  padding: 0 2px 12px;
  border-radius: 0 0 5px 5px;
  background: linear-gradient(180deg, #1b2531 0%, #070a0f 100%);
  color: #f8fafc;
  box-shadow: 0 10px 20px rgb(0 0 0 / 28%);
}

.piano-key.active {
  background: #2db4a8;
  color: #041414;
  box-shadow: inset 0 0 0 2px rgb(255 255 255 / 40%);
}

.piano-key span {
  font-size: 0.82rem;
  font-weight: 700;
}

kbd {
  min-inline-size: 28px;
  padding: 3px 6px;
  border-radius: 4px;
  background: rgb(255 255 255 / 76%);
  color: #172033;
  font-size: 0.72rem;
  font-weight: 700;
}

.piano-key-black kbd {
  background: rgb(255 255 255 / 18%);
  color: #f8fafc;
}

@media (max-width: 760px) {
  .keyboard-shell {
    --white-key-width: 48px;
    --white-key-height: 190px;
    inline-size: 100%;
  }
}
</style>
