import { useStore } from "reactflow";

export function useZoom() {
  const isZoomMiddle = useStore((s) => s.transform[2] >= 0.6);
  const isZoomClose = useStore((s) => s.transform[2] >= 1);
  return { isZoomMiddle, isZoomClose };
}
