import React, { Suspense, useState, useMemo } from "react";
import {
  useTasks,
  useQueueStats,
  useArchiveTask,
  useDeleteTask,
  useRunTask,
  useRetryWithPayload,
  TaskInfo,
} from "@/hooks/use-tasks";
import { useQueues } from "@/hooks/use-queues";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Play,
  Archive,
  Trash2,
  RefreshCw,
  AlertCircle,
  CheckCircle,
  Clock,
  XCircle,
  Zap,
  Mail,
  ChevronDown,
  ChevronUp,
  Edit,
  Eye,
} from "lucide-react";
import { toast } from "@/lib/toast";
import { formatDistanceToNow } from "date-fns";

const STATE_COLORS: Record<string, string> = {
  active: "bg-blue-500",
  pending: "bg-yellow-500",
  scheduled: "bg-purple-500",
  retry: "bg-orange-500",
  archived: "bg-red-500",
  completed: "bg-green-500",
};

const STATE_ICONS: Record<
  string,
  React.ComponentType<{ className?: string }>
> = {
  active: RefreshCw,
  pending: Clock,
  scheduled: Clock,
  retry: AlertCircle,
  archived: Archive,
  completed: CheckCircle,
};

function formatTime(timeStr?: string) {
  if (!timeStr) return "-";
  try {
    const date = new Date(timeStr);
    return formatDistanceToNow(date, { addSuffix: true });
  } catch {
    return timeStr;
  }
}

function TaskStateCell({ state }: { state: string }) {
  const Icon = STATE_ICONS[state] || AlertCircle;
  return (
    <div className="flex items-center gap-1.5">
      <div
        className={`w-2 h-2 rounded-full ${STATE_COLORS[state] || "bg-gray-500"}`}
      />
      <Icon className="h-3.5 w-3.5 text-muted-foreground" />
      <span className="text-xs capitalize">{state}</span>
    </div>
  );
}

