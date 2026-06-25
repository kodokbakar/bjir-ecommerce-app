import { useEffect, useState } from "react";
import { AlertTriangle, X } from "lucide-react";

import {
  API_SERVER_ERROR_EVENT,
  type ApiServerErrorEventDetail,
} from "../services/api";

function GlobalNetworkError() {
  const [error, setError] = useState<ApiServerErrorEventDetail | null>(null);

  useEffect(() => {
    function handleServerError(event: Event) {
      const customEvent = event as CustomEvent<ApiServerErrorEventDetail>;
      setError(customEvent.detail);
    }

    window.addEventListener(API_SERVER_ERROR_EVENT, handleServerError);

    return () => {
      window.removeEventListener(API_SERVER_ERROR_EVENT, handleServerError);
    };
  }, []);

  if (!error) {
    return null;
  }

  return (
    <aside className="global-network-error" role="alert" aria-live="assertive">
      <AlertTriangle className="h-5 w-5" aria-hidden="true" />

      <div>
        <strong>Server error {error.status}</strong>
        <p>{error.message}</p>
      </div>

      <button
        type="button"
        aria-label="Dismiss server error"
        onClick={() => setError(null)}
      >
        <X className="h-4 w-4" aria-hidden="true" />
      </button>
    </aside>
  );
}

export default GlobalNetworkError;
