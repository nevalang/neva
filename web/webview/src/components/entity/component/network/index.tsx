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