import { defineStore } from 'pinia';

import type { Composition, CreateCompositionPayload } from '@/types/composition';

const API_BASE = '/api/v1';

export const useCompositionsStore = defineStore('compositions', {
  state: () => ({
    compositions: [] as Composition[],
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async fetchCompositions() {
      this.loading = true;
      this.error = null;

      try {
        const response = await fetch(`${API_BASE}/compositions`);
        if (!response.ok) throw new Error(`Fetch failed with ${response.status}`);
        this.compositions = await response.json();
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unknown fetch error';
      } finally {
        this.loading = false;
      }
    },

    async createComposition(payload: CreateCompositionPayload) {
      const formData = new FormData();
      formData.append('composition[title]', payload.title);
      formData.append('composition[bpm]', String(payload.bpm ?? 120));
      formData.append('composition[key_signature]', payload.key_signature ?? 'C');
      formData.append('composition[is_public]', String(payload.is_public ?? true));
      formData.append('composition[midi_data]', JSON.stringify(payload.midi_data ?? {}));

      if (payload.audio) {
        formData.append('composition[audio]', payload.audio, 'recording.webm');
      }

      const response = await fetch(`${API_BASE}/compositions`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) throw new Error(`Create failed with ${response.status}`);
      const composition = (await response.json()) as Composition;
      this.compositions.unshift(composition);
      return composition;
    },
  },
});
