// Simple toast implementation
type ToastType = "success" | "error" | "info";

interface ToastOptions {
  duration?: number;
}

class Toast {
  private show(message: string, type: ToastType, options?: ToastOptions) {
    const duration = options?.duration || 3000;

    // Create toast element
    const toast = document.createElement("div");
    toast.className = `fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg z-50 transform transition-all duration-300 translate-y-full opacity-0 ${
      type === "success"
        ? "bg-green-600 text-white"
        : type === "error"
          ? "bg-red-600 text-white"
          : "bg-gray-800 text-white"
    }`;
    toast.textContent = message;

    document.body.appendChild(toast);

    // Animate in
    requestAnimationFrame(() => {
      toast.classList.remove("translate-y-full", "opacity-0");
    });

    // Remove after duration
    setTimeout(() => {
      toast.classList.add("translate-y-full", "opacity-0");
      setTimeout(() => {
        document.body.removeChild(toast);
      }, 300);
    }, duration);
  }

  success(message: string, options?: ToastOptions) {
    this.show(message, "success", options);
  }

  error(message: string, options?: ToastOptions) {
    this.show(message, "error", options);
  }

  info(message: string, options?: ToastOptions) {
    this.show(message, "info", options);
  }
}

export const toast = new Toast();
