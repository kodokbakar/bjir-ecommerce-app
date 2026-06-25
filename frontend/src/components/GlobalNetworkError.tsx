import { useEffect } from "react";

import {
  API_SERVER_ERROR_EVENT,
  type ApiServerErrorEventDetail,
} from "../services/api";
import { useToast } from "../context/toast"

function GlobalNetworkError() {
  const { showToast } = useToast();

  useEffect(() => {
    function handleServerError(event: Event) {
      const customEvent = event as CustomEvent<ApiServerErrorEventDetail>;

      showToast(
        {
          type: "error",
          title: `Server error ${customEvent.detail.status}`,
          message: customEvent.detail.message,
        },
        {
          duration: 6000,
        },
      );
    }

    window.addEventListener(API_SERVER_ERROR_EVENT, handleServerError);

    return () => {
      window.removeEventListener(API_SERVER_ERROR_EVENT, handleServerError);
    };
  }, [showToast]);

  return null;
}

export default GlobalNetworkError;
