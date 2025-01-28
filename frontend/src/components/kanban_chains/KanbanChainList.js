import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

const KanbanChainList = () => {
  const [kanbanChains, setKanbanChains] = useState(null);

  useEffect(() => {
    const fetchKanbanChains = async () => {
      try {
        const response = await api.get('/kanban-chains');
        console.log("API response:", response);
        setKanbanChains(response.data);
      } catch (error) {
        console.error('Error fetching kanban chains:', error);
        setKanbanChains([]); // Set to empty array on error
      }
    };
    fetchKanbanChains();
  }, []);

  const handleDelete = async (id) => {
    try {
      await api.delete(`/kanban-chains/${id}`);
      setKanbanChains(kanbanChains.filter(chain => chain.id !== id));
    } catch (error) {
        console.error("Error deleting kanban chain:", error);
    }
  };


  return (
    <div>
      <h2>Kanban Chains</h2>
          <Link to="/kanban-chains/new">Create New Kanban Chain</Link>
          {kanbanChains && kanbanChains.length > 0 ? (
          <table>
              <thead>
              <tr>
                  <th>ID</th>
                  <th>Customer ID</th>
                  <th>Supplier ID</th>
                  <th>Product ID</th>
                  <th>Lead Time</th>
                  <th>Container Type</th>
                  <th>Quantity</th>
                  <th>Status Chain</th>
                  <th>No Of Active Kanbans</th>
                <th>Actions</th>
              </tr>
              </thead>
              <tbody>
                {kanbanChains?.map(chain => (
                 <tr key={chain.id}>
                    <td>{chain.id}</td>
                    <td>{chain.cliente_id}</td>
                    <td>{chain.fornitore_id}</td>
                    <td>{chain.prodotto_codice}</td>
                    <td>{chain.leadtime_days}</td>
                    <td>{chain.tipo_contenitore}</td>
                    <td>{chain.quantity}</td>
                    <td>{chain.status_chain_id}</td>
                     <td>{chain.no_of_active_kanbans}</td>
                    <td>
                       <Link to={`/kanban-chains/${chain.id}/edit`}>Edit</Link>
                      <button onClick={() => handleDelete(chain.id)}>Delete</button>
                   </td>
                  </tr>
                ))}
              </tbody>
          </table>
          ) : (
            <p>No Kanban chains found</p>
          )}
    </div>
  );
};

export default KanbanChainList;