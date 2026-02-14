import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function TablePagination(props: {
  page: number;
  totalPages: number;
  previousLabel: string;
  nextLabel: string;
  pageLabel: string;
  rowsPerPageLabel: string;
  rowsPerPage: number;
  rowsPerPageOptions: number[];
  onRowsPerPageChange: (value: number) => void;
  onPrevious: () => void;
  onNext: () => void;
}) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2 mt-3 px-2">
      <div className="flex items-center gap-2">
        <span className="text-sm text-muted-foreground">{props.rowsPerPageLabel}</span>
        <Select
          value={String(props.rowsPerPage)}
          onValueChange={(value) => props.onRowsPerPageChange(Number(value))}
        >
          <SelectTrigger className="w-[110px] h-8">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {props.rowsPerPageOptions.map((value) => (
              <SelectItem key={value} value={String(value)}>
                {value}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="flex items-center gap-2">
        <div className="text-sm text-muted-foreground">{props.pageLabel}</div>
        <Button
          variant="outline"
          size="sm"
          disabled={props.page <= 1}
          onClick={props.onPrevious}
        >
          {props.previousLabel}
        </Button>
        <Button
          variant="outline"
          size="sm"
          disabled={props.page >= props.totalPages}
          onClick={props.onNext}
        >
          {props.nextLabel}
        </Button>
      </div>
    </div>
  );
}
