import React, { useState, useRef } from 'react';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import Breadboard from './components/Breadboard';
import ComponentList from './components/ComponentList';
import CustomDragLayer from './components/CustomDragLayer';
import './css/App.css';

function App() {
  const [state, setState] = useState({
    components: [],
    connections: [],
    customConnectionPoints: []
  });
  const [message, setMessage] = useState(null);
  const fileInputRef = useRef(null);

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
        
        if (!Array.isArray(parsedCircuit.components) || 
            !Array.isArray(parsedCircuit.connections) ||
            !Array.isArray(parsedCircuit.customConnectionPoints)) {
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

  const downloadCircuit = () => {
    const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(state));
    const downloadAnchorNode = document.createElement('a');
    downloadAnchorNode.setAttribute("href", dataStr);
    downloadAnchorNode.setAttribute("download", "circuit.json");
    document.body.appendChild(downloadAnchorNode);
    downloadAnchorNode.click();
    downloadAnchorNode.remove();
    showMessage('Circuit downloaded successfully!');
  };

  const uploadCircuit = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const parsedCircuit = JSON.parse(e.target.result);
          if (!Array.isArray(parsedCircuit.components) || 
              !Array.isArray(parsedCircuit.connections) ||
              !Array.isArray(parsedCircuit.customConnectionPoints)) {
            throw new Error('Invalid circuit file');
          }
          setState(parsedCircuit);
          showMessage('Circuit uploaded successfully!');
        } catch (error) {
          console.error('Error parsing circuit file:', error);
          showMessage('Error uploading circuit. The file may be corrupted or in an invalid format.');
        }
      };
      reader.readAsText(file);
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
        <div className="top-controls">
          <button onClick={saveCircuit}>Save</button>
          <button onClick={loadCircuit}>Load</button>
          <button onClick={downloadCircuit}>Download</button>
          <button onClick={() => fileInputRef.current.click()}>Upload</button>
          <input 
            type="file" 
            ref={fileInputRef} 
            style={{ display: 'none' }} 
            onChange={uploadCircuit} 
            accept=".json"
          />
        </div>
        <ComponentList components={components} />
        <Breadboard 
          state={state} 
          setState={setState} 
        />
        <CustomDragLayer />
        {message && <div className="message-bar">{message}</div>}
      </div>
    </DndProvider>
  );
}

export default App;