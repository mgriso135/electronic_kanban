import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

const StatusChainList = () => {
    const [statusChains, setStatusChains] = useState(null);

    useEffect(() => {
      const fetchStatusChains = async () => {
        try {
          const response = await api.get('/status-chains');
          console.log("API response:", response)
          setStatusChains(response.data);
        } catch (error) {
            console.error('Error fetching status chains:', error);
          setStatusChains([]);
        }
      };
        fetchStatusChains();
    }, []);

    const handleDelete = async (id) => {
      try {
        await api.delete(`/status-chains/${id}`);
        setStatusChains(statusChains.filter(chain => chain.status_chain_id !== id));
      } catch (error) {
          console.error("Error deleting status chain:", error)
      }
  };


    return (
      <div>
        <h2>Status Chains</h2>
          <Link to="/status-chains/new">Create New Status Chain</Link>
          {statusChains && statusChains.length > 0 ? (
          <table>
              <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {statusChains.map(chain => (
                  <tr key={chain.status_chain_id}>
                    <td>{chain.status_chain_id}</td>
                    <td>{chain.name}</td>
                    <td>
                        <Link to={`/status-chains/${chain.status_chain_id}/edit`}>Edit</Link>
                        <button onClick={() => handleDelete(chain.status_chain_id)}>Delete</button>
                    </td>
                  </tr>
                ))}
                </tbody>
            </table>
             ) : (
                <p>No status chains found</p>
          )}
      </div>
    );
};

export default StatusChainList;