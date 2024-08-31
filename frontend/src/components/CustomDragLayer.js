import React from 'react';
import { useDragLayer } from 'react-dnd';
import { getComponentIcon } from './ComponentIcons';

const CustomDragLayer = () => {
  const { isDragging, item, currentOffset } = useDragLayer((monitor) => ({
    item: monitor.getItem(),
    currentOffset: monitor.getSourceClientOffset(),
    isDragging: monitor.isDragging(),
  }));

  if (!isDragging) {
    return null;
  }

  return (
    <div style={{
      position: 'fixed',
      pointerEvents: 'none',
      zIndex: 100,
      left: 0,
      top: 0,
      width: '100%',
      height: '100%',
    }}>
      <div style={{
        position: 'absolute',
        width: '66px',
        height: '50px',
        left: currentOffset?.x,
        top: currentOffset?.y,
        opacity: 0.8,
      }}>
        {getComponentIcon(item.type)}
      </div>
    </div>
  );
};

export default CustomDragLayer;