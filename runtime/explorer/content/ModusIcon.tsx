import React from "react";

const ModusLogo = ({ className }: { className?: string }) => (
  <svg
    id="Layer_1"
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 150 150"
    className={className}
  >
    <defs>
      <style>
        {`
          .cls-1{fill:#3fe6a1;}
          .cls-2{fill:#299e6e;}
          .cls-3{fill:#383838;}
        `}
      </style>
    </defs>
    <polygon
      className="cls-1"
      points="46.95 95.14 39.95 58.84 44.85 36.24 57.65 34.74 78.65 49.54 92.75 41.54 86.05 54.94 90.05 70.74 94.45 73.84 105.55 68.04 102.85 80.04 131.35 68.84 100.05 107.44 87.55 107.74 75.15 99.24 67.85 88.24 60.45 89.24 46.95 95.14"
    />
    <polygon
      className="cls-2"
      points="57.65 34.74 46.95 95.14 39.95 58.84 44.85 36.24 57.65 34.74"
    />
    <polygon
      className="cls-2"
      points="71.43 44.45 60.45 89.24 86.03 55.05 92.75 41.54 78.65 49.54 71.43 44.45"
    />
    <polygon
      className="cls-2"
      points="94.45 73.84 105.55 68.04 102.85 80.04 131.35 68.84 100.05 107.44 102.85 80.04 94.45 73.84"
    />
    <polygon
      className="cls-2"
      points="80.95 73.84 84.25 80.74 90.95 76.64 80.95 73.84"
    />
  </svg>
);

export default ModusLogo;
