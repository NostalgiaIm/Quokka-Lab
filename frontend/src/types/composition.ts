export interface Composition {
  id: number;
  title: string;
  audio_url: string | null;
  midi_data: Record<string, unknown>;
  bpm: number;
  key_signature: string;
  is_public: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateCompositionPayload {
  title: string;
  bpm?: number;
  key_signature?: string;
  is_public?: boolean;
  midi_data?: Record<string, unknown>;
  audio?: Blob;
}
