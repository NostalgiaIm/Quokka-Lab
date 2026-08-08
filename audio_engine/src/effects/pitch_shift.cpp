#include "effects/pitch_shift.h"

#include <algorithm>
#include <cmath>
#include <vector>

namespace quokka::effects {

void ApplyNaivePitchShift(dsp::AudioBuffer& buffer, float semitones) {
  if (std::abs(semitones) < 0.001f || buffer.samples.empty()) {
    return;
  }

  // MVP 占位实现：这里先做线性重采样；生产质量应替换为 Phase Vocoder。
  const float ratio = std::pow(2.0f, semitones / 12.0f);
  std::vector<float> output(buffer.samples.size(), 0.0f);

  for (int frame = 0; frame < buffer.frames(); ++frame) {
    const float source_frame = frame * ratio;
    const int base = static_cast<int>(source_frame);
    const float fraction = source_frame - base;

    if (base + 1 >= buffer.frames()) {
      break;
    }

    for (int channel = 0; channel < buffer.channels; ++channel) {
      const size_t a = static_cast<size_t>(base * buffer.channels + channel);
      const size_t b = static_cast<size_t>((base + 1) * buffer.channels + channel);
      const size_t out = static_cast<size_t>(frame * buffer.channels + channel);
      output[out] = buffer.samples[a] * (1.0f - fraction) + buffer.samples[b] * fraction;
    }
  }

  buffer.samples = std::move(output);
}

}  // namespace quokka::effects
