import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Slider } from "@/components/ui/slider";
import { Move, ZoomIn, ZoomOut } from "lucide-react";
import { useCallback, useEffect, useRef, useState, type PointerEvent as ReactPointerEvent } from "react";

const MIN_ZOOM = 0.5;
const MAX_ZOOM = 3;
const DEFAULT_ZOOM = 1;
const PREVIEW_SIZE = 280;

export function ImageCropperDialog({ open, onOpenChange, imageFile, onCropComplete, outputSize = 256 }: { open: boolean; onOpenChange: (open: boolean) => void; imageFile: File | null; onCropComplete: (croppedBlob: Blob) => void; outputSize?: number; }) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const imageRef = useRef<HTMLImageElement | null>(null);
  const [zoom, setZoom] = useState(DEFAULT_ZOOM);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const [imageLoaded, setImageLoaded] = useState(false);

  useEffect(() => {
    if (!imageFile) return void setImageLoaded(false);
    const img = new Image(); const url = URL.createObjectURL(imageFile);
    img.onload = () => { imageRef.current = img; setImageLoaded(true); setZoom(DEFAULT_ZOOM); setPosition({ x: 0, y: 0 }); };
    img.onerror = () => setImageLoaded(false); img.src = url;
    return () => URL.revokeObjectURL(url);
  }, [imageFile]);

  const drawPreview = useCallback(() => {
    const canvas = canvasRef.current; const img = imageRef.current; if (!canvas || !img || !imageLoaded) return;
    const ctx = canvas.getContext("2d"); if (!ctx) return; ctx.clearRect(0, 0, PREVIEW_SIZE, PREVIEW_SIZE); ctx.save(); ctx.beginPath(); ctx.arc(PREVIEW_SIZE / 2, PREVIEW_SIZE / 2, PREVIEW_SIZE / 2, 0, Math.PI * 2); ctx.clip();
    const { drawWidth, drawHeight } = getDrawSize(img, zoom, PREVIEW_SIZE); const x = (PREVIEW_SIZE - drawWidth) / 2 + position.x; const y = (PREVIEW_SIZE - drawHeight) / 2 + position.y; ctx.drawImage(img, x, y, drawWidth, drawHeight); ctx.restore();
    ctx.strokeStyle = "rgba(255, 255, 255, 0.3)"; ctx.lineWidth = 2; ctx.beginPath(); ctx.arc(PREVIEW_SIZE / 2, PREVIEW_SIZE / 2, PREVIEW_SIZE / 2 - 1, 0, Math.PI * 2); ctx.stroke();
  }, [zoom, position, imageLoaded]);
  useEffect(() => { drawPreview(); }, [drawPreview]);

  const onPointerDown = (e: ReactPointerEvent<HTMLCanvasElement>) => { e.preventDefault(); setIsDragging(true); setDragStart({ x: e.clientX - position.x, y: e.clientY - position.y }); (e.target as HTMLElement).setPointerCapture(e.pointerId); };
  const onPointerMove = (e: ReactPointerEvent<HTMLCanvasElement>) => isDragging && setPosition({ x: e.clientX - dragStart.x, y: e.clientY - dragStart.y });
  const onPointerUp = (e: ReactPointerEvent<HTMLCanvasElement>) => { setIsDragging(false); (e.target as HTMLElement).releasePointerCapture(e.pointerId); };

  const handleApply = () => {
    const img = imageRef.current; if (!img) return; const outputCanvas = document.createElement("canvas"); outputCanvas.width = outputSize; outputCanvas.height = outputSize;
    const ctx = outputCanvas.getContext("2d"); if (!ctx) return; const scaleFactor = outputSize / PREVIEW_SIZE; ctx.beginPath(); ctx.arc(outputSize / 2, outputSize / 2, outputSize / 2, 0, Math.PI * 2); ctx.clip();
    const { drawWidth, drawHeight } = getDrawSize(img, zoom, PREVIEW_SIZE); const w = drawWidth * scaleFactor; const h = drawHeight * scaleFactor; const x = (outputSize - w) / 2 + position.x * scaleFactor; const y = (outputSize - h) / 2 + position.y * scaleFactor;
    ctx.drawImage(img, x, y, w, h); outputCanvas.toBlob((blob) => { if (blob) { onCropComplete(blob); onOpenChange(false); } }, "image/png", 1.0);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}><DialogContent className="sm:max-w-md" showCloseButton={false} aria-describedby="crop-dialog-description"><DialogHeader><DialogTitle>Adjust your logo</DialogTitle><DialogDescription id="crop-dialog-description">Drag to reposition and use the slider to zoom in or out.</DialogDescription></DialogHeader><div className="flex flex-col items-center gap-6 py-4"><div className="relative rounded-full overflow-hidden bg-muted" style={{ width: PREVIEW_SIZE, height: PREVIEW_SIZE }}><canvas ref={canvasRef} width={PREVIEW_SIZE} height={PREVIEW_SIZE} className={isDragging ? "cursor-grabbing" : "cursor-grab"} onPointerDown={onPointerDown} onPointerMove={onPointerMove} onPointerUp={onPointerUp} onPointerLeave={onPointerUp} aria-label="Image preview. Drag to reposition." role="img" />{imageLoaded && !isDragging && <div className="absolute inset-0 flex items-center justify-center pointer-events-none"><div className="bg-black/40 text-white text-xs px-2 py-1 rounded-full flex items-center gap-1 opacity-70"><Move className="size-3" /><span>Drag to move</span></div></div>}{!imageLoaded && <div className="absolute inset-0 flex items-center justify-center"><span className="text-muted-foreground text-sm">Loading image...</span></div>}</div><div className="flex items-center gap-3 w-full max-w-[280px]"><ZoomOut className="size-4 text-muted-foreground shrink-0" /><Slider value={[zoom]} onValueChange={(v) => setZoom(v[0])} min={MIN_ZOOM} max={MAX_ZOOM} step={0.01} className="flex-1" aria-label="Zoom level" /><ZoomIn className="size-4 text-muted-foreground shrink-0" /></div></div><DialogFooter className="gap-2 sm:gap-0"><Button variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button><Button onClick={handleApply} disabled={!imageLoaded}>Apply</Button></DialogFooter></DialogContent></Dialog>
  );
}

function getDrawSize(img: HTMLImageElement, zoom: number, previewSize: number) {
  const ratio = img.width / img.height;
  if (ratio > 1) { const drawHeight = previewSize * zoom; return { drawWidth: drawHeight * ratio, drawHeight }; }
  const drawWidth = previewSize * zoom; return { drawWidth, drawHeight: drawWidth / ratio };
}
