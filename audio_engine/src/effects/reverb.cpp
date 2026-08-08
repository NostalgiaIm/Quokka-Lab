#include "effects/reverb.h"

#include <algorithm>
#include <array>
#include <cmath>
#include <vector>

namespace quokka::effects {
namespace {

class CombFilter {
 public:
  CombFilter(int delay_samples, float feedback)
      : delay_(std::max(delay_samples, 1)), feedback_(feedback), memory_(static_cast<size_t>(delay_), 0.0f) {}

  float Process(float input) {
    const float delayed = memory_[index_];
    memory_[index_] = input + delayed * feedback_;
    index_ = (index_ + 1) % memory_.size();
    return delayed;
  }

 private:
  int delay_;
  float feedback_;
  std::vector<float> memory_;
  size_t index_ = 0;
};

float ClampSample(float value) {
  return std::clamp(value, -1.0f, 1.0f);
}

}  // namespace

void ApplySchroederReverb(dsp::AudioBuffer& buffer, float mix) {
  mix = std::clamp(mix, 0.0f, 1.0f);
  if (mix <= 0.0f || buffer.samples.empty()) {
    return;
  }

  // Schroeder 混响常用若干接近质数关系的延迟长度，以减少明显的金属振铃。
  const std::array<float, 4> delay_ms = {29.7f, 37.1f, 41.1f, 43.7f};
  const std::array<float, 4> feedback = {0.742f, 0.733f, 0.715f, 0.697f};
  std::vector<std::array<CombFilter, 4>> filters;
  filters.reserve(static_cast<size_t>(buffer.channels));

  for (int channel = 0; channel < buffer.channels; ++channel) {
    filters.push_back({
        CombFilter(static_cast<int>(buffer.sample_rate * delay_ms[0] / 1000.0f) + channel * 23, feedback[0]),
        CombFilter(static_cast<int>(buffer.sample_rate * delay_ms[1] / 1000.0f) + channel * 29, feedback[1]),
        CombFilter(static_cast<int>(buffer.sample_rate * delay_ms[2] / 1000.0f) + channel * 31, feedback[2]),
        CombFilter(static_cast<int>(buffer.sample_rate * delay_ms[3] / 1000.0f) + channel * 37, feedback[3]),
    });
  }

  for (int frame = 0; frame < buffer.frames(); ++frame) {
    for (int channel = 0; channel < buffer.channels; ++channel) {
      const size_t sample_index = static_cast<size_t>(frame * buffer.channels + channel);
      const float dry = buffer.samples[sample_index];
      float wet = 0.0f;

      for (auto& filter : filters[static_cast<size_t>(channel)]) {
        wet += filter.Process(dry);
      }
      wet *= 0.25f;

      buffer.samples[sample_index] = ClampSample(dry * (1.0f - mix) + wet * mix);
    }
  }
}

}  // namespace quokka::effects
