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
  const [message, setMessage] = useState(null);

  const components = ['resistor', 'capacitor', 'inductor', 'diode', 'led', 'transistor', 'ic'];

  const saveCircuit = () => {
    localStorage.setItem('savedCircuit', JSON.stringify(state));
    showMessage('Circuit saved successfully!');
  };

  const loadCircuit = () => {
    try {
      const savedCircuit = localStorage.getItem('savedCircuit');
      if (savedCircuit) {
        const parsedCircuit = JSON.parse(savedCircuit);
        
        if (!Array.isArray(parsedCircuit.components) || !Array.isArray(parsedCircuit.connections)) {
          throw new Error('Invalid saved circuit data');
        }
        
        setState(parsedCircuit);
        showMessage('Circuit loaded successfully!');
      } else {
        showMessage('No saved circuit found!');
      }
    } catch (error) {
      console.error('Error loading circuit:', error);
      showMessage('Error loading circuit. The saved data may be corrupted.');
    }
  };

  const showMessage = (msg) => {
    setMessage(msg);
    setTimeout(() => setMessage(null), 3000);
  };

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="App">
        <h1>Breadboard Simulator</h1>
        <ComponentList components={components} />
        <Breadboard 
          state={state} 
          setState={setState} 
          onSave={saveCircuit} 
          onLoad={loadCircuit}
        />
        <CustomDragLayer />
        {message && <div className="message-bar">{message}</div>}
      </div>
    </DndProvider>
  );
}

export default App;