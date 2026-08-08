#pragma once

#include "dsp/audio_buffer.h"

namespace quokka::effects {

void ApplyFeedbackDelay(dsp::AudioBuffer& buffer, float delay_ms, float feedback, float mix);

}  // namespace quokka::effects
