import { useCallback, useContext, useEffect, useMemo, useState } from "react";
import { Node } from "reactflow";
import { buildFileNodes } from "./build_file_nodes";
import { FileContext } from "../app";
import { Flow } from "../flow";
import { InterfaceNode } from "../flow/nodes/interface_node";
import { TypeNode } from "../flow/nodes/type_node";
import { ConstNode } from "../flow/nodes/const_node";
import { getLayoutedNodes } from "./get_layouted_nodes";
import { useNavigate } from "react-router-dom";

const flowNodeTypes = {
  type: TypeNode,
  const: ConstNode,
  interface: InterfaceNode,
  component: InterfaceNode,
};

export function Editor() {
  const navigate = useNavigate();
  const fileContext = useContext(FileContext);

  const nodes = useMemo(
    () => buildFileNodes(fileContext.state),
    [fileContext.state]
  );

  const [layoutedNodes, setLayoutedNodes] = useState<Node[]>([]);
  useEffect(() => {
    (async () => {
      setLayoutedNodes(await getLayoutedNodes(nodes));
    })();
  }, [nodes]);

  const handleNodeClick = useCallback(
    (node: Node) => {
      console.log({ node });
      if (node.type === "component") {
        navigate(`/${node.data.entityName}`);
      }
    },
    [navigate]
  );

  return (
    <div className="editor">
      <Flow
        nodes={layoutedNodes}
        edges={[]}
        nodeTypes={flowNodeTypes}
        onNodeClick={handleNodeClick}
        title=""
      />
    </div>
  );
}
