#pragma once

#include "dsp/audio_buffer.h"

namespace quokka::effects {

void ApplyNaivePitchShift(dsp::AudioBuffer& buffer, float semitones);

}  // namespace quokka::effects
