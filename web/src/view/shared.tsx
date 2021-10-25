import * as React from "react"
import { useRef } from "react"

interface DraggableProps {
  onDrag(event: MouseEvent)
  onDragEnd(event: MouseEvent)
}

export const Draggable = (props: React.PropsWithChildren<DraggableProps>) => {
  const elemRef = useRef(null)
  const dragProps = useRef(null)

  const initializeDrag = event => {
    const { target, clientX, clientY } = event
    const { offsetTop, offsetLeft } = target
    const { left, top } = elemRef.current.getBoundingClientRect()

    dragProps.current = {
      dragStartLeft: left - offsetLeft,
      dragStartTop: top - offsetTop,
      dragStartX: clientX,
      dragStartY: clientY,
    }

    window.addEventListener("mousemove", drag, false)
    window.addEventListener("mouseup", stopDragging, false)
  }

  const drag = (event: MouseEvent) => {
    elemRef.current.style.transform = `translate(${
      dragProps.current.dragStartLeft +
      event.clientX -
      dragProps.current.dragStartX
    }px, ${
      dragProps.current.dragStartTop +
      event.clientY -
      dragProps.current.dragStartY
    }px)`
    props.onDrag(event)
  }

  const stopDragging = (event: MouseEvent) => {
    window.removeEventListener("mousemove", drag, false)
    window.removeEventListener("mouseup", stopDragging, false)
    props.onDragEnd(event)
  }

  return (
    <div onMouseDown={initializeDrag} ref={elemRef}>
      {props.children}
    </div>
  )
}
