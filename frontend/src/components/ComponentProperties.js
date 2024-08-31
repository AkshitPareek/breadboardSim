import React from 'react';

const ComponentProperties = ({ component, onUpdate }) => {
  if (!component) return null;

  const handleChange = (e) => {
    const { name, value } = e.target;
    onUpdate(component.id, { ...component.properties, [name]: parseFloat(value) });
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
                value={component.properties.resistance}
                onChange={handleChange}
              />
            </label>
            <label>
              Power Rating (W):
              <input
                type="number"
                name="powerRating"
                value={component.properties.powerRating || 0.25}
                onChange={handleChange}
              />
            </label>
            <label>
              Tolerance (%):
              <input
                type="number"
                name="tolerance"
                value={component.properties.tolerance || 5}
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
                value={component.properties.capacitance}
                onChange={handleChange}
              />
            </label>
            <label>
              Voltage Rating (V):
              <input
                type="number"
                name="voltageRating"
                value={component.properties.voltageRating || 50}
                onChange={handleChange}
              />
            </label>
            <label>
              Type:
              <select
                name="type"
                value={component.properties.type || 'ceramic'}
                onChange={handleChange}
              >
                <option value="ceramic">Ceramic</option>
                <option value="electrolytic">Electrolytic</option>
                <option value="film">Film</option>
              </select>
            </label>
          </>
        );
      case 'inductor':
        return (
          <>
            <label>
              Inductance (mH):
              <input
                type="number"
                name="inductance"
                value={component.properties.inductance}
                onChange={handleChange}
              />
            </label>
            <label>
              Current Rating (A):
              <input
                type="number"
                name="currentRating"
                value={component.properties.currentRating || 1}
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
                value={component.properties.forwardVoltage}
                onChange={handleChange}
              />
            </label>
            <label>
              Max Current (mA):
              <input
                type="number"
                name="maxCurrent"
                value={component.properties.maxCurrent || 20}
                onChange={handleChange}
              />
            </label>
            {component.type === 'led' && (
              <label>
                Color:
                <select
                  name="color"
                  value={component.properties.color || 'red'}
                  onChange={handleChange}
                >
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
                value={component.properties.transistorType || 'npn'}
                onChange={handleChange}
              >
                <option value="npn">NPN</option>
                <option value="pnp">PNP</option>
              </select>
            </label>
            <label>
              Gain (hFE):
              <input
                type="number"
                name="gain"
                value={component.properties.gain}
                onChange={handleChange}
              />
            </label>
            <label>
              Max Collector Current (A):
              <input
                type="number"
                name="maxCollectorCurrent"
                value={component.properties.maxCollectorCurrent || 1}
                onChange={handleChange}
              />
            </label>
          </>
        );
      case 'ic':
        return (
          <>
            <label>
              Type:
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
      <h3>{component.type.charAt(0).toUpperCase() + component.type.slice(1)} Properties</h3>
      {renderProperties()}
    </div>
  );
};

export default ComponentProperties;