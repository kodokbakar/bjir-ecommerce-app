import { Component, type ErrorInfo, type ReactNode } from "react";

interface ErrorBoundaryProps {
  children: ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = {
    hasError: false,
  };

  static getDerivedStateFromError(): ErrorBoundaryState {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("Unhandled UI error:", error, errorInfo);
  }

  private handleReload = () => {
    window.location.reload();
  };

  render() {
    if (this.state.hasError) {
      return (
        <main className="error-boundary-shell">
          <section
            className="error-boundary-card"
            role="alert"
            aria-labelledby="error-boundary-title"
          >
            <span className="products-eyebrow">Runtime Fault</span>

            <h1 id="error-boundary-title">Terjadi kesalahan</h1>

            <p>
              Aplikasi mengalami masalah tak terduga. Muat ulang halaman untuk
              mencoba lagi.
            </p>

            <button
              className="error-boundary-button"
              type="button"
              onClick={this.handleReload}
            >
              Muat ulang
            </button>
          </section>
        </main>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
