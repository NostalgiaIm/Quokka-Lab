import type { NoteEvent, Waveform } from '@/types/audio';

interface ActiveNote {
  oscillator: OscillatorNode;
  gain: GainNode;
}

export interface AudioEngine {
  analyser: AnalyserNode;
  destination: MediaStreamAudioDestinationNode;
  events: NoteEvent[];
  masterGain: GainNode;
  playNote: (frequency: number, duration?: number, waveform?: Waveform) => void;
  startNote: (note: string, frequency: number, waveform?: Waveform) => void;
  stopNote: (note: string) => void;
  setVolume: (volume: number) => void;
  resume: () => Promise<void>;
}

const ATTACK_SECONDS = 0.012;
const RELEASE_SECONDS = 0.08;

/**
 * 创建 Quokka Studio 在浏览器端使用的合成器音频图。
 *
 * 音频链路：
 * OscillatorNode -> GainNode -> master GainNode -> AnalyserNode
 *                                          |-> 扬声器播放
 *                                          |-> MediaRecorder 录制流
 */
export function useAudio(): AudioEngine {
  const AudioContextConstructor = window.AudioContext ?? window.webkitAudioContext;
  if (!AudioContextConstructor) {
    throw new Error('Web Audio API is not available in this browser.');
  }

  const context = new AudioContextConstructor();
  const masterGain = context.createGain();
  const analyser = context.createAnalyser();
  const destination = context.createMediaStreamDestination();
  const activeNotes = new Map<string, ActiveNote>();
  const events: NoteEvent[] = [];
  let workletReady = false;
  let workletPromise: Promise<void> | null = null;

  masterGain.gain.value = 0.6;
  analyser.fftSize = 2048;

  masterGain.connect(analyser);
  analyser.connect(context.destination);
  analyser.connect(destination);

  async function resume(): Promise<void> {
    await ensureOutputWorklet();

    if (context.state === 'suspended') {
      await context.resume();
    }
  }

  async function ensureOutputWorklet(): Promise<void> {
    if (workletReady || workletPromise || !context.audioWorklet) {
      await workletPromise;
      return;
    }

    workletPromise = context.audioWorklet
      .addModule(new URL('../worklets/quokka-output-processor.js', import.meta.url).href)
      .then(() => {
        const outputWorklet = new AudioWorkletNode(context, 'quokka-output-processor');
        masterGain.disconnect();
        masterGain.connect(outputWorklet);
        outputWorklet.connect(analyser);
        workletReady = true;
      })
      .catch((error) => {
        console.warn('AudioWorklet 不可用，已回退到直接 Web Audio 链路。', error);
      });

    await workletPromise;
  }

  function setVolume(volume: number): void {
    const safeVolume = Math.min(Math.max(volume, 0), 1);
    masterGain.gain.setTargetAtTime(safeVolume, context.currentTime, 0.01);
  }

  function createVoice(frequency: number, waveform: Waveform = 'sine'): ActiveNote {
    const oscillator = context.createOscillator();
    const gain = context.createGain();
    const now = context.currentTime;

    oscillator.type = waveform;
    oscillator.frequency.setValueAtTime(frequency, now);
    gain.gain.setValueAtTime(0, now);
    gain.gain.linearRampToValueAtTime(1, now + ATTACK_SECONDS);

    oscillator.connect(gain);
    gain.connect(masterGain);
    oscillator.start(now);

    return { oscillator, gain };
  }

  function releaseVoice(voice: ActiveNote): void {
    const now = context.currentTime;

    voice.gain.gain.cancelScheduledValues(now);
    voice.gain.gain.setValueAtTime(voice.gain.gain.value, now);
    voice.gain.gain.linearRampToValueAtTime(0, now + RELEASE_SECONDS);
    voice.oscillator.stop(now + RELEASE_SECONDS + 0.01);
  }

  function startNote(note: string, frequency: number, waveform: Waveform = 'sine'): void {
    if (activeNotes.has(note)) return;

    const voice = createVoice(frequency, waveform);
    activeNotes.set(note, voice);
    events.push({ type: 'note_on', note, frequency, at: performance.now() });
  }

  function stopNote(note: string): void {
    const voice = activeNotes.get(note);
    if (!voice) return;

    releaseVoice(voice);
    activeNotes.delete(note);
    events.push({ type: 'note_off', note, frequency: voice.oscillator.frequency.value, at: performance.now() });
  }

  function playNote(frequency: number, duration = 0.45, waveform: Waveform = 'sine'): void {
    const noteId = `preview:${frequency}:${performance.now()}`;
    startNote(noteId, frequency, waveform);
    window.setTimeout(() => stopNote(noteId), duration * 1000);
  }

  return {
    analyser,
    destination,
    events,
    masterGain,
    playNote,
    startNote,
    stopNote,
    setVolume,
    resume,
  };
}

declare global {
  interface Window {
    webkitAudioContext?: typeof AudioContext;
  }
}
