import React, { useState } from 'react';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import Breadboard from './components/Breadboard';
import ComponentList from './components/ComponentList';
import CustomDragLayer from './components/CustomDragLayer';
import './App.css';

function App() {
  const [state, setState] = useState({
    components: [],
    connections: []
  });

  const components = ['resistor', 'capacitor', 'inductor', 'diode', 'led', 'transistor', 'ic'];

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="App">
        <h1>Breadboard Simulator</h1>
        <ComponentList components={components} />
        <Breadboard state={state} setState={setState} />
        <CustomDragLayer />
      </div>
    </DndProvider>
  );
}

export default App;