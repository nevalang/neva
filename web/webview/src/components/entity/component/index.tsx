import { useMemo } from "react";
import { ComponentViewState } from "../../../core/file_view_state";
import { Flow } from "../../flow";33
import { buildComponentNetwork } from "./network/helpers";
import { InterfaceNode } from "../../flow/nodes/interface_node";

interface IComponentProps {
  viewState: ComponentViewState;
  entityName: string;
}

const nodeTypes = { component: InterfaceNode };

export function Component(props: IComponentProps) {
  const { nodes, edges } = useMemo(
    () => buildComponentNetwork(props.entityName, props.viewState),
    [props.entityName, props.viewState]
  );
  return <Flow nodes={nodes} edges={edges} nodeTypes={nodeTypes} />;
}
