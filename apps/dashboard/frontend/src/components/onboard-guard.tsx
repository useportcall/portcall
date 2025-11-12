import { useCreateApp, useListApps } from "@/hooks";
import { zodResolver } from "@hookform/resolvers/zod";
import { ReactNode } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import FloatingLogoutButton from "./floating-logout-button";
import { Button } from "./ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Form, FormControl, FormField, FormItem, FormLabel } from "./ui/form";
import { Input } from "./ui/input";

const FormSchema = z.object({
  name: z.string().min(1, { message: "Please enter a name" }),
});

export default function OnboardGuard({ children }: { children: ReactNode }) {
  const { data: apps } = useListApps();

  const { mutate } = useCreateApp();

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: "",
    },
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    console.log("Creating app with data:", data);
    return mutate(data);
  }

  if (!apps) {
    return <></>;
  }

  if (apps.data.length) {
    return children;
  }

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10 space-mono-regular">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <div className="flex flex-col gap-6">
          <div className="w-full flex justify-center space-x-2">
            <div className="rounded-full h-2 w-2 bg-cyan-800" />
          </div>
          <Card>
            <CardHeader className="text-center">
              <CardTitle className="text-xl">
                Create your Portcall app
              </CardTitle>
              <CardDescription>Add a name for your app!</CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                  <div className="grid gap-6">
                    <FormField
                      control={form.control}
                      name="name"
                      render={({ field }) => {
                        return (
                          <FormItem className="grid">
                            <FormLabel className="text-start">Name</FormLabel>
                            <FormControl>
                              <Input
                                required
                                type="text"
                                placeholder="your app/project name"
                                {...field}
                              />
                            </FormControl>
                          </FormItem>
                        );
                      }}
                    />
                    <Button type="submit" className="w-full">
                      Let&apos;s go
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </div>
      </div>
      <FloatingLogoutButton />
    </div>
  );
}
