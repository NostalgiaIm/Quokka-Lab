#pragma once

#include <filesystem>

#include "dsp/audio_buffer.h"

namespace quokka::dsp {

AudioBuffer ReadWav(const std::filesystem::path& path);
void WriteWav(const std::filesystem::path& path, const AudioBuffer& buffer);

}  // namespace quokka::dsp
