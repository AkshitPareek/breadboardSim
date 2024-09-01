import React, { useRef, useState } from 'react';
import { useDrop } from 'react-dnd';
import Component from './Component';
import Wire from './Wire';
import CustomConnectionPoint from './CustomConnectionPoint';
import ComponentProperties from './ComponentProperties';
import { GRID_SIZE, COMPONENT_WIDTH, COMPONENT_HEIGHT } from '../constants';
import '../css/Breadboard.css';

const COMPONENT_PROPERTIES = {
  resistor: { resistance: 0, powerRating: 0, tolerance: 0 },
  capacitor: { capacitance: 0, voltageRating: 0, capacitorType: 'Electrolytic' },
  inductor: { inductance: 0, currentRating: 0 },
  diode: { forwardVoltage: 0, maxCurrent: 0 },
  led: { forwardVoltage: 0, maxCurrent: 0, color: 'Red' },
  transistor: { transistorType: 'NPN', gain: 0, maxCollectorCurrent: 0 },
  ic: { icType: '', description: '' },
  battery: { voltage: 0, capacity: 0 },
  power_supply: { voltage: 0, maxCurrent: 0 }
};

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
          moveComponent(item.id, x, y);
        } else {
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
      properties: { ...COMPONENT_PROPERTIES[type] }, // Initialize properties based on type
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

  const handleComponentClick = (component) => {
    setSelectedComponent(component);
    setIsPropertiesPanelOpen(true);
  };

  const closePropertiesPanel = () => {
    setIsPropertiesPanelOpen(false);
    setSelectedComponent(null);
  };

  const updateComponentProperties = (id, newProperties) => {
    console.log(`Updating properties for component ${id}:`, newProperties);
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

  return (
    <div className="breadboard-layout">
      <div className="breadboard-wrapper">
        <div 
          ref={(node) => {
            drop(node);
            boardRef.current = node;
          }}
          className="breadboard-container"
          onClick={handleGridClick}
        >
          <svg width="800" height="400" style={{position: 'absolute', top: 0, left: 0, pointerEvents: 'none'}}>
            {state.connections.map((connection, index) => (
              <Wire key={index} connection={connection} components={state.components} customPoints={state.customConnectionPoints} />
            ))}
          </svg>
          {state.customConnectionPoints.map(point => (
            <CustomConnectionPoint
              key={point.id}
              point={point}
              onClick={() => handleCustomConnectionPointClick(point.id)}
              isActive={wireStart && wireStart.customPointId === point.id}
            />
          ))}
          {state.components.map(component => (
            <div
              key={component.id}
              style={{
                position: 'absolute',
                left: `${component.position.x * GRID_SIZE}px`,
                top: `${component.position.y * GRID_SIZE}px`,
              }}
            >
              <Component
                component={component}
                onMove={handleConnectionPointClick}
                activeWireStart={wireStart}
                onClick={() => handleComponentClick(component)}
              />
            </div>
          ))}
        </div>
        {isPropertiesPanelOpen && (
          <ComponentProperties
            component={selectedComponent}
            onUpdate={updateComponentProperties}
            onClose={closePropertiesPanel}
          />
        )}
      </div>
    </div>
  );
};

export default Breadboard;