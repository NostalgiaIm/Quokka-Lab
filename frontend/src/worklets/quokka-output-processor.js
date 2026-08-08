class QuokkaOutputProcessor extends AudioWorkletProcessor {
  /**
   * 在最终输出阶段做一个轻微的软限幅。
   * 这样多个音符叠加时不容易产生刺耳的数字削波。
   */
  process(inputs, outputs) {
    const input = inputs[0];
    const output = outputs[0];

    for (let channel = 0; channel < output.length; channel += 1) {
      const source = input[channel] || input[0];
      const target = output[channel];

      for (let index = 0; index < target.length; index += 1) {
        const sample = source ? source[index] : 0;
        target[index] = Math.tanh(sample * 1.1);
      }
    }

    return true;
  }
}

registerProcessor('quokka-output-processor', QuokkaOutputProcessor);
