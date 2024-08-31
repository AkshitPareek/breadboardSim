import React from 'react';

const IconWrapper = ({ children, color = "#000" }) => (
  <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    {children}
  </svg>
);

const ResistorIcon = () => (
  <IconWrapper>
    <rect x="10" y="18" width="20" height="4" fill="#000" />
    <rect x="5" y="15" width="30" height="10" stroke="#000" strokeWidth="2" fill="none" />
  </IconWrapper>
);

const CapacitorIcon = () => (
  <IconWrapper>
    <rect x="18" y="5" width="4" height="30" fill="#000" />
    <rect x="5" y="15" width="12" height="10" fill="#000" />
    <rect x="23" y="15" width="12" height="10" fill="#000" />
  </IconWrapper>
);

const InductorIcon = () => (
  <IconWrapper>
    <path d="M5 20 Q10 20 15 10 Q20 0 25 10 Q30 20 35 20" stroke="#000" strokeWidth="2" fill="none" />
  </IconWrapper>
);

const DiodeIcon = () => (
  <IconWrapper>
    <polygon points="15,10 15,30 30,20" fill="#000" />
    <line x1="30" y1="10" x2="30" y2="30" stroke="#000" strokeWidth="2" />
  </IconWrapper>
);

const LEDIcon = () => (
  <IconWrapper>
    <circle cx="20" cy="20" r="10" fill="#ff0" stroke="#000" strokeWidth="2" />
    <path d="M15 20 L25 20 M20 15 L20 25" stroke="#000" strokeWidth="2" />
  </IconWrapper>
);

const TransistorIcon = () => (
  <IconWrapper>
    <circle cx="20" cy="20" r="15" fill="none" stroke="#000" strokeWidth="2" />
    <path d="M10 30 L30 10" stroke="#000" strokeWidth="2" />
    <circle cx="15" cy="25" r="2" fill="#000" />
    <circle cx="25" cy="15" r="2" fill="#000" />
    <circle cx="25" cy="25" r="2" fill="#000" />
  </IconWrapper>
);

const ICIcon = () => (
  <IconWrapper>
    <rect x="5" y="10" width="30" height="20" fill="#000" />
    <circle cx="10" cy="15" r="2" fill="#fff" />
    <rect x="8" y="30" width="4" height="5" fill="#000" />
    <rect x="28" y="30" width="4" height="5" fill="#000" />
  </IconWrapper>
);

export const getComponentIcon = (type) => {
  switch (type) {
    case 'resistor':
      return <ResistorIcon />;
    case 'capacitor':
      return <CapacitorIcon />;
    case 'inductor':
      return <InductorIcon />;
    case 'diode':
      return <DiodeIcon />;
    case 'led':
      return <LEDIcon />;
    case 'transistor':
      return <TransistorIcon />;
    case 'ic':
      return <ICIcon />;
    default:
      return <span>?</span>;
  }
};