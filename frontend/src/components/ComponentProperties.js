import React, { useState, useEffect } from 'react';
import '../css/ComponentProperties.css';

const COMPONENT_PROPERTIES = {
  resistor: [
    { name: 'resistance', label: 'Resistance (Ω)', type: 'number', unit: 'ohm' },
    { name: 'powerRating', label: 'Power Rating (W)', type: 'number', unit: 'watt' },
    { name: 'tolerance', label: 'Tolerance (%)', type: 'number', unit: 'percent' },
  ],
  capacitor: [
    { name: 'capacitance', label: 'Capacitance (µF)', type: 'number', unit: 'microfarad' },
    { name: 'voltageRating', label: 'Voltage Rating (V)', type: 'number', unit: 'volt' },
    { name: 'capacitorType', label: 'Type', type: 'select', options: ['Electrolytic', 'Ceramic', 'Film'] }
  ],
  inductor: [
    { name: 'inductance', label: 'Inductance (H)', type: 'number', unit: 'henry' },
    { name: 'currentRating', label: 'Current Rating (A)', type: 'number', unit: 'ampere' }
  ],
  diode: [
    { name: 'forwardVoltage', label: 'Forward Voltage (V)', type: 'number', unit: 'volt' },
    { name: 'maxCurrent', label: 'Max Current (mA)', type: 'number', unit: 'milliampere' }
  ],
  led: [
    { name: 'forwardVoltage', label: 'Forward Voltage (V)', type: 'number', unit: 'volt' },
    { name: 'maxCurrent', label: 'Max Current (mA)', type: 'number', unit: 'milliampere' },
    { name: 'color', label: 'Color', type: 'select', options: ['Red', 'Green', 'Blue', 'Yellow'] }
  ],
  transistor: [
    { name: 'transistorType', label: 'Type', type: 'select', options: ['NPN', 'PNP'] },
    { name: 'gain', label: 'Gain (hFE)', type: 'number' },
    { name: 'maxCollectorCurrent', label: 'Max Collector Current (A)', type: 'number', unit: 'ampere' }
  ],
  ic: [
    { name: 'icType', label: 'IC Type', type: 'text' },
    { name: 'description', label: 'Description', type: 'textarea' }
  ],
  battery: [
    { name: 'voltage', label: 'Voltage (V)', type: 'number', unit: 'volt' },
    { name: 'capacity', label: 'Capacity (mAh)', type: 'number', unit: 'milliamp-hour' }
  ],
  power_supply: [
    { name: 'voltage', label: 'Voltage (V)', type: 'number', unit: 'volt' },
    { name: 'maxCurrent', label: 'Max Current (A)', type: 'number', unit: 'ampere' }
  ]
};

const ComponentProperties = ({ component, onUpdate, onClose }) => {
  const [localProperties, setLocalProperties] = useState(component.properties);

  useEffect(() => {
    setLocalProperties(component.properties);
  }, [component]);

  const handleChange = (name, value, type) => {
    console.log(`Changing ${name} to ${value} (type: ${type})`);
    let updatedValue = value;

    if (type === 'number') {
      const numericRegex = /^-?\d*\.?\d*$/;
      if (numericRegex.test(value)) {
        updatedValue = value === '' ? '' : parseFloat(value);
      } else {
        return; // Invalid numeric input, don't update
      }
    }

    const newProperties = { ...localProperties, [name]: updatedValue };
    setLocalProperties(newProperties);
    onUpdate(component.id, newProperties);
  };

  const renderInput = (prop) => {
    const value = localProperties[prop.name] || '';

    switch (prop.type) {
      case 'number':
        return (
          <input
            type="text"
            name={prop.name}
            value={value}
            onChange={(e) => handleChange(prop.name, e.target.value, 'number')}
            onBlur={(e) => {
              const validatedValue = e.target.value === '' ? '' : parseFloat(e.target.value) || 0;
              handleChange(prop.name, validatedValue, 'number');
            }}
          />
        );
      case 'select':
        return (
          <select
            name={prop.name}
            value={value}
            onChange={(e) => handleChange(prop.name, e.target.value, 'select')}
          >
            <option value="">Select {prop.label}</option>
            {prop.options.map((option) => (
              <option key={option} value={option.toLowerCase()}>{option}</option>
            ))}
          </select>
        );
      case 'textarea':
        return (
          <textarea
            name={prop.name}
            value={value}
            onChange={(e) => handleChange(prop.name, e.target.value, 'textarea')}
          />
        );
      default:
        return (
          <input
            type="text"
            name={prop.name}
            value={value}
            onChange={(e) => handleChange(prop.name, e.target.value, 'text')}
          />
        );
    }
  };

  const renderProperties = () => {
    const properties = COMPONENT_PROPERTIES[component.type] || [];
    return properties.map((prop) => (
      <label key={prop.name}>
        {prop.label}:
        {renderInput(prop)}
        {prop.unit && <span className="unit">{prop.unit}</span>}
      </label>
    ));
  };

  return (
    <div className="component-properties">
      <button className="close-button" onClick={onClose}>&times;</button>
      <h3>{component.type.charAt(0).toUpperCase() + component.type.slice(1)} Properties</h3>
      {renderProperties()}
    </div>
  );
};

export default ComponentProperties;