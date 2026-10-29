import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import Header from './Header';
import { searchMulti } from '../../api/catalog';

vi.mock('../../hooks/useAuth', () => ({default:()=>({auth:null})}));
vi.mock('../../api/catalog', () => ({searchMulti:vi.fn()}));

const renderHeader = () => render(<MemoryRouter><Header handleLogout={vi.fn()}/></MemoryRouter>);

describe('Header navigation and autocomplete', () => {
  it('exposes a keyboard-operable responsive navigation toggle', async () => {
    renderHeader();
    const toggle=screen.getByRole('button',{name:'Menu'});
    expect(toggle).toHaveAttribute('aria-expanded','false');
    await userEvent.click(toggle);
    expect(toggle).toHaveAttribute('aria-expanded','true');
    expect(document.getElementById('main-nav')).toHaveClass('open');
  });

  it('renders person photography and poster artwork in suggestions', async () => {
    searchMulti.mockResolvedValue({data:{data:{results:[
      {id:1,media_type:'person',name:'Ada Actor',profile_path:'/person.jpg',known_for_department:'Acting'},
      {id:2,media_type:'movie',title:'Reel Story',poster_path:'/poster.jpg',release_date:'2026-01-01'},
    ]}}});
    renderHeader();
    fireEvent.change(screen.getByRole('combobox'),{target:{value:'re'}});
    await waitFor(()=>expect(searchMulti).toHaveBeenCalled(),{timeout:1000});
    expect(await screen.findByText('Ada Actor')).toBeInTheDocument();
    expect(screen.getByText('Reel Story')).toBeInTheDocument();
    const images=document.querySelectorAll('.autocomplete__art');
    expect(images[0]).toHaveAttribute('src',expect.stringContaining('/person.jpg'));
    expect(images[1]).toHaveAttribute('src',expect.stringContaining('/poster.jpg'));
  });
});
