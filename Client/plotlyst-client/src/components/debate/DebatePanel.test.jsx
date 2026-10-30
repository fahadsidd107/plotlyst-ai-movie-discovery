import { fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { beforeEach, expect, test, vi } from 'vitest';
import DebatePanel from './DebatePanel';
import { debateApi } from '../../api/debate';

vi.mock('../../api/debate', () => ({ debateApi: { movie: vi.fn() }, debateMessage: error => error.message }));

beforeEach(() => vi.clearAllMocks());

test('shows an honest insufficient-data state without sentiment claims', async () => {
  debateApi.movie.mockResolvedValue({ status: 'insufficient', review_count: 2, minimum_reviews: 5, rating_distribution: { 5: 2 }, loved: [], disliked: [] });
  render(<MemoryRouter><DebatePanel movieId="10" title="Example" /></MemoryRouter>);
  expect(await screen.findByText('2 of 5 reviews collected.')).toBeInTheDocument();
  expect(screen.getByText(/No community sentiment has been generated yet/)).toBeInTheDocument();
});

test('expands grounded evidence and links to its original review', async () => {
  debateApi.movie.mockResolvedValue({ status: 'available', review_count: 5, minimum_reviews: 5, analysis_method: 'deterministic-grounded-v1', rating_distribution: { 1: 1, 2: 0, 3: 1, 4: 1, 5: 2 }, loved: [{ theme: 'acting', summary: 'Acting was repeatedly praised.', mention_count: 1, evidence: [{ review_id: 'abc123', author_name: 'Mina', stars: 5, excerpt: 'A remarkable performance.' }] }], disliked: [], polarization: { summary: 'The ballot leans positive.' } });
  render(<MemoryRouter><DebatePanel movieId="10" title="Example" /></MemoryRouter>);
  fireEvent.click(await screen.findByRole('button', { name: /Open evidence/ }));
  expect(screen.getByText('A remarkable performance.')).toBeInTheDocument();
  expect(screen.getByRole('link', { name: /Original review/ })).toHaveAttribute('href', '/movie/10?review=abc123#review-abc123');
});
