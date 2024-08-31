import React from 'react';
import { DOT_SIZE } from '../constants';

const ConnectionPoint = ({ point, index, isActive, onClick }) => (
  <div
    className={`connection-point ${isActive ? 'active' : ''}`}
    style={{
      position: 'absolute',
      left: `${point.x + DOT_SIZE/2}px`,
      top: `${point.y - DOT_SIZE/2}px`,
      width: `${DOT_SIZE}px`,
      height: `${DOT_SIZE}px`,
      borderRadius: '50%',
      backgroundColor: isActive ? 'green' : 'red',
      cursor: 'pointer',
      transition: 'background-color 0.3s ease',
    }}
    onClick={(e) => {
      e.stopPropagation();
      onClick(index);
    }}
  />
);

export default ConnectionPoint;