function TaskDetailDialog({
  task,
  open,
  onClose,
  onArchive,
  onDelete,
  onRun,
  onRetry,
}: {
  task: TaskInfo | null;
  open: boolean;
  onClose: () => void;
  onArchive: (task: TaskInfo) => void;
  onDelete: (task: TaskInfo) => void;
  onRun: (task: TaskInfo) => void;
  onRetry: (task: TaskInfo, payload: string) => void;
}) {
  const [editedPayload, setEditedPayload] = useState("");
  const [isEditing, setIsEditing] = useState(false);

  const formattedPayload = useMemo(() => {
    if (!task?.payload) return "{}";
    try {
      return JSON.stringify(JSON.parse(task.payload), null, 2);
    } catch {
      return task.payload;
    }
  }, [task?.payload]);

  const handleEdit = () => {
    setEditedPayload(formattedPayload);
    setIsEditing(true);
  };

  const handleRetry = () => {
    if (!task) return;
    onRetry(task, isEditing ? editedPayload : formattedPayload);
    setIsEditing(false);
  };

  if (!task) return null;

  return (
    <Dialog open={open} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-hidden flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <span className="font-mono text-sm">{task.type}</span>
            <TaskStateCell state={task.state} />
          </DialogTitle>
          <DialogDescription className="font-mono text-xs">
            ID: {task.id}
          </DialogDescription>
        </DialogHeader>

        <ScrollArea className="flex-1 pr-4">
          <div className="space-y-4">
            {/* Error Section */}
            {task.last_err && (
              <div className="rounded-md border border-red-200 bg-red-50 p-3">
                <div className="flex items-center gap-2 text-red-700 mb-2">
                  <XCircle className="h-4 w-4" />
                  <span className="font-medium text-sm">Error</span>
                  {task.last_failed_at && (
                    <span className="text-xs text-red-500 ml-auto">
                      {formatTime(task.last_failed_at)}
                    </span>
                  )}
                </div>
                <pre className="text-xs text-red-600 whitespace-pre-wrap font-mono overflow-auto max-h-32">
                  {task.last_err}
                </pre>
              </div>
            )}

            {/* Metadata */}
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div>
                <span className="text-muted-foreground">Queue:</span>{" "}
                <span className="font-medium">{task.queue}</span>
              </div>
              <div>
                <span className="text-muted-foreground">Retried:</span>{" "}
                <span className="font-medium">
                  {task.retried} / {task.max_retry}
                </span>
              </div>
              {task.next_process_at && (
                <div>
                  <span className="text-muted-foreground">Next process:</span>{" "}
                  <span className="font-medium">
                    {formatTime(task.next_process_at)}
                  </span>
                </div>
              )}
              {task.completed_at && (
                <div>
                  <span className="text-muted-foreground">Completed:</span>{" "}
                  <span className="font-medium">
                    {formatTime(task.completed_at)}
                  </span>
                </div>
              )}
            </div>

            <Separator />

            {/* Payload */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <Label className="text-sm font-medium">Payload</Label>
                {(task.state === "archived" || task.state === "retry") && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleEdit}
                    className="h-7"
                  >
                    <Edit className="h-3 w-3 mr-1" />
                    Edit & Retry
                  </Button>
                )}
              </div>
              {isEditing ? (
                <Textarea
                  value={editedPayload}
                  onChange={(e) => setEditedPayload(e.target.value)}
                  className="font-mono text-xs min-h-[200px]"
                />
              ) : (
                <pre className="bg-muted rounded-md p-3 text-xs font-mono overflow-auto max-h-60">
                  {formattedPayload}
                </pre>
              )}
            </div>
          </div>
        </ScrollArea>

        <Separator />
        <DialogFooter className="gap-2 sm:gap-0">
          <div className="flex gap-2 w-full sm:w-auto">
            {(task.state === "pending" ||
              task.state === "retry" ||
              task.state === "scheduled") && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => onArchive(task)}
                className="text-orange-600 border-orange-200 hover:bg-orange-50"
              >
                <Archive className="h-4 w-4 mr-1" />
                Archive
              </Button>
            )}
            {(task.state === "archived" || task.state === "retry") && (
              <>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onRun(task)}
                  className="text-green-600 border-green-200 hover:bg-green-50"
                >
                  <Play className="h-4 w-4 mr-1" />
                  Run Now
                </Button>
                <Button variant="default" size="sm" onClick={handleRetry}>
                  <RefreshCw className="h-4 w-4 mr-1" />
                  {isEditing ? "Retry with Changes" : "Retry"}
                </Button>
              </>
            )}
            {task.state === "archived" && (
              <Button
                variant="destructive"
                size="sm"
                onClick={() => onDelete(task)}
              >
                <Trash2 className="h-4 w-4 mr-1" />
                Delete
              </Button>
            )}
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function QueueStatsCards({ queue }: { queue: string }) {
  const { data: stats } = useQueueStats();
  const queueStats = stats?.find((s) => s.queue === queue);

  if (!queueStats) return null;

  const items = [
    { label: "Active", value: queueStats.active, color: "text-blue-600" },
    { label: "Pending", value: queueStats.pending, color: "text-yellow-600" },
    { label: "Retry", value: queueStats.retry, color: "text-orange-600" },
    { label: "Archived", value: queueStats.archived, color: "text-red-600" },
    {
      label: "Completed",
      value: queueStats.completed,
      color: "text-green-600",
    },
  ];

  return (
    <div className="flex gap-3 flex-wrap">
      {items.map((item) => (
        <div key={item.label} className="flex items-center gap-1.5 text-sm">
          <span className="text-muted-foreground">{item.label}:</span>
          <span className={`font-semibold ${item.color}`}>{item.value}</span>
        </div>
      ))}
    </div>
  );
}

