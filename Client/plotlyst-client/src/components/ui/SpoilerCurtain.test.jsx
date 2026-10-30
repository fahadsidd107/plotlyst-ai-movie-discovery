import { fireEvent, render, screen } from '@testing-library/react';
import { expect, test, vi } from 'vitest';
import SpoilerCurtain from './SpoilerCurtain';

test('conceals spoiler content and supports an accessible deliberate reveal', () => {
  const reveal = vi.fn();
  render(<SpoilerCurtain label="Episode description" reason="Beyond your progress" onReveal={reveal} />);
  expect(screen.getByText('Beyond your progress')).toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'Reveal Episode description' }));
  expect(reveal).toHaveBeenCalledOnce();
});
