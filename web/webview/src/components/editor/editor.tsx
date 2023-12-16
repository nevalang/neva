import { useContext, useMemo } from "react";
import { buildGraph } from "./build_graph";
import { FileContext } from "../app";
import { Flow } from "../flow";
import { InterfaceNode } from "../flow/nodes/interface_node";
import { TypeNode } from "../flow/nodes/type_node";
import { ConstNode } from "../flow/nodes/const_node";

const flowNodeTypes = {
  type: TypeNode,
  const: ConstNode,
  interface: InterfaceNode,
  component: InterfaceNode,
};

export function Editor() {
  const fileContext = useContext(FileContext);
  const { nodes, edges } = useMemo(
    () => buildGraph(fileContext.state),
    [fileContext]
  );

  return <Flow nodes={nodes} edges={edges} nodeTypes={flowNodeTypes} />;
}
