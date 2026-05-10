'use client';

import { useState, useEffect } from 'react';
import axios from 'axios';

export default function TestDashboard() {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    console.log('Test');
  }, []);

  return (
    <div className="container-fluid p-0">
      <div>Test</div>
    </div>
  );
}
