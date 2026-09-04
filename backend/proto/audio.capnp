@0xdbc0333fc70f1a11;

struct AudioJob {
  inputPath @0 :Text;
  outputPath @1 :Text;
  reverb @2 :Float32;
  delayMs @3 :Float32;
  pitchSemitones @4 :Float32;
}

struct AudioJobResult {
  ok @0 :Bool;
  outputPath @1 :Text;
  error @2 :Text;
}
