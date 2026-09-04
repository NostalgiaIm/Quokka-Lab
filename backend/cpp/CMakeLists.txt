cmake_minimum_required(VERSION 3.20)
project(quokka_audio_engine VERSION 0.1.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

find_package(PkgConfig REQUIRED)
pkg_check_modules(SNDFILE REQUIRED IMPORTED_TARGET sndfile)

add_executable(quokka_audio
  src/main.cpp
  src/effects/reverb.cpp
  src/effects/delay.cpp
  src/effects/pitch_shift.cpp
  src/dsp/wav_io.cpp
)

target_include_directories(quokka_audio PRIVATE src)
target_link_libraries(quokka_audio PRIVATE PkgConfig::SNDFILE)

if(MSVC)
  target_compile_options(quokka_audio PRIVATE /W4)
else()
  target_compile_options(quokka_audio PRIVATE -Wall -Wextra -Wpedantic)
endif()
