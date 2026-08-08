export type Waveform = OscillatorType;

export interface NoteDefinition {
  id: string;
  name: string;
  frequency: number;
  keyboard: string;
  color: 'white' | 'black';
  octave: number;
  whiteIndex?: number;
}

export interface NoteEvent {
  type: 'note_on' | 'note_off';
  note: string;
  frequency: number;
  at: number;
}
