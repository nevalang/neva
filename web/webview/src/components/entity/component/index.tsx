import { useMemo } from "react";
import { ComponentViewState } from "../../../core/file_view_state";
import { Flow } from "../../flow";
import { buildComponentNetwork } from "./build_component_net";
import { InterfaceNode } from "../../flow/nodes/interface_node";
import { getLayoutedNodes } from "./get_layouted_nodes";

interface IComponentProps {
  viewState: ComponentViewState;
  entityName: string;
}

const nodeTypes = {
  interface: InterfaceNode,
  component: InterfaceNode,
};

export function Component(props: IComponentProps) {
  const { nodes, edges } = useMemo(
    () =>
      getLayoutedNodes(
        buildComponentNetwork(props.entityName, props.viewState)
      ),
    [props.entityName, props.viewState]
  );

  return (
    <Flow nodes={nodes} edges={edges} nodeTypes={nodeTypes} nodesDraggable />
  );
}

// const onNodeMouseEnter = useCallback(
//     (_: MouseEvent, hoveredNode: Node) => {
//       if (hoveredNode.type !== "component") {
//         return;
//       }
//       const { newEdges, newNodes } = handleNodeMouseEnter(
//         hoveredNode,
//         edgesState,
//         nodesState
//       );
//       setEdgesState(newEdges);
//       setNodesState(newNodes);
//     },
//     [edgesState, nodesState, setEdgesState, setNodesState]
//   );

//   const onNodeMouseLeave = useCallback(() => {
//     const { newEdges, newNodes } = handleNodeMouseLeave(edgesState, nodesState);
//     setEdgesState(newEdges);
//     setNodesState(newNodes);
//   }, [edgesState, nodesState, setEdgesState, setNodesState]);
