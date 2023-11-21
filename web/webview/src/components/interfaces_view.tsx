import { Interface } from "../generated/sourcecode";
import { InterfaceView } from "./interface_view";

interface IInterfacesViewProps {
  interfaces: Array<{ name: string; entity: Interface }>;
}

export function InterfacesView(props: IInterfacesViewProps) {
  return (
    <>
      {props.interfaces.map((entry) => {
        const { name, entity } = entry;
        return <InterfaceView name={name} entity={entity} />;
      })}
    </>
  );
}
