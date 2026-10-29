import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import LoadingScreen from './LoadingScreen';

describe('LoadingScreen', () => {
  it('announces progress while keeping the projector machine decorative', () => {
    render(<LoadingScreen/>);
    expect(screen.getByRole('status')).toHaveTextContent('The programme starts shortly');
    expect(screen.getByText('Loading Plotlyst')).toHaveClass('visually-hidden');
    expect(screen.getByRole('status').querySelector('.plotlyst-loader__machine')).toHaveAttribute('aria-hidden','true');
  });

  it('supports the compact route state', () => {
    render(<LoadingScreen compact/>);
    expect(screen.getByRole('status')).toHaveClass('plotlyst-loader--compact');
  });
});
