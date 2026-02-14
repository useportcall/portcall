import { useEffect, useRef, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

export function useInputSaveIndicator() {
  const [visible, setVisible] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const showSaved = () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
    }

    setVisible(true);

    timerRef.current = setTimeout(() => {
      setVisible(false);
    }, 3000);
  };

  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);

  return { visible, showSaved };
}

export function InputSaveIndicator({ visible }: { visible: boolean }) {
  return (
    <AnimatePresence>
      {visible && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.3 }}
          className="text-sm text-muted-foreground font-medium"
        >
          Saved.
        </motion.div>
      )}
    </AnimatePresence>
  );
}
