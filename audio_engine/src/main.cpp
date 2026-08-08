#include <cstdlib>
#include <exception>
#include <filesystem>
#include <iostream>
#include <stdexcept>
#include <string>
#include <string_view>

#include "dsp/wav_io.h"
#include "effects/delay.h"
#include "effects/pitch_shift.h"
#include "effects/reverb.h"

namespace {

struct Options {
  std::filesystem::path input;
  std::filesystem::path output;
  float reverb = 0.0f;
  float delay_ms = 0.0f;
  float delay_feedback = 0.35f;
  float pitch_semitones = 0.0f;
};

void PrintUsage() {
  std::cerr << "Usage: quokka_audio input.wav output.wav [--reverb 0..1] "
               "[--delay-ms N] [--delay-feedback 0..0.95] [--pitch-semitones N]\n";
}

float ParseFloat(const char* value, std::string_view flag) {
  char* end = nullptr;
  const float parsed = std::strtof(value, &end);
  if (end == value || *end != '\0') {
    throw std::invalid_argument("Invalid numeric value for " + std::string(flag));
  }
  return parsed;
}

Options ParseArgs(int argc, char** argv) {
  if (argc < 3) {
    PrintUsage();
    throw std::invalid_argument("Missing input or output path");
  }

  Options options;
  options.input = argv[1];
  options.output = argv[2];

  for (int index = 3; index < argc; ++index) {
    const std::string_view flag = argv[index];
    if (index + 1 >= argc) {
      throw std::invalid_argument("Missing value for " + std::string(flag));
    }

    if (flag == "--reverb") {
      options.reverb = ParseFloat(argv[++index], flag);
    } else if (flag == "--delay-ms") {
      options.delay_ms = ParseFloat(argv[++index], flag);
    } else if (flag == "--delay-feedback") {
      options.delay_feedback = ParseFloat(argv[++index], flag);
    } else if (flag == "--pitch-semitones") {
      options.pitch_semitones = ParseFloat(argv[++index], flag);
    } else {
      throw std::invalid_argument("Unknown flag: " + std::string(flag));
    }
  }

  return options;
}

}  // namespace

int main(int argc, char** argv) {
  try {
    const Options options = ParseArgs(argc, argv);
    quokka::dsp::AudioBuffer buffer = quokka::dsp::ReadWav(options.input);

    quokka::effects::ApplySchroederReverb(buffer, options.reverb);
    const float delay_mix = options.delay_ms > 0.0f ? 0.35f : 0.0f;
    quokka::effects::ApplyFeedbackDelay(buffer, options.delay_ms, options.delay_feedback, delay_mix);
    quokka::effects::ApplyNaivePitchShift(buffer, options.pitch_semitones);

    quokka::dsp::WriteWav(options.output, buffer);
    std::cout << "Processed " << options.input << " -> " << options.output << '\n';
    return 0;
  } catch (const std::exception& error) {
    std::cerr << "quokka_audio error: " << error.what() << '\n';
    return 1;
  }
}
