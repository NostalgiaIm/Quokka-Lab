<script setup lang="ts">
import type { Composition } from '@/types/composition';

const props = defineProps<{
  compositions: Composition[];
  loading: boolean;
  error: string | null;
}>();
</script>

<template>
  <section class="playlist" aria-label="Public compositions">
    <div class="playlist-header">
      <h2>Compositions</h2>
      <span>{{ props.compositions.length }}</span>
    </div>

    <p v-if="props.loading" class="muted">Loading compositions...</p>
    <p v-else-if="props.error" class="error">{{ props.error }}</p>
    <p v-else-if="props.compositions.length === 0" class="muted">No compositions yet.</p>

    <article v-for="composition in props.compositions" :key="composition.id" class="composition-card">
      <div>
        <h3>{{ composition.title }}</h3>
        <p>{{ composition.key_signature }} · {{ composition.bpm }} BPM</p>
      </div>
      <audio v-if="composition.audio_url" :src="composition.audio_url" controls />
    </article>
  </section>
</template>

<style scoped>
.playlist {
  display: grid;
  gap: 12px;
}

.playlist-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

h2,
h3,
p {
  margin: 0;
}

.playlist-header h2 {
  font-size: 1.1rem;
}

.playlist-header span {
  min-inline-size: 30px;
  padding: 4px 8px;
  border-radius: 999px;
  background: #e6f5f3;
  color: #12675f;
  text-align: center;
  font-weight: 800;
}

.composition-card {
  display: grid;
  gap: 10px;
  padding: 14px;
  border: 1px solid #d7dee8;
  border-radius: 6px;
  background: #fff;
}

.composition-card p,
.muted {
  color: #647184;
}

.error {
  color: #b42318;
}

audio {
  inline-size: 100%;
}
</style>
