import React from 'react';
import ConnectionPoint from './ConnectionPoint';

const CustomConnectionPoint = ({ point, onClick, isActive }) => (
  <ConnectionPoint
    point={point}
    isActive={isActive}
    onClick={() => onClick(point.id)}
  />
);

export default CustomConnectionPoint;