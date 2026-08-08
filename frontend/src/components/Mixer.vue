<script setup lang="ts">
import type { Waveform } from '@/types/audio';

const waveform = defineModel<Waveform>('waveform', { required: true });
const volume = defineModel<number>('volume', { required: true });
const visualizerMode = defineModel<'waveform' | 'spectrum'>('visualizerMode', { required: true });
</script>

<template>
  <section class="mixer" aria-label="Synth controls">
    <label>
      Waveform
      <select v-model="waveform">
        <option value="sine">Sine</option>
        <option value="square">Square</option>
        <option value="sawtooth">Saw</option>
        <option value="triangle">Triangle</option>
      </select>
    </label>

    <label>
      Volume
      <input v-model.number="volume" min="0" max="1" step="0.01" type="range" />
    </label>

    <label>
      Visualizer
      <select v-model="visualizerMode">
        <option value="waveform">Waveform</option>
        <option value="spectrum">Spectrum</option>
      </select>
    </label>
  </section>
</template>

<style scoped>
.mixer {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

label {
  display: grid;
  gap: 8px;
  color: #526173;
  font-size: 0.86rem;
  font-weight: 700;
}

select,
input {
  min-block-size: 38px;
  inline-size: 100%;
}

select {
  border: 1px solid #cfd8e3;
  border-radius: 6px;
  background: #fff;
  color: #162030;
  padding: 0 10px;
}

@media (max-width: 760px) {
  .mixer {
    grid-template-columns: 1fr;
  }
}
</style>
