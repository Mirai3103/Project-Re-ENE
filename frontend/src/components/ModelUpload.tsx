import React, { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { RippleButton } from "@/components/ui/shadcn-io/ripple-button";
import { Upload, FileArchive, Loader2 } from "lucide-react";
import { useModelUpload } from "@/hooks/useModelUpload";
import { Events } from "@wailsio/runtime";
import { ModelService } from "@wailsbindings/services";
interface ModelUploadProps {
  onUploadSuccess?: () => void;
}

export default function ModelUpload({ onUploadSuccess }: ModelUploadProps) {
  const { isUploading, error, handleFileUpload, resetError } = useModelUpload();
  const [isDragging, setIsDragging] = useState(false);

  // Setup Wails drag and drop
  useEffect(() => {
    const handleFileDrop = async (event: any) => {
      setIsDragging(false);
      const file = event.data.files[0];
      if (!file) return;
      console.log("Dropped file:", file);
      return;
      const filePath = file.path;

      const isZip = filePath.toLowerCase().endsWith(".zip");

      if (!isZip) {
        alert("Please upload a .zip file containing the Live2D model");
        return;
      }

      await ModelService.UploadModel(filePath);
      onUploadSuccess && onUploadSuccess?.();
    };

    // Register drop handler with useDropTarget = true for CSS classes
    Events.On("common:WindowFilesDropped", handleFileDrop);

    return () => {
      Events.Off("common:WindowFilesDropped");
    };
  }, [handleFileUpload, onUploadSuccess]);

  // Handle errors
  useEffect(() => {
    if (error) {
      alert(error);
      resetError();
    }
  }, [error, resetError]);

  const handleBrowseClick = async () => {
    const result = await handleFileUpload();
    if (result.success && onUploadSuccess) {
      setTimeout(onUploadSuccess, 500);
    }
  };

  return (
    <Card className="border-2 border-primary/20 bg-card/50 backdrop-blur-sm relative overflow-hidden">
      <div className="absolute -right-12 -top-12 -z-10 h-32 w-32 rounded-full bg-primary/20 blur-3xl" />

      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-2xl">
          <Upload className="h-6 w-6 text-primary" />
          Upload Model
        </CardTitle>
        <CardDescription>
          Upload a .zip file containing your Live2D model
        </CardDescription>
      </CardHeader>

      <CardContent>
        <div
          className={`
            relative border-2 border-dashed rounded-2xl p-8
            transition-all duration-300
            ${
              isDragging || isUploading
                ? "border-primary bg-primary/10 scale-105"
                : "border-primary/30 bg-background/30 hover:border-primary/60"
            }
          `}
          style={{ "--wails-drop-target": "drop" } as React.CSSProperties}
          onDragEnter={() => setIsDragging(true)}
          onDragLeave={() => setIsDragging(false)}
        >
          <div className="flex flex-col items-center gap-4">
            <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
              {isUploading ? (
                <Loader2 className="h-8 w-8 text-primary animate-spin" />
              ) : (
                <FileArchive className="h-8 w-8 text-primary" />
              )}
            </div>

            <div className="text-center space-y-2">
              <p className="text-lg font-semibold text-foreground">
                {isUploading
                  ? "Uploading..."
                  : isDragging
                  ? "Drop your file here!"
                  : "Click to upload or drag & drop"}
              </p>
              <p className="text-sm text-muted-foreground">
                .ZIP files only (Max 100MB)
              </p>
            </div>

            <RippleButton
              type="button"
              className="mt-2"
              onClick={handleBrowseClick}
              disabled={isUploading}
            >
              {isUploading ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Uploading...
                </>
              ) : (
                <>
                  <Upload className="h-4 w-4 mr-2" />
                  Browse Files
                </>
              )}
            </RippleButton>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
