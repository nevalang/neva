import { useMemo } from "react";
import { Link } from "react-router-dom";
import { ComponentViewState } from "../../../core/file_view_state";
import { Flow } from "../../flow";
import { InterfaceNode } from "../../flow/nodes/interface_node";
import { buildComponentNetwork } from "./build_component_net";
import { getLayoutedNodes } from "./get_layouted_nodes";
import { handleNodeMouseEnter, handleNodeMouseLeave } from "./mouse_handlers";
import { NormalEdge } from "./edge";

interface IComponentProps {
  viewState: ComponentViewState;
  entityName: string;
}

const nodeTypes = {
  interface: InterfaceNode,
  component: InterfaceNode,
};

const edgeTypes = {
  normal: NormalEdge,
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
        edgeTypes={edgeTypes}
        nodesDraggable
        leftTopPanel={<Link to="/">Back</Link>}
        onNodeMouseEnter={handleNodeMouseEnter}
        onNodeMouseLeave={handleNodeMouseLeave}
      />
    </div>
  );
}
