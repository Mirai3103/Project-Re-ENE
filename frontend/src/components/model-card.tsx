import { ModelService, Model as Live2dModel } from "@wailsbindings/services";
import { Badge, CheckCircle, Download } from "lucide-react";
import { useQuery } from "@/lib/query";
import type { ILive2DModel } from "@/types/models";
interface ModelCardProps {
  model: Live2dModel;
  setSelectedModel: (model: Live2dModel & { data: ILive2DModel }) => void;
  selectedModel: Live2dModel | undefined;
}
export default function ModelCard({
  model,
  setSelectedModel,
  selectedModel,
}: ModelCardProps) {
  const { data, error, isLoading } = useQuery({
    queryKey: ["models", model.id],
    queryFn: async () => {
      const res = await fetch("/models" + model.path);
      const data = await res.json();
      return data as ILive2DModel;
    },
  });

  return (
    <div
      key={model.id}
      onClick={() => setSelectedModel({ ...model, data: data as ILive2DModel })}
      className={`
      group relative p-4 rounded-xl cursor-pointer
      transition-all duration-300
      ${
        selectedModel?.id === model.id
          ? "bg-primary/10 border-2 border-primary shadow-lg"
          : "bg-background/50 border-2 border-transparent hover:border-primary/30 hover:bg-background/70"
      }
    `}
    >
      <div className="flex items-start gap-4">
        {/* Thumbnail */}

        {/* Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2 mb-1">
            <h3 className="font-semibold text-foreground truncate">
              {model.name}
            </h3>
            {model.is_active && (
              <Badge className="bg-primary/20 text-primary border-primary/30 shrink-0">
                <CheckCircle className="h-3 w-3 mr-1" />
                Active
              </Badge>
            )}
          </div>

          <div className="flex items-center gap-3 text-xs text-muted-foreground">
            <span>{Math.round((model.size ?? 0) / 1024 / 1024)} MB</span>
            <span>•</span>
            <span>{data?.Groups?.length} motions</span>
            <span>•</span>
            <span>{data?.FileReferences?.Expressions?.length} expressions</span>
          </div>
        </div>
      </div>
    </div>
  );
}
