#pragma once

#include <vector>

namespace quokka::dsp {

struct AudioBuffer {
  int sample_rate = 44100;
  int channels = 2;
  std::vector<float> samples;

  [[nodiscard]] int frames() const {
    if (channels <= 0) {
      return 0;
    }
    return static_cast<int>(samples.size() / static_cast<size_t>(channels));
  }
};

}  // namespace quokka::dsp
