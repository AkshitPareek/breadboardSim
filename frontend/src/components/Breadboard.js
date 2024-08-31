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
    setState(prevState => ({
      ...prevState,
      components: prevState.components.map(comp =>
        comp.id === id
          ? { ...comp, rotation: ((comp.rotation || 0) + 90) % 360 }
          : comp
      ),
    }));
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
    setState(prevState => ({
      ...prevState,
      components: prevState.components.map(comp =>
        comp.id === id
          ? { ...comp, properties: { ...comp.properties, ...newProperties } }
          : comp
      ),
    }));
  };

  const handleGridClick = (e) => {
    if (e.metaKey || e.ctrlKey) {
      const rect = boardRef.current.getBoundingClientRect();
      const x = Math.floor((e.clientX - rect.left) / GRID_SIZE) * GRID_SIZE;
      const y = Math.floor((e.clientY - rect.top) / GRID_SIZE) * GRID_SIZE;
      
      setState(prevState => ({
        ...prevState,
        customConnectionPoints: [...prevState.customConnectionPoints, { x, y, id: Date.now() }]
      }));
    }
  };

  const handleCustomConnectionPointClick = (pointId) => {
    if (!wireStart) {
      setWireStart({ customPointId: pointId });
    } else {
      if (wireStart.customPointId !== pointId) {
        setState(prevState => ({
          ...prevState,
          connections: [
            ...prevState.connections,
            { from: wireStart, to: { customPointId: pointId } }
          ]
        }));
      }
      setWireStart(null);
    }
  };

  const renderCustomConnectionPoints = () => {
    return state.customConnectionPoints.map((point) => (
      <div
        key={point.id}
        className={`connection-point ${wireStart && wireStart.customPointId === point.id ? 'active' : ''}`}
        style={{
          position: 'absolute',
          left: `${point.x}px`,
          top: `${point.y}px`,
          width: `${DOT_SIZE}px`,
          height: `${DOT_SIZE}px`,
        }}
        onClick={() => handleCustomConnectionPointClick(point.id)}
      />
    ));
  };

  const renderWires = () => {
    return state.connections.map((connection, index) => {
      let start, end;

      try {
        if (connection.from.customPointId) {
          const startPoint = state.customConnectionPoints.find(p => p.id === connection.from.customPointId);
          if (!startPoint) throw new Error('Start point not found');
          start = { x: startPoint.x + DOT_SIZE / 2, y: startPoint.y + DOT_SIZE / 2 };
        } else {
          const startComponent = state.components.find(c => c.id === connection.from.componentId);
          if (!startComponent) throw new Error('Start component not found');
          const startPoint = startComponent.connectionPoints[connection.from.pointIndex];
          if (!startPoint) throw new Error('Start connection point not found');
          start = {
            x: startComponent.position.x * GRID_SIZE + startPoint.x + DOT_SIZE,
            y: startComponent.position.y * GRID_SIZE + startPoint.y,
          };
        }

        if (connection.to.customPointId) {
          const endPoint = state.customConnectionPoints.find(p => p.id === connection.to.customPointId);
          if (!endPoint) throw new Error('End point not found');
          end = { x: endPoint.x + DOT_SIZE / 2, y: endPoint.y + DOT_SIZE / 2 };
        } else {
          const endComponent = state.components.find(c => c.id === connection.to.componentId);
          if (!endComponent) throw new Error('End component not found');
          const endPoint = endComponent.connectionPoints[connection.to.pointIndex];
          if (!endPoint) throw new Error('End connection point not found');
          end = {
            x: endComponent.position.x * GRID_SIZE + endPoint.x + DOT_SIZE,
            y: endComponent.position.y * GRID_SIZE + endPoint.y,
          };
        }

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
      } catch (error) {
        console.error(`Error rendering wire ${index}:`, error.message);
        return null; // Skip rendering this wire
      }
    }).filter(Boolean); // Remove any null entries (failed renders)
  };

  return (
    <div className="breadboard-layout">
      <div className="breadboard-wrapper">
        {/* Remove this section if it exists
        <div className="breadboard-controls">
          <button onClick={() => {}}>Save</button>
          <button onClick={() => {}}>Load</button>
        </div>
        */}
        <div 
          ref={(node) => {
            drop(node);
            boardRef.current = node;
          }}
          className="breadboard-container"
          onClick={handleGridClick}
        >
          <svg width="800" height="400" style={{position: 'absolute', top: 0, left: 0, pointerEvents: 'none'}}>
            {renderWires()}
          </svg>
          {renderCustomConnectionPoints()}
          {(state.components || []).map(component => (
            <DraggableComponent
              key={component.id}
              component={component}
              onMove={handleConnectionPointClick}
              activeWireStart={wireStart}
              onRotate={() => rotateComponent(component.id)}
              onClick={() => handleComponentClick(component)}
            />
          ))}
        </div>
        {isPropertiesPanelOpen && (
          <div className="properties-panel">
            <ComponentProperties
              component={selectedComponent}
              onUpdate={updateComponentProperties}
              onClose={closePropertiesPanel}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default Breadboard;