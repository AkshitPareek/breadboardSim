import React from 'react';
import '../css/ComponentProperties.css';

const ComponentProperties = ({ component, onUpdate, onClose }) => {
  if (!component) return null;

  const handleChange = (e) => {
    const { name, value } = e.target;
    onUpdate(component.id, { [name]: value });
  };

  const renderProperties = () => {
    switch (component.type) {
      case 'resistor':
        return (
          <>
            <label>
              Resistance (Ω):
              <input
                type="number"
                name="resistance"
                value={component.properties.resistance || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Power Rating (W):
              <input
                type="number"
                name="powerRating"
                value={component.properties.powerRating || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Tolerance (%):
              <input
                type="number"
                name="tolerance"
                value={component.properties.tolerance || ''}
                onChange={handleChange}
              />
            </label>
          </>
        );
      case 'capacitor':
        return (
          <>
            <label>
              Capacitance (µF):
              <input
                type="number"
                name="capacitance"
                value={component.properties.capacitance || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Voltage Rating (V):
              <input
                type="number"
                name="voltageRating"
                value={component.properties.voltageRating || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Type:
              <select
                name="capacitorType"
                value={component.properties.capacitorType || ''}
                onChange={handleChange}
              >
                <option value="">Select Type</option>
                <option value="electrolytic">Electrolytic</option>
                <option value="ceramic">Ceramic</option>
                <option value="film">Film</option>
              </select>
            </label>
          </>
        );
      case 'inductor':
        return (
          <>
            <label>
              Inductance (H):
              <input
                type="number"
                name="inductance"
                value={component.properties.inductance || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Current Rating (A):
              <input
                type="number"
                name="currentRating"
                value={component.properties.currentRating || ''}
                onChange={handleChange}
              />
            </label>
          </>
        );
      case 'diode':
      case 'led':
        return (
          <>
            <label>
              Forward Voltage (V):
              <input
                type="number"
                name="forwardVoltage"
                value={component.properties.forwardVoltage || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Max Current (mA):
              <input
                type="number"
                name="maxCurrent"
                value={component.properties.maxCurrent || ''}
                onChange={handleChange}
              />
            </label>
            {component.type === 'led' && (
              <label>
                Color:
                <select
                  name="color"
                  value={component.properties.color || ''}
                  onChange={handleChange}
                >
                  <option value="">Select Color</option>
                  <option value="red">Red</option>
                  <option value="green">Green</option>
                  <option value="blue">Blue</option>
                  <option value="yellow">Yellow</option>
                </select>
              </label>
            )}
          </>
        );
      case 'transistor':
        return (
          <>
            <label>
              Type:
              <select
                name="transistorType"
                value={component.properties.transistorType || ''}
                onChange={handleChange}
              >
                <option value="">Select Type</option>
                <option value="npn">NPN</option>
                <option value="pnp">PNP</option>
              </select>
            </label>
            <label>
              Gain (hFE):
              <input
                type="number"
                name="gain"
                value={component.properties.gain || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Max Collector Current (A):
              <input
                type="number"
                name="maxCollectorCurrent"
                value={component.properties.maxCollectorCurrent || ''}
                onChange={handleChange}
              />
            </label>
          </>
        );
      case 'ic':
        return (
          <>
            <label>
              IC Type:
              <input
                type="text"
                name="icType"
                value={component.properties.icType || ''}
                onChange={handleChange}
              />
            </label>
            <label>
              Description:
              <textarea
                name="description"
                value={component.properties.description || ''}
                onChange={handleChange}
              />
            </label>
          </>
        );
      default:
        return null;
    }
  };

  return (
    <div className="component-properties">
      <button className="close-button" onClick={onClose}>&times;</button>
      <h3>{component.type.charAt(0).toUpperCase() + component.type.slice(1)} Properties</h3>
      {renderProperties()}
      {/* Remove the rotation input */}
    </div>
  );
};

export default ComponentProperties;