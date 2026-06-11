import { render, screen } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router-dom'
import { beforeEach, describe, expect, it, type Mock, vi } from 'vitest'
import { useAuth } from '../hooks/useAuth'
import ProtectedRoute from './ProtectedRoute'

vi.mock('../hooks/useAuth', () => ({
  useAuth: vi.fn(),
}))

const mockUseAuth = useAuth as unknown as Mock

const renderProtectedRoute = (initialPath = '/dashboard') => {
  const router = createMemoryRouter(
    [
      {
        path: '/',
        element: <ProtectedRoute />,
        children: [
          {
            path: 'dashboard',
            element: <p>Dashboard privat</p>,
          },
        ],
      },
      {
        path: '/login',
        element: <p>Halaman login</p>,
      },
    ],
    {
      initialEntries: [initialPath],
    },
  )

  render(<RouterProvider router={router} />)
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    mockUseAuth.mockReset()
  })

  it('shows loading state while auth is initializing', () => {
    mockUseAuth.mockReturnValue({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: true,
      login: vi.fn(),
      logout: vi.fn(),
    })

    renderProtectedRoute()

    expect(screen.getByText('Memeriksa sesi...')).toBeInTheDocument()
  })

  it('redirects unauthenticated user to login page', async () => {
    mockUseAuth.mockReturnValue({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
    })

    renderProtectedRoute()

    expect(await screen.findByText('Halaman login')).toBeInTheDocument()
  })

  it('renders protected content for authenticated user', () => {
    mockUseAuth.mockReturnValue({
      user: {
        id: 1,
        name: 'Bintang',
        email: 'bintang@example.com',
        role: 'customer',
      },
      token: 'valid-token',
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
    })

    renderProtectedRoute()

    expect(screen.getByText('Dashboard privat')).toBeInTheDocument()
  })
})