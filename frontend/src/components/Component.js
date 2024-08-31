import React from 'react';
import { useDrag } from 'react-dnd';
import { getComponentIcon } from './ComponentIcons';
import ConnectionPoint from './ConnectionPoint';
import { COMPONENT_WIDTH, COMPONENT_HEIGHT } from '../constants';

const Component = ({ component, onMove, activeWireStart, onClick }) => {
  const [{ isDragging }, drag] = useDrag(() => ({
    type: 'component',
    item: { id: component.id, type: component.type },
    collect: (monitor) => ({
      isDragging: !!monitor.isDragging(),
    }),
  }));

  return (
    <div
      ref={drag}
      className="component-item"
      style={{ 
        opacity: isDragging ? 0.5 : 1,
        cursor: 'move',
        width: `${COMPONENT_WIDTH}px`,
        height: `${COMPONENT_HEIGHT}px`,
      }}
      onClick={() => onClick && onClick(component)}
    >
      <div className="component-icon">
        {getComponentIcon(component.type)}
      </div>
      {component.connectionPoints && component.connectionPoints.map((point, index) => (
        <ConnectionPoint
          key={index}
          point={point}
          index={index}
          isActive={activeWireStart && activeWireStart.componentId === component.id && activeWireStart.pointIndex === index}
          onClick={() => onMove && onMove(component.id, index)}
        />
      ))}
    </div>
  );
};

export default Component;