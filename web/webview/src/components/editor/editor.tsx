import { useEffect, useMemo, useState } from "react";
import { Node } from "reactflow";
import { FileViewState } from "../../core/file_view_state";
import { buildGraph } from "./helpers/build_graph";
import { Flow } from "./flow";
import getLayoutedNodes from "./helpers/get_layouted_nodes";

interface IEditorProps {
  fileViewState: FileViewState;
}

export function Editor(props: IEditorProps) {
  const { nodes, edges } = useMemo(
    () => buildGraph(props.fileViewState),
    [props.fileViewState]
  );

  const [layoutedNodes, setLayoutedNodes] = useState<Node[]>([]);

  useEffect(() => {
    (async () => {
      const newLayoutedNodes = await getLayoutedNodes(nodes, edges);
      setLayoutedNodes(newLayoutedNodes);
    })();
  }, [edges, nodes]);

  if (layoutedNodes.length === 0) {
    return null;
  }

  return <Flow nodes={layoutedNodes} edges={edges} />;
}