function TasksTableContent({
  queue,
  state,
  taskType,
  limit,
}: {
  queue: string;
  state: string;
  taskType: string;
  limit: number;
}) {
  const { data } = useTasks({
    queue,
    state: state || undefined,
    taskType: taskType || undefined,
    limit,
  });
  const archiveMutation = useArchiveTask();
  const deleteMutation = useDeleteTask();
  const runMutation = useRunTask();
  const retryMutation = useRetryWithPayload();

  const [selectedTask, setSelectedTask] = useState<TaskInfo | null>(null);
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());

  const handleArchive = async (task: TaskInfo) => {
    try {
      await archiveMutation.mutateAsync({ queue: task.queue, taskId: task.id });
      toast.success("Task archived");
      setSelectedTask(null);
    } catch (err) {
      toast.error("Failed to archive task");
    }
  };

  const handleDelete = async (task: TaskInfo) => {
    try {
      await deleteMutation.mutateAsync({ queue: task.queue, taskId: task.id });
      toast.success("Task deleted");
      setSelectedTask(null);
    } catch (err) {
      toast.error("Failed to delete task");
    }
  };

  const handleRun = async (task: TaskInfo) => {
    try {
      await runMutation.mutateAsync({ queue: task.queue, taskId: task.id });
      toast.success("Task queued to run");
      setSelectedTask(null);
    } catch (err) {
      toast.error("Failed to run task");
    }
  };

  const handleRetry = async (task: TaskInfo, payload: string) => {
    try {
      await retryMutation.mutateAsync({
        queue: task.queue,
        taskType: task.type,
        payload,
      });
      toast.success("New task created");
      setSelectedTask(null);
    } catch (err) {
      toast.error("Failed to create task");
    }
  };

  const toggleExpand = (id: string) => {
    setExpandedRows((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  if (!data?.tasks || data.tasks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
        <CheckCircle className="h-8 w-8 mb-2" />
        <p>No tasks found</p>
      </div>
    );
  }

  return (
    <>
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="w-[40px]"></TableHead>
            <TableHead className="w-[200px]">Task Type</TableHead>
            <TableHead className="w-[100px]">State</TableHead>
            <TableHead className="w-[80px]">Retries</TableHead>
            <TableHead className="w-[140px]">Last Failed</TableHead>
            <TableHead className="w-[100px]">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.tasks.map((task) => {
            const isExpanded = expandedRows.has(task.id);
            const hasError = !!task.last_err;

            return (
              <React.Fragment key={task.id}>
                <TableRow
                  className={`cursor-pointer ${hasError ? "bg-red-50/50" : ""}`}
                  onClick={() => toggleExpand(task.id)}
                >
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-6 w-6 p-0"
                      onClick={(e) => {
                        e.stopPropagation();
                        toggleExpand(task.id);
                      }}
                    >
                      {isExpanded ? (
                        <ChevronUp className="h-4 w-4" />
                      ) : (
                        <ChevronDown className="h-4 w-4" />
                      )}
                    </Button>
                  </TableCell>
                  <TableCell>
                    <span className="font-mono text-xs">{task.type}</span>
                  </TableCell>
                  <TableCell>
                    <TaskStateCell state={task.state} />
                  </TableCell>
                  <TableCell>
                    <span className="text-xs">
                      {task.retried} / {task.max_retry}
                    </span>
                  </TableCell>
                  <TableCell>
                    <span className="text-xs text-muted-foreground">
                      {formatTime(task.last_failed_at)}
                    </span>
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 px-2"
                      onClick={(e) => {
                        e.stopPropagation();
                        setSelectedTask(task);
                      }}
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
                {isExpanded && (
                  <TableRow className="bg-muted/30 hover:bg-muted/30">
                    <TableCell colSpan={6} className="p-0">
                      <div className="px-4 py-3 space-y-2">
                        <div className="flex items-center gap-4 text-xs">
                          <span className="text-muted-foreground">ID:</span>
                          <code className="font-mono">{task.id}</code>
                        </div>
                        {task.last_err && (
                          <div className="bg-red-50 border border-red-200 rounded p-2">
                            <div className="flex items-center gap-1 text-red-700 text-xs mb-1">
                              <XCircle className="h-3 w-3" />
                              <span className="font-medium">Error:</span>
                            </div>
                            <pre className="text-xs text-red-600 font-mono whitespace-pre-wrap">
                              {task.last_err}
                            </pre>
                          </div>
                        )}
                        <div className="flex gap-2">
                          <Button
                            variant="outline"
                            size="sm"
                            className="h-7 text-xs"
                            onClick={(e) => {
                              e.stopPropagation();
                              setSelectedTask(task);
                            }}
                          >
                            <Eye className="h-3 w-3 mr-1" />
                            View Details
                          </Button>
                          {(task.state === "pending" ||
                            task.state === "retry" ||
                            task.state === "scheduled") && (
                            <Button
                              variant="outline"
                              size="sm"
                              className="h-7 text-xs text-orange-600 border-orange-200"
                              onClick={(e) => {
                                e.stopPropagation();
                                handleArchive(task);
                              }}
                            >
                              <Archive className="h-3 w-3 mr-1" />
                              Archive
                            </Button>
                          )}
                          {(task.state === "archived" ||
                            task.state === "retry") && (
                            <Button
                              variant="outline"
                              size="sm"
                              className="h-7 text-xs text-green-600 border-green-200"
                              onClick={(e) => {
                                e.stopPropagation();
                                handleRun(task);
                              }}
                            >
                              <Play className="h-3 w-3 mr-1" />
                              Run Now
                            </Button>
                          )}
                        </div>
                      </div>
                    </TableCell>
                  </TableRow>
                )}
              </React.Fragment>
            );
          })}
        </TableBody>
      </Table>

      <TaskDetailDialog
        task={selectedTask}
        open={!!selectedTask}
        onClose={() => setSelectedTask(null)}
        onArchive={handleArchive}
        onDelete={handleDelete}
        onRun={handleRun}
        onRetry={handleRetry}
      />
    </>
  );
}

