import { Component, type ErrorInfo, type ReactNode } from 'react'
import { C } from '../styles/tokens'

interface ErrorBoundaryProps {
  children: ReactNode
}

interface ErrorBoundaryState {
  hasError: boolean
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = {
    hasError: false,
  }

  static getDerivedStateFromError(): ErrorBoundaryState {
    return { hasError: true }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Unhandled UI error:', error, errorInfo)
  }

  private handleReload = () => {
    window.location.reload()
  }

  render() {
    if (this.state.hasError) {
      return (
        <div
          style={{
            minHeight: '100vh',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: C.secondary,
            padding: 24,
          }}
        >
          <div
            role="alert"
            style={{
              width: '100%',
              maxWidth: 420,
              background: '#fff',
              border: `1px solid ${C.border}`,
              borderRadius: 12,
              padding: 24,
              textAlign: 'center',
            }}
          >
            <h1
              style={{
                margin: '0 0 8px',
                color: C.textDark,
                fontSize: 20,
                fontWeight: 600,
              }}
            >
              Terjadi kesalahan
            </h1>

            <p
              style={{
                margin: '0 0 20px',
                color: C.textMuted,
                fontSize: 14,
                lineHeight: 1.6,
              }}
            >
              Aplikasi mengalami masalah tak terduga. Muat ulang halaman untuk mencoba lagi.
            </p>

            <button
              type="button"
              onClick={this.handleReload}
              style={{
                height: 40,
                padding: '0 16px',
                border: 'none',
                borderRadius: 8,
                background: C.primary,
                color: '#fff',
                cursor: 'pointer',
                fontWeight: 500,
              }}
            >
              Muat ulang
            </button>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}

export default ErrorBoundary