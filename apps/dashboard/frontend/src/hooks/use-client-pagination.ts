import { useEffect, useMemo, useState } from "react";

export function useClientPagination<T>(items: T[], initialRowsPerPage = 20) {
  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(initialRowsPerPage);
  const total = items.length;
  const totalPages = Math.max(1, Math.ceil(total / rowsPerPage));
  const safePage = Math.min(page, totalPages);

  useEffect(() => {
    setPage((value) => Math.min(value, totalPages));
  }, [totalPages]);

  const pagedItems = useMemo(() => {
    const start = (safePage - 1) * rowsPerPage;
    return items.slice(start, start + rowsPerPage);
  }, [items, rowsPerPage, safePage]);

  return {
    page: safePage,
    total,
    totalPages,
    rowsPerPage,
    pagedItems,
    onPrevious: () => setPage((value) => Math.max(1, value - 1)),
    onNext: () => setPage((value) => Math.min(totalPages, value + 1)),
    onRowsPerPageChange: (value: number) => {
      setRowsPerPage(value);
      setPage(1);
    },
    resetPage: () => setPage(1),
  };
}
