// src/components/save-indicator.tsx
import { useEffect, useState, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";

export function SaveIndicator() {
  const [visible, setVisible] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    function handleSaved() {
      // Clear any existing timer to prevent multiple timers running
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }

      setVisible(true);

      // Store the timer ID in a ref
      timerRef.current = setTimeout(() => {
        setVisible(false);
      }, 3000);
    }

    window.addEventListener("saved", handleSaved);

    // This cleanup function runs when the component unmounts
    return () => {
      window.removeEventListener("saved", handleSaved);
      // Clear the timer in the cleanup to prevent the setState on unmounted component
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);

  return (
    <AnimatePresence>
      {visible && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.3 }}
          className="self-center text-sm text-slate-400 font-medium"
        >
          Saved.
        </motion.div>
      )}
    </AnimatePresence>
  );
}
