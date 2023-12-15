import { useContext, useEffect, useMemo, useState } from "react";
import { Node } from "reactflow";
import { buildGraph } from "./helpers/build_graph";
import { Flow } from "./flow";
import getLayoutedNodes from "./helpers/get_layouted_nodes";
import { FileContext } from "../app";

export function Editor() {
  const fileContext = useContext(FileContext);
  const { nodes, edges } = useMemo(
    () => buildGraph(fileContext.state),
    [fileContext]
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
