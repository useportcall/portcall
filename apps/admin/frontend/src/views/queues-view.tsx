import { Suspense, useState } from "react";
import { useQueues, useEnqueueTask, TaskPayload } from "@/hooks/use-queues";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Play, Clock, Zap, CheckCircle } from "lucide-react";
import { toast } from "@/lib/toast";

function QueuesContent() {
  const { data: queues } = useQueues();
  const enqueueMutation = useEnqueueTask();

  const [selectedTask, setSelectedTask] = useState<{
    queue: string;
    task: string;
    params: { name: string; type: string; required: boolean }[];
  } | null>(null);

  const [formValues, setFormValues] = useState<Record<string, string>>({});

  const handleOpenDialog = (
    queue: string,
    taskName: string,
    params: { name: string; type: string; required: boolean }[],
  ) => {
    setSelectedTask({ queue, task: taskName, params });
    setFormValues({});
  };

  const handleEnqueue = async () => {
    if (!selectedTask) return;

    const payload: TaskPayload = {
      queue_name: selectedTask.queue,
      task_type: selectedTask.task,
      payload: {},
    };

    // Build payload from form values
    selectedTask.params.forEach((param) => {
      const value = formValues[param.name];
      if (value !== undefined && value !== "") {
        // Parse based on type
        if (
          param.type === "number" ||
          param.type === "integer" ||
          param.type === "int"
        ) {
          payload.payload[param.name] = parseInt(value, 10);
        } else if (param.type === "float" || param.type === "float64") {
          payload.payload[param.name] = parseFloat(value);
        } else if (param.type === "boolean" || param.type === "bool") {
          payload.payload[param.name] = value === "true";
        } else {
          payload.payload[param.name] = value;
        }
      }
    });

    try {
      await enqueueMutation.mutateAsync(payload);
      toast.success(`Task "${selectedTask.task}" enqueued successfully`);
      setSelectedTask(null);
      setFormValues({});
    } catch (error) {
      toast.error(`Failed to enqueue task: ${(error as Error).message}`);
    }
  };

  const getQueueColor = (queueName: string) => {
    if (queueName.includes("billing")) return "text-green-600";
    if (queueName.includes("email")) return "text-blue-600";
    return "text-purple-600";
  };

  const getQueueIcon = (queueName: string) => {
    if (queueName.includes("billing")) return Zap;
    if (queueName.includes("email")) return Clock;
    return Play;
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Queue Jobs</h1>
        <p className="text-muted-foreground mt-1">
          Trigger background queue tasks manually
        </p>
      </div>

      <Tabs defaultValue={queues[0]?.name || "billing_queue"}>
        <TabsList>
          {queues.map((queue) => {
            const Icon = getQueueIcon(queue.name);
            return (
              <TabsTrigger
                key={queue.name}
                value={queue.name}
                className="flex items-center gap-2"
              >
                <Icon className={`h-4 w-4 ${getQueueColor(queue.name)}`} />
                {queue.name
                  .replace("_", " ")
                  .replace(/\b\w/g, (l) => l.toUpperCase())}
                <Badge variant="secondary" className="ml-1">
                  {queue.tasks.length}
                </Badge>
              </TabsTrigger>
            );
          })}
        </TabsList>

        {queues.map((queue) => (
          <TabsContent key={queue.name} value={queue.name} className="mt-4">
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {queue.tasks.map((task) => (
                <Card key={task.name} className="flex flex-col">
                  <CardHeader className="pb-2">
                    <CardTitle className="text-lg flex items-center justify-between">
                      <span className="truncate">{task.name}</span>
                      <Badge variant="outline" className="ml-2 shrink-0">
                        {task.params.length} params
                      </Badge>
                    </CardTitle>
                    <CardDescription className="text-xs">
                      {task.description || "No description available"}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="flex-1 flex flex-col justify-end">
                    {task.params.length > 0 && (
                      <div className="mb-3 space-y-1">
                        <div className="text-xs text-muted-foreground font-medium">
                          Parameters:
                        </div>
                        <div className="flex flex-wrap gap-1">
                          {task.params.map((param) => (
                            <Badge
                              key={param.name}
                              variant={param.required ? "default" : "secondary"}
                              className="text-xs"
                            >
                              {param.name}
                              {param.required && "*"}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}
                    <Button
                      className="w-full"
                      onClick={() =>
                        handleOpenDialog(queue.name, task.name, task.params)
                      }
                    >
                      <Play className="h-4 w-4 mr-2" />
                      Enqueue
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          </TabsContent>
        ))}
      </Tabs>

      {/* Enqueue Dialog */}
      <Dialog
        open={!!selectedTask}
        onOpenChange={(open) => !open && setSelectedTask(null)}
      >
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Enqueue Task</DialogTitle>
            <DialogDescription>
              Fill in the parameters to enqueue "{selectedTask?.task}" to{" "}
              {selectedTask?.queue}
            </DialogDescription>
          </DialogHeader>
          <ScrollArea className="max-h-[400px] pr-4">
            <div className="space-y-4">
              {selectedTask?.params.length === 0 && (
                <p className="text-sm text-muted-foreground">
                  This task has no parameters.
                </p>
              )}
              {selectedTask?.params.map((param) => (
                <div key={param.name} className="space-y-2">
                  <Label
                    htmlFor={param.name}
                    className="flex items-center gap-2"
                  >
                    {param.name}
                    {param.required && (
                      <span className="text-destructive">*</span>
                    )}
                    <Badge variant="outline" className="text-xs">
                      {param.type}
                    </Badge>
                  </Label>
                  {param.type === "boolean" || param.type === "bool" ? (
                    <Select
                      value={formValues[param.name] || ""}
                      onValueChange={(value) =>
                        setFormValues((prev) => ({
                          ...prev,
                          [param.name]: value,
                        }))
                      }
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select value" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="true">True</SelectItem>
                        <SelectItem value="false">False</SelectItem>
                      </SelectContent>
                    </Select>
                  ) : (
                    <Input
                      id={param.name}
                      type={
                        param.type === "number" ||
                        param.type === "integer" ||
                        param.type === "int" ||
                        param.type === "float" ||
                        param.type === "float64"
                          ? "number"
                          : "text"
                      }
                      placeholder={`Enter ${param.name}`}
                      value={formValues[param.name] || ""}
                      onChange={(e) =>
                        setFormValues((prev) => ({
                          ...prev,
                          [param.name]: e.target.value,
                        }))
                      }
                    />
                  )}
                </div>
              ))}
            </div>
          </ScrollArea>
          <Separator />
          <DialogFooter>
            <Button variant="outline" onClick={() => setSelectedTask(null)}>
              Cancel
            </Button>
            <Button
              onClick={handleEnqueue}
              disabled={enqueueMutation.isPending}
            >
              {enqueueMutation.isPending ? (
                <>
                  <Clock className="h-4 w-4 mr-2 animate-spin" />
                  Enqueueing...
                </>
              ) : (
                <>
                  <CheckCircle className="h-4 w-4 mr-2" />
                  Enqueue Task
                </>
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

export default function QueuesView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <div className="animate-pulse text-muted-foreground">Loading...</div>
        </div>
      }
    >
      <QueuesContent />
    </Suspense>
  );
}
