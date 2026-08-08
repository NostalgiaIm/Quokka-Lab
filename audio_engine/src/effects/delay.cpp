#include "effects/delay.h"

#include <algorithm>
#include <vector>

namespace quokka::effects {

void ApplyFeedbackDelay(dsp::AudioBuffer& buffer, float delay_ms, float feedback, float mix) {
  mix = std::clamp(mix, 0.0f, 1.0f);
  feedback = std::clamp(feedback, 0.0f, 0.95f);
  if (delay_ms <= 0.0f || mix <= 0.0f) {
    return;
  }

  const int delay_frames = std::max(1, static_cast<int>(buffer.sample_rate * delay_ms / 1000.0f));
  std::vector<float> memory(static_cast<size_t>(delay_frames * buffer.channels), 0.0f);
  size_t index = 0;

  for (float& sample : buffer.samples) {
    const float delayed = memory[index];
    memory[index] = sample + delayed * feedback;
    sample = std::clamp(sample * (1.0f - mix) + delayed * mix, -1.0f, 1.0f);
    index = (index + 1) % memory.size();
  }
}

}  // namespace quokka::effects
