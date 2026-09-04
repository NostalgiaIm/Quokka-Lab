#pragma once

#include "dsp/audio_buffer.h"

namespace quokka::effects {

// 使用四个并联梳状滤波器实现一个轻量 Schroeder 风格混响。
void ApplySchroederReverb(dsp::AudioBuffer& buffer, float mix);

}  // namespace quokka::effects
