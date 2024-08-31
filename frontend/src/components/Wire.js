import React from 'react';
import { GRID_SIZE, DOT_SIZE } from '../constants';

const Wire = ({ connection, components, customPoints }) => {
  const getPointCoordinates = (pointInfo) => {
    if (pointInfo.customPointId) {
      const customPoint = customPoints.find(p => p.id === pointInfo.customPointId);
      return customPoint ? { x: customPoint.x + DOT_SIZE, y: customPoint.y} : null;
    } else {
      const component = components.find(c => c.id === pointInfo.componentId);
      if (!component) return null;
      const point = component.connectionPoints[pointInfo.pointIndex];
      if (!point) return null;
      return {
        x: component.position.x * GRID_SIZE + point.x + DOT_SIZE,
        y: component.position.y * GRID_SIZE + point.y
      };
    }
  };

  const start = getPointCoordinates(connection.from);
  const end = getPointCoordinates(connection.to);

  if (!start || !end) return null;

  return (
    <line
      x1={start.x}
      y1={start.y}
      x2={end.x}
      y2={end.y}
      stroke="black"
      strokeWidth="2"
    />
  );
};

export default Wire;