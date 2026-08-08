#include "dsp/wav_io.h"

#include <sndfile.h>

#include <stdexcept>
#include <string>

namespace quokka::dsp {

AudioBuffer ReadWav(const std::filesystem::path& path) {
  SF_INFO info{};
  SNDFILE* file = sf_open(path.string().c_str(), SFM_READ, &info);
  if (file == nullptr) {
    throw std::runtime_error("Failed to open input WAV: " + path.string());
  }

  AudioBuffer buffer;
  buffer.sample_rate = info.samplerate;
  buffer.channels = info.channels;
  buffer.samples.resize(static_cast<size_t>(info.frames * info.channels));

  const sf_count_t read_count = sf_readf_float(file, buffer.samples.data(), info.frames);
  sf_close(file);

  if (read_count != info.frames) {
    throw std::runtime_error("Failed to read all WAV frames: " + path.string());
  }

  return buffer;
}

void WriteWav(const std::filesystem::path& path, const AudioBuffer& buffer) {
  SF_INFO info{};
  info.channels = buffer.channels;
  info.samplerate = buffer.sample_rate;
  info.format = SF_FORMAT_WAV | SF_FORMAT_PCM_16;

  SNDFILE* file = sf_open(path.string().c_str(), SFM_WRITE, &info);
  if (file == nullptr) {
    throw std::runtime_error("Failed to open output WAV: " + path.string());
  }

  const sf_count_t frames = static_cast<sf_count_t>(buffer.frames());
  const sf_count_t written = sf_writef_float(file, buffer.samples.data(), frames);
  sf_close(file);

  if (written != frames) {
    throw std::runtime_error("Failed to write all WAV frames: " + path.string());
  }
}

}  // namespace quokka::dsp
