import { useState } from "react";
import { ModelService } from "@wailsbindings/services";

export function useModelUpload() {
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const openChooseFileDialog = async () => {
    setIsUploading(true);
    setError(null);

    try {
      await ModelService.ChooseModel();
      setIsUploading(false);
      return { success: true };
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Upload failed";
      setError(errorMessage);
      setIsUploading(false);
      return { success: false, error: errorMessage };
    }
  };

  const resetError = () => setError(null);

  return {
    isUploading,
    error,
    resetError,
    handleFileUpload: openChooseFileDialog,
  };
}
