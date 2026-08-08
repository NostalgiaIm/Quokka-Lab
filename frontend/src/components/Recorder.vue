<script setup lang="ts">
const props = defineProps<{
  audioUrl: string | null;
  isPaused: boolean;
  isRecording: boolean;
}>();

const emit = defineEmits<{
  download: [];
  pause: [];
  resume: [];
  start: [];
  stop: [];
  upload: [];
}>();
</script>

<template>
  <section class="recorder" aria-label="Recording controls">
    <div class="transport">
      <button v-if="!props.isRecording" class="record-button" type="button" @click="emit('start')">
        Record
      </button>
      <button v-else class="stop-button" type="button" @click="emit('stop')">Stop</button>

      <button v-if="props.isRecording && !props.isPaused" type="button" @click="emit('pause')">Pause</button>
      <button v-if="props.isRecording && props.isPaused" type="button" @click="emit('resume')">Resume</button>
      <button type="button" :disabled="!props.audioUrl" @click="emit('download')">Download</button>
      <button type="button" :disabled="!props.audioUrl" @click="emit('upload')">Upload</button>
    </div>

    <audio v-if="props.audioUrl" class="preview" :src="props.audioUrl" controls />
  </section>
</template>

<style scoped>
.recorder {
  display: grid;
  gap: 14px;
}

.transport {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

button {
  min-block-size: 40px;
  padding: 0 16px;
  border: 1px solid #cfd8e3;
  border-radius: 6px;
  background: #fff;
  color: #162030;
  font-weight: 700;
  cursor: pointer;
}

button:disabled {
  cursor: not-allowed;
  opacity: 0.45;
}

.record-button {
  border-color: #e65252;
  background: #e65252;
  color: #fff;
}

.stop-button {
  border-color: #15202f;
  background: #15202f;
  color: #fff;
}

.preview {
  inline-size: 100%;
}
</style>
