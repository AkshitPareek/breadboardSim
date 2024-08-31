import React from 'react';
import { useDrag } from 'react-dnd';
import { getComponentIcon } from './ComponentIcons';
import '../css/ComponentList.css';

const DraggableComponent = ({ component }) => {
  const [{ isDragging }, drag] = useDrag(() => ({
    type: 'component',
    item: { type: component },
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
      }}
    >
      <div className="component-icon">
        {getComponentIcon(component)}
      </div>
      <div className="component-label">{component}</div>
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