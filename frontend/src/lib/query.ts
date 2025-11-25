import { useEffect, useState, useCallback } from "react";

type UseQueryOptions<TData, TError = Error> = {
  queryKey: unknown[];
  queryFn: () => Promise<TData>;
  enabled?: boolean;
};

type UseQueryResult<TData, TError> = {
  data: TData | undefined;
  error: TError | undefined;
  isLoading: boolean;
  refetch: () => Promise<void>;
};

export function useQuery<TData, TError = unknown>(
  options: UseQueryOptions<TData, TError>
): UseQueryResult<TData, TError> {
  const { queryKey, queryFn, enabled = true } = options;

  const [data, setData] = useState<TData>();
  const [error, setError] = useState<TError>();
  const [isLoading, setLoading] = useState<boolean>(enabled);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(undefined);

    try {
      const result = await queryFn();
      setData(result);
    } catch (err) {
      setError(err as TError);
    } finally {
      setLoading(false);
    }
  }, [queryFn, JSON.stringify(queryKey)]);

  useEffect(() => {
    if (!enabled) return;
    fetchData();
  }, [enabled, fetchData]);

  return { data, error, isLoading, refetch: fetchData };
}
