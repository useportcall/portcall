import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { useGetCompany, useUpdateCompany } from "@/hooks";
import { cn } from "@/lib/utils";
import { ImagePlus, Loader2, Trash2, Upload } from "lucide-react";
import { useId, useRef, useState } from "react";
import { toast } from "sonner";
import { ImageCropperDialog } from "./image-cropper-dialog";

const ALLOWED_TYPES = ["image/jpeg", "image/png", "image/gif", "image/webp"];
const MAX_FILE_SIZE = 2 * 1024 * 1024;

export function CompanyLogoInput() {
  const { data: company } = useGetCompany();
  const { mutateAsync, isPending } = useUpdateCompany();
  const id = useId();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [cropperOpen, setCropperOpen] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const hasLogo = Boolean(company.data.icon_logo_url);

  return (
    <div className="space-y-3"><Label htmlFor={id}>Icon logo</Label><div className="flex items-center gap-4"><input id={id} ref={fileInputRef} type="file" accept={ALLOWED_TYPES.join(",")} className="hidden" onChange={(e) => { const file = e.target.files?.[0]; if (!file) return; if (fileInputRef.current) fileInputRef.current.value = ""; const error = validateFile(file); if (error) { toast.error(error); return; } setSelectedFile(file); setCropperOpen(true); }} disabled={isPending} aria-label="Upload company logo" />
      <button type="button" onClick={() => !isPending && fileInputRef.current?.click()} disabled={isPending} className={cn("relative group rounded-full focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2", isPending && "opacity-60 cursor-not-allowed")} aria-label={hasLogo ? "Change company logo" : "Upload company logo"}><Avatar className="size-20 border-2 border-dashed border-muted-foreground/30 bg-muted/50 transition-colors group-hover:border-primary/50 group-hover:bg-muted">{hasLogo ? <AvatarImage src={company.data.icon_logo_url ?? undefined} alt="Company logo" className="object-cover" /> : null}<AvatarFallback className="bg-transparent">{isPending ? <Loader2 className="size-6 text-muted-foreground animate-spin" /> : <ImagePlus className="size-6 text-muted-foreground" />}</AvatarFallback></Avatar>{hasLogo && !isPending && <div className="absolute inset-0 rounded-full bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"><Upload className="size-5 text-white" /></div>}</button>
      <div className="flex flex-col gap-1.5"><Button type="button" variant="outline" size="sm" onClick={() => fileInputRef.current?.click()} disabled={isPending} className="gap-1.5">{isPending ? <Loader2 className="size-3.5 animate-spin" /> : <Upload className="size-3.5" />}{hasLogo ? "Change" : "Upload"}</Button>{hasLogo && <Button type="button" variant="ghost" size="sm" onClick={async () => { if (!hasLogo || isPending) return; try { const formData = new FormData(); formData.append("data", JSON.stringify({ remove_logo: true })); await mutateAsync(formData); toast.success("Logo removed successfully"); } catch { toast.error("Failed to remove logo. Please try again."); } }} disabled={isPending} className="gap-1.5 text-muted-foreground hover:text-destructive"><Trash2 className="size-3.5" />Remove</Button>}</div></div><p className="text-xs text-muted-foreground">JPEG, PNG, GIF, or WebP. Max 2MB.</p>
      <ImageCropperDialog open={cropperOpen} onOpenChange={(open) => { setCropperOpen(open); if (!open) setSelectedFile(null); }} imageFile={selectedFile} onCropComplete={async (croppedBlob) => { try { const formData = new FormData(); formData.append("data", JSON.stringify({})); formData.append("logo", new File([croppedBlob], "logo.png", { type: "image/png" })); await mutateAsync(formData); toast.success("Logo updated successfully"); setSelectedFile(null); } catch { toast.error("Failed to upload logo. Please try again."); } }} outputSize={256} />
    </div>
  );
}

function validateFile(file: File) {
  if (!ALLOWED_TYPES.includes(file.type)) return "Invalid file type. Please use JPEG, PNG, GIF, or WebP.";
  if (file.size > MAX_FILE_SIZE) return `File too large (${formatFileSize(file.size)}). Maximum size is 2MB.`;
  return null;
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
