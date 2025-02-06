import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

const StatusList = () => {
    const [statuses, setStatuses] = useState(null);

    useEffect(() => {
      const fetchStatuses = async () => {
        try {
          const response = await api.get('/statuses');
          console.log("API response:", response);
          setStatuses(response.data);
        } catch (error) {
          console.error('Error fetching statuses:', error);
          setStatuses([]);
        }
      };
        fetchStatuses();
    }, []);

    const handleDelete = async (id) => {
      try {
        await api.delete(`/statuses/${id}`);
          setStatuses(statuses.filter(status => status.status_id !== id))
      } catch (error) {
        console.error('Error deleting status:', error);
      }
  };

    return (
      <div>
          <h2>Statuses</h2>
          <Link to="/statuses/new">Create New Status</Link>
          {statuses && statuses.length > 0 ? (
          <table>
              <thead>
                <tr>
                      <th>ID</th>
                      <th>Name</th>
                      <th>Color</th>
                      <th>Actions</th>
                </tr>
                </thead>
              <tbody>
                {statuses.map(status => (
                  <tr key={status.status_id}>
                      <td>{status.status_id}</td>
                      <td>{status.name}</td>
                      <td>{status.color}</td>
                      <td>
                        <Link to={`/statuses/${status.status_id}/edit`}>Edit</Link>
                        <button onClick={() => handleDelete(status.status_id)}>Delete</button>
                      </td>
                  </tr>
                ))}
                </tbody>
          </table>
           ) : (
                <p>No statuses found</p>
          )}
      </div>
    );
};

export default StatusList;