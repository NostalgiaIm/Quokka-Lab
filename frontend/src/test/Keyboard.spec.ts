import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';

import Keyboard from '@/components/Keyboard.vue';

describe('Keyboard', () => {
  it('renders two octaves from C4 to B5', () => {
    const wrapper = mount(Keyboard, {
      props: {
        waveform: 'sine',
      },
    });

    expect(wrapper.text()).toContain('C4');
    expect(wrapper.text()).toContain('B5');
    expect(wrapper.findAll('.piano-key-white')).toHaveLength(14);
    expect(wrapper.findAll('.piano-key-black')).toHaveLength(10);
  });
});
