import React from 'react';
import { useDrag } from 'react-dnd';
import { getComponentIcon } from './ComponentIcons';
import './ComponentList.css';

const DOT_SIZE = 10;

export const DraggableComponent = ({ component, onMove, activeWireStart, onRotate, onClick }) => {
  const [{ isDragging }, drag] = useDrag(() => ({
    type: 'component',
    item: () => {
      if (typeof component === 'string') {
        // This is a new component from the list
        return { type: component };
      } else {
        // This is an existing component on the board
        return { id: component.id, type: component.type };
      }
    },
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
        position: component.position ? 'absolute' : 'relative',
        left: component.position ? `${component.position.x * 20}px` : 'auto',
        top: component.position ? `${component.position.y * 20}px` : 'auto',
        transform: component.rotation ? `rotate(${component.rotation}deg)` : 'none',
      }}
      onClick={() => onClick && onClick(component)}
      onDoubleClick={() => onRotate && onRotate(component.id)}
    >
      <div className="component-icon">
        {getComponentIcon(component.type || component)}
      </div>
      {!component.position && <div className="component-label">{component.type || component}</div>}
      {component.connectionPoints && component.connectionPoints.map((point, index) => (
        <div
          key={index}
          className={`connection-point ${activeWireStart && activeWireStart.componentId === component.id && activeWireStart.pointIndex === index ? 'active' : ''}`}
          style={{
            position: 'absolute',
            left: `${point.x - DOT_SIZE/2}px`,
            top: `${point.y - DOT_SIZE/2}px`,
            width: `${DOT_SIZE}px`,
            height: `${DOT_SIZE}px`,
            borderRadius: '50%',
            backgroundColor: activeWireStart && activeWireStart.componentId === component.id && activeWireStart.pointIndex === index ? 'green' : 'red',
            cursor: 'pointer',
            transition: 'background-color 0.3s ease',
          }}
          onClick={(e) => {
            e.stopPropagation();
            onMove && onMove(component.id, index);
          }}
        />
      ))}
    </div>
  );
};

const ComponentList = ({ components }) => {
  return (
    <div className="component-list">
      {components.map((component, index) => (
        <DraggableComponent key={index} component={component} />
      ))}
    </div>
  );
};

export default ComponentList;