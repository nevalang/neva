import { useMemo } from "react";
import { ComponentViewState } from "../../../core/file_view_state";
import { Flow } from "../../flow";
import { buildComponentNetwork } from "./build_component_net";
import { InterfaceNode } from "../../flow/nodes/interface_node";
import { getLayoutedNodes } from "./get_layouted_nodes";
import { Link } from "react-router-dom";

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
    <div className="entity">
      <Flow
        title={props.entityName}
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        nodesDraggable
        leftTopPanel={<Link to="/">Go Back</Link>}
      />
    </div>
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
