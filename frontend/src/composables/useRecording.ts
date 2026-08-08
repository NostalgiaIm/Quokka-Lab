import { ref, type Ref } from 'vue';

export interface RecorderState {
  audioUrl: Ref<string | null>;
  isPaused: Ref<boolean>;
  isRecording: Ref<boolean>;
  recordedBlob: Ref<Blob | null>;
  downloadRecording: () => void;
  pauseRecording: () => void;
  resetRecording: () => void;
  resumeRecording: () => void;
  startRecording: (stream: MediaStream) => void;
  stopRecording: () => Promise<Blob | null>;
}

const audioUrl = ref<string | null>(null);
const isPaused = ref(false);
const isRecording = ref(false);
const recordedBlob = ref<Blob | null>(null);

let recorder: MediaRecorder | null = null;
let chunks: BlobPart[] = [];
let stopResolver: ((blob: Blob | null) => void) | null = null;

/**
 * 封装 Web Audio 输出流上的 MediaRecorder。
 * 录制结束后会生成 Blob URL，供页面立即播放或下载。
 */
export function useRecording(): RecorderState {
  function resetRecording(): void {
    if (audioUrl.value) {
      URL.revokeObjectURL(audioUrl.value);
    }

    audioUrl.value = null;
    recordedBlob.value = null;
    chunks = [];
  }

  function startRecording(stream: MediaStream): void {
    resetRecording();
    const mimeType = pickMimeType();
    recorder = new MediaRecorder(stream, mimeType ? { mimeType } : undefined);

    recorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        chunks.push(event.data);
      }
    };

    recorder.onstop = () => {
      const blob = chunks.length > 0 ? new Blob(chunks, { type: recorder?.mimeType || 'audio/webm' }) : null;
      recordedBlob.value = blob;
      audioUrl.value = blob ? URL.createObjectURL(blob) : null;
      isRecording.value = false;
      isPaused.value = false;
      stopResolver?.(blob);
      stopResolver = null;
    };

    recorder.start();
    isRecording.value = true;
  }

  function pauseRecording(): void {
    if (recorder?.state === 'recording') {
      recorder.pause();
      isPaused.value = true;
    }
  }

  function resumeRecording(): void {
    if (recorder?.state === 'paused') {
      recorder.resume();
      isPaused.value = false;
    }
  }

  function stopRecording(): Promise<Blob | null> {
    if (!recorder || recorder.state === 'inactive') {
      return Promise.resolve(recordedBlob.value);
    }

    return new Promise((resolve) => {
      stopResolver = resolve;
      recorder?.stop();
    });
  }

  function downloadRecording(): void {
    if (!audioUrl.value) return;

    const link = document.createElement('a');
    link.href = audioUrl.value;
    link.download = `quokka-recording-${Date.now()}.webm`;
    link.click();
  }

  return {
    audioUrl,
    isPaused,
    isRecording,
    recordedBlob,
    downloadRecording,
    pauseRecording,
    resetRecording,
    resumeRecording,
    startRecording,
    stopRecording,
  };
}

function pickMimeType(): string {
  const supportedTypes = ['audio/webm;codecs=opus', 'audio/webm', 'audio/ogg;codecs=opus'];
  return supportedTypes.find((type) => MediaRecorder.isTypeSupported(type)) ?? '';
}
