import { useEffect, useRef, useState } from "react";

/**
 * Hook for managing copy-to-clipboard functionality with tooltip feedback.
 * Provides state management for open/copied states and cleanup.
 */
export function useCopyToClipboard() {
    const [open, setOpen] = useState(false);
    const [copied, setCopied] = useState(false);
    const copyRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const ref = useRef<boolean>(false);

    useEffect(() => {
        return () => {
            if (copyRef.current) {
                clearTimeout(copyRef.current);
            }
        };
    }, []);

    function copyToClipboard(text: string) {
        // Clear any existing timer to prevent multiple timers running
        if (copyRef.current) {
            clearTimeout(copyRef.current);
        }

        // Set ref to prevent tooltip flickering
        ref.current = true;

        // Copy to clipboard
        navigator.clipboard.writeText(text);

        setCopied(true);

        // Reset after 1 second
        copyRef.current = setTimeout(() => {
            ref.current = false;
            setOpen(false);
            setCopied(false);
        }, 1000);
    }

    return {
        open,
        setOpen,
        copied,
        copyToClipboard,
        ref,
    };
}