function JobsContent() {
  const { data: queues } = useQueues();
  const [selectedQueue, setSelectedQueue] = useState("billing_queue");
  const [selectedState, setSelectedState] = useState("all");
  const [selectedTaskType, setSelectedTaskType] = useState("all");
  const [limit, setLimit] = useState(50);

  // Get task types for the selected queue
  const taskTypes = useMemo(() => {
    const queue = queues.find((q) => q.name === selectedQueue);
    return queue?.tasks.map((t) => t.name) || [];
  }, [queues, selectedQueue]);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Queue Jobs</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            Monitor and manage queue tasks
          </p>
        </div>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-4">
          <div className="flex flex-wrap gap-4">
            {/* Queue Selector */}
            <div className="flex items-center gap-2">
              <Label className="text-sm whitespace-nowrap">Queue:</Label>
              <Select value={selectedQueue} onValueChange={setSelectedQueue}>
                <SelectTrigger className="w-[160px] h-8">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="billing_queue">
                    <div className="flex items-center gap-2">
                      <Zap className="h-4 w-4 text-green-600" />
                      billing_queue
                    </div>
                  </SelectItem>
                  <SelectItem value="email_queue">
                    <div className="flex items-center gap-2">
                      <Mail className="h-4 w-4 text-blue-600" />
                      email_queue
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* State Filter */}
            <div className="flex items-center gap-2">
              <Label className="text-sm whitespace-nowrap">State:</Label>
              <Select value={selectedState} onValueChange={setSelectedState}>
                <SelectTrigger className="w-[130px] h-8">
                  <SelectValue placeholder="All states" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All states</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="scheduled">Scheduled</SelectItem>
                  <SelectItem value="retry">Retry</SelectItem>
                  <SelectItem value="archived">Archived</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Task Type Filter */}
            <div className="flex items-center gap-2">
              <Label className="text-sm whitespace-nowrap">Task:</Label>
              <Select
                value={selectedTaskType}
                onValueChange={setSelectedTaskType}
              >
                <SelectTrigger className="w-[180px] h-8">
                  <SelectValue placeholder="All types" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All types</SelectItem>
                  {taskTypes.map((type) => (
                    <SelectItem key={type} value={type}>
                      {type}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Limit */}
            <div className="flex items-center gap-2">
              <Label className="text-sm whitespace-nowrap">Show:</Label>
              <Select
                value={limit.toString()}
                onValueChange={(v) => setLimit(parseInt(v))}
              >
                <SelectTrigger className="w-[80px] h-8">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="50">50</SelectItem>
                  <SelectItem value="100">100</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Stats */}
            <div className="ml-auto">
              <Suspense
                fallback={
                  <span className="text-sm text-muted-foreground">
                    Loading stats...
                  </span>
                }
              >
                <QueueStatsCards queue={selectedQueue} />
              </Suspense>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          <Suspense
            fallback={
              <div className="flex items-center justify-center py-12">
                <RefreshCw className="h-6 w-6 animate-spin text-muted-foreground" />
              </div>
            }
          >
            <TasksTableContent
              queue={selectedQueue}
              state={selectedState === "all" ? "" : selectedState}
              taskType={selectedTaskType === "all" ? "" : selectedTaskType}
              limit={limit}
            />
          </Suspense>
        </CardContent>
      </Card>
    </div>
  );
}

export default function JobsView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <div className="animate-pulse text-muted-foreground">Loading...</div>
        </div>
      }
    >
      <JobsContent />
    </Suspense>
  );
}
