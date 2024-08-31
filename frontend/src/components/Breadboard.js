import React, { useRef, useState } from 'react';
import { useDrop } from 'react-dnd';
import { DraggableComponent } from './ComponentList';
import ComponentProperties from './ComponentProperties';
import { getComponentIcon } from './ComponentIcons';
import './Breadboard.css';

const GRID_SIZE = 20;
const COMPONENT_WIDTH = 66;
const COMPONENT_HEIGHT = 50;
const DOT_SIZE = 10;

const Breadboard = ({ state, setState }) => {
  const boardRef = useRef(null);
  const [wireStart, setWireStart] = useState(null);
  const [selectedComponent, setSelectedComponent] = useState(null);
  const [isPropertiesPanelOpen, setIsPropertiesPanelOpen] = useState(false);

  const [, drop] = useDrop(() => ({
    accept: 'component',
    drop: (item, monitor) => {
      const boardRect = boardRef.current.getBoundingClientRect();
      const clientOffset = monitor.getClientOffset();
      const initialClientOffset = monitor.getInitialClientOffset();
      const initialSourceClientOffset = monitor.getInitialSourceClientOffset();
      
      if (clientOffset && initialClientOffset && initialSourceClientOffset) {
        const dx = clientOffset.x - initialClientOffset.x;
        const dy = clientOffset.y - initialClientOffset.y;

        const x = Math.floor((initialSourceClientOffset.x + dx - boardRect.left) / GRID_SIZE);
        const y = Math.floor((initialSourceClientOffset.y + dy - boardRect.top) / GRID_SIZE);

        if (item.id) {
          // Move existing component
          moveComponent(item.id, x, y);
        } else {
          // Add new component
          addComponent(item.type, x, y);
        }
      }
    },
  }));

  const addComponent = (type, x, y) => {
    const newComponent = {
      id: `${type}-${Date.now()}`,
      type,
      position: { x, y },
      rotation: 0,
      properties: {},
      connectionPoints: [
        { x: 0, y: COMPONENT_HEIGHT / 2 },
        { x: COMPONENT_WIDTH, y: COMPONENT_HEIGHT / 2 },
      ],
    };
    setState(prevState => ({
      ...prevState,
      components: [...prevState.components, newComponent],
    }));
  };

  const moveComponent = (id, x, y) => {
    setState(prevState => ({
      ...prevState,
      components: prevState.components.map(comp =>
        comp.id === id ? { ...comp, position: { x, y } } : comp
      ),
    }));
  };

  const handleConnectionPointClick = (componentId, pointIndex) => {
    if (!wireStart) {
      setWireStart({ componentId, pointIndex });
    } else {
      if (wireStart.componentId !== componentId) {
        setState(prevState => ({
          ...prevState,
          connections: [
            ...prevState.connections,
            { from: wireStart, to: { componentId, pointIndex } }
          ]
        }));
      }
      setWireStart(null);
    }
  };

  const rotateComponent = (id) => {
    // Implementation of rotateComponent function
  };

  const handleComponentClick = (component) => {
    setSelectedComponent(component);
    setIsPropertiesPanelOpen(true);
  };

  const closePropertiesPanel = () => {
    setIsPropertiesPanelOpen(false);
    setSelectedComponent(null);
  };

  const updateComponentProperties = (id, newProperties) => {
    // Implementation of updateComponentProperties function
  };

  const renderWires = () => {
    return state.connections.map((connection, index) => {
      const startComponent = state.components.find(c => c.id === connection.from.componentId);
      const endComponent = state.components.find(c => c.id === connection.to.componentId);
      
      if (!startComponent || !endComponent) return null;

      const startPoint = startComponent.connectionPoints[connection.from.pointIndex];
      const endPoint = endComponent.connectionPoints[connection.to.pointIndex];

      const start = {
        x: startComponent.position.x * GRID_SIZE + startPoint.x + DOT_SIZE,
        y: startComponent.position.y * GRID_SIZE + startPoint.y,
      };
      const end = {
        x: endComponent.position.x * GRID_SIZE + endPoint.x + DOT_SIZE,
        y: endComponent.position.y * GRID_SIZE + endPoint.y,
      };

      return (
        <line
          key={index}
          x1={start.x}
          y1={start.y}
          x2={end.x}
          y2={end.y}
          stroke="black"
          strokeWidth="2"
        />
      );
    });
  };

  return (
    <div className="breadboard-layout">
      <div className="breadboard-wrapper">
        <div 
          ref={(node) => {
            drop(node);
            boardRef.current = node;
          }}
          className="breadboard-container"
        >
          <svg width="800" height="400" style={{position: 'absolute', top: 0, left: 0, pointerEvents: 'none'}}>
            {renderWires()}
          </svg>
          {(state.components || []).map(component => (
            <DraggableComponent
              key={component.id}
              component={component}
              onMove={handleConnectionPointClick}
              activeWireStart={wireStart}
              onRotate={rotateComponent}
              onClick={() => handleComponentClick(component)}
            />
          ))}
        </div>
      </div>
      {isPropertiesPanelOpen && (
        <div className="properties-panel">
          <button className="close-button" onClick={closePropertiesPanel}>×</button>
          <ComponentProperties
            component={selectedComponent}
            onUpdate={updateComponentProperties}
          />
        </div>
      )}
    </div>
  );
};

export default Breadboard;