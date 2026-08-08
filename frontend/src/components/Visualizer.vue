<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';

const props = defineProps<{
  analyser: AnalyserNode | null;
  mode: 'waveform' | 'spectrum';
}>();

const canvas = ref<HTMLCanvasElement | null>(null);
let animationFrame = 0;

function draw(): void {
  const target = canvas.value;
  const analyser = props.analyser;
  if (!target || !analyser) return;

  const context = target.getContext('2d');
  if (!context) return;

  const rect = target.getBoundingClientRect();
  const pixelRatio = window.devicePixelRatio || 1;
  target.width = Math.floor(rect.width * pixelRatio);
  target.height = Math.floor(rect.height * pixelRatio);
  context.scale(pixelRatio, pixelRatio);

  context.clearRect(0, 0, rect.width, rect.height);
  context.fillStyle = '#f6f8fb';
  context.fillRect(0, 0, rect.width, rect.height);

  if (props.mode === 'waveform') {
    drawWaveform(context, analyser, rect.width, rect.height);
  } else {
    drawSpectrum(context, analyser, rect.width, rect.height);
  }

  animationFrame = requestAnimationFrame(draw);
}

function drawWaveform(context: CanvasRenderingContext2D, analyser: AnalyserNode, width: number, height: number): void {
  const data = new Uint8Array(analyser.fftSize);
  analyser.getByteTimeDomainData(data);

  context.lineWidth = 2;
  context.strokeStyle = '#1a9187';
  context.beginPath();

  data.forEach((value, index) => {
    const x = (index / (data.length - 1)) * width;
    const y = (value / 255) * height;
    if (index === 0) {
      context.moveTo(x, y);
    } else {
      context.lineTo(x, y);
    }
  });

  context.stroke();
}

function drawSpectrum(context: CanvasRenderingContext2D, analyser: AnalyserNode, width: number, height: number): void {
  const data = new Uint8Array(analyser.frequencyBinCount);
  analyser.getByteFrequencyData(data);

  const barWidth = width / data.length;
  data.forEach((value, index) => {
    const normalized = value / 255;
    const barHeight = normalized * height;
    context.fillStyle = `rgb(${32 + normalized * 30}, ${110 + normalized * 110}, ${156 + normalized * 60})`;
    context.fillRect(index * barWidth, height - barHeight, Math.max(barWidth - 1, 1), barHeight);
  });
}

watch(() => props.mode, () => {
  cancelAnimationFrame(animationFrame);
  draw();
});

onMounted(draw);

onBeforeUnmount(() => {
  cancelAnimationFrame(animationFrame);
});
</script>

<template>
  <canvas ref="canvas" class="visualizer" aria-label="Audio visualizer"></canvas>
</template>

<style scoped>
.visualizer {
  display: block;
  inline-size: 100%;
  block-size: 180px;
  border: 1px solid #d7dee8;
  border-radius: 6px;
}
</style>